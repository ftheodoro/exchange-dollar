package main

import (
	"context"
	"encoding/json"
	"github.com/ftheodoro/exchange-dollar/model"
	"io"
	"net/http"
	"time"
)

const CoinName = "USDBRL"
const TimeRequest = 2

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)

}
func Index(w http.ResponseWriter, r *http.Request) {
	jsonData, err, httpStatus := requestExhange()
	if err != nil {
		w.WriteHeader(httpStatus)
		return
	}

	var exchangeRate model.ExchangeRate

	if err = json.Unmarshal(jsonData, &exchangeRate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchangeRate)

}

func requestExhange() ([]byte, error, int) {
	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TimeRequest)
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
