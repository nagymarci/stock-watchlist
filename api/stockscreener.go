package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nagymarci/stock-watchlist/model"
)

type StockClient struct {
	host string
}

func NewStockClient(h string) *StockClient {
	return &StockClient{
		host: h,
	}
}

func (sc *StockClient) RegisterStock(symbol string) error {

	resp, err := http.Post(sc.host+symbol, "", nil)

	if err != nil {
		return fmt.Errorf("Failed to register stock [%s] with error [%v]", symbol, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return fmt.Errorf("Failed to register [%s], status code [%d], response [%v]", symbol, resp.StatusCode, response)
	}

	return nil
}

func (sc *StockClient) Get(symbol string) (model.StockData, error) {
	resp, err := http.Get(sc.host + symbol)

	stockData := model.StockData{}

	if err != nil {
		return stockData, fmt.Errorf("Failed to get stock [%s] with error [%v]", symbol, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return stockData, fmt.Errorf("Failed to get [%s], status code [%d], response [%v]", symbol, resp.StatusCode, response)
	}

	err = json.NewDecoder(resp.Body).Decode(&stockData)

	if err != nil {
		return stockData, fmt.Errorf("Failed to deserialize data for [%s], error: [%v]", symbol, err)
	}

	return stockData, nil
}

type sp500DivYield struct {
	Yield      float64
	NextUpdate time.Time
	Mux        sync.Mutex
}

//Sp500DivYield stores information of the S&P500 dividend yield, and when we should update it next
var sp500 sp500DivYield

//TODO move this logic to stock-screener service
func (sc *StockClient) GetSP500DivYield() float64 {
	now := time.Now()
	if sp500.NextUpdate.Before(now) {
		log.Println("Before lock")
		sp500.Mux.Lock()
		log.Println("After lock")

		defer sp500.Mux.Unlock()

		if sp500.NextUpdate.Before(now) {
			yield, err := getSp500DivYield()
			if err != nil {
				log.Printf("Failed to update sp500 dividend yield: [%v]\n", err)
				log.Println("Using old sp500 dividend yield")

			} else {
				nextUpdateInterval, err := time.ParseDuration("12h")
				if err != nil {
					log.Printf("Error when parsing duration [%v]\n", err)
				}
				sp500.Yield = yield
				sp500.NextUpdate = now.Add(nextUpdateInterval)
				log.Println("SP500 dividend yield updated")
			}
		}
	}

	return sp500.Yield
}

func getSp500DivYield() (float64, error) {
	host := os.Getenv("SP500_URL")

	resp, err := http.Get(host)

	if err != nil {
		return 0, fmt.Errorf("Failed to get SP500 div yield: [%v]", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return 0, fmt.Errorf("Failed to get SP500 div yield:, status code [%d], response [%v]", resp.StatusCode, response)
	}

	var response float64
	fmt.Fscan(resp.Body, &response)

	return response, nil
}
