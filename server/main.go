package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ftheodoro/exchange-dollar/config"
	"github.com/ftheodoro/exchange-dollar/model"
)

const CoinName = "USDBRL"
const TimeRequest = 200
const TimeSaveDataBase = 10

func main() {

	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)

}
func Index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*TimeSaveDataBase)
	defer cancel()
	jsonData, err, httpStatus := requestExhange()
	if err != nil {
		w.WriteHeader(httpStatus)
		return
	}

	var exchangeRate model.ExchangeRate
	w.Header().Set("Content-Type", "application/json")
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
	err = db.WithContext(ctx).Create(exchangeRate).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("timeout exceeded for saving to database"))
		return
	}

	json.NewEncoder(w).Encode(exchangeRate)

}

func requestExhange() ([]byte, error, int) {
	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*TimeRequest)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []byte(""), err, http.StatusBadRequest
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err, http.StatusRequestTimeout
	}
	defer resp.Body.Close()

	var apiData map[string]interface{}

	body, err := io.ReadAll(resp.Body)

	if err = json.Unmarshal(body, &apiData); err != nil {
		return []byte(""), err, http.StatusBadRequest
	}

	jsonData, err := json.Marshal(apiData[CoinName])

	if err != nil {
		return []byte(""), err, http.StatusBadGateway
	}
	return jsonData, nil, http.StatusOK
}
