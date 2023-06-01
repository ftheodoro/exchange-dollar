package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	requrl, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	req, err := http.DefaultClient.Do(requrl)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var exchangeRate string

	if err := json.Unmarshal(res, &exchangeRate); err != nil {
		panic(err)
	}
	saveDataFile(exchangeRate)
	io.Copy(os.Stdout, req.Body)
}

func saveDataFile(exchange string) {
	file, err := os.Create("exchangerate.txt")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err != nil {
		panic(err)
	}
	dolarRate := fmt.Sprintf(" Dólar: %s", exchange)
	_, err = file.WriteString(dolarRate)
	if err != nil {
		panic(err)
	}

}
