package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ftheodoro/exchange-dollar/config"
	"github.com/ftheodoro/exchange-dollar/model"
)

const CoinName = "USDBRL"
const TimeAPIRequest = 200
const TimeSaveDataBase = 10
const TimeClientRequest = 300

func main() {
	log.Println("server started")
	http.HandleFunc("/cotacao", Index)
	http.ListenAndServe(":8080", nil)

}
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	select {
	case <-time.After(time.Millisecond * TimeClientRequest):

		log.Println("client-initiated request")
		defer log.Println("finish request ")

		jsonData, httpStatus, err := requestExhange()
		if err != nil {
			log.Println(err)
			w.WriteHeader(httpStatus)
			return
		}

		var exchangeRate model.ExchangeRate

		if err = json.Unmarshal(jsonData, &exchangeRate); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		db, err := config.ConnDB()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error creating database."))

			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*TimeSaveDataBase)
		defer cancel()
		err = db.WithContext(ctx).Create(exchangeRate).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("timeout exceeded for saving to database"))
			return
		}

		json.NewEncoder(w).Encode(exchangeRate.Bid)
	case <-r.Context().Done():
		log.Println("request cancel client")

	}

}

func requestExhange() ([]byte, int, error) {
	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*TimeAPIRequest)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []byte(""), http.StatusBadRequest, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Timeout request")
		return []byte(""), http.StatusRequestTimeout, err
	}
	defer resp.Body.Close()

	var apiData map[string]interface{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error read body request")
		return []byte(""), http.StatusBadRequest, err
	}

	if err = json.Unmarshal(body, &apiData); err != nil {
		log.Println("Error unmarshal body")
		return []byte(""), http.StatusBadRequest, err
	}

	jsonData, err := json.Marshal(apiData[CoinName])

	if err != nil {
		log.Println("Error marshal body")
		return []byte(""), http.StatusBadGateway, err
	}
	return jsonData, http.StatusOK, nil
}
