package service

import (
	"testing"

	"github.com/nagymarci/stock-watchlist/model"
)

type mockSp500Client struct{}

func (sp *mockSp500Client) GetSP500DivYield() float64 {
	return 1.0
}

func TestStockCalculate(t *testing.T) {
	t.Run("math works", func(t *testing.T) {
		stockService := NewStockService(&mockSp500Client{})

		stock := model.StockData{}
		stock.Ticker = "INTC"
		stock.Dividend = 0.33
		stock.Eps = 5.43
		stock.Price = 49.28
		stock.DividendYield5yr.Avg = 2.62
		stock.DividendYield5yr.Max = 3.65
		stock.PeRatio5yr.Avg = 14.89
		stock.PeRatio5yr.Min = 8.79

		expectedResult := model.CalculatedStockInfo{}
		expectedResult.Ticker = stock.Ticker
		expectedResult.AnnualDividend = 1.32
		expectedResult.CurrentPe = 9.075506445672191
		expectedResult.OptInPe = 11.84
		expectedResult.PeColor = "green"
		expectedResult.Price = stock.Price
		expectedResult.OptInPrice = 37.714285714285715
		expectedResult.PriceColor = "red"
		expectedResult.DividendYield = 2.678571428571429
		expectedResult.OptInYield = 3.5
		expectedResult.DividendColor = "yellow"

		result := stockService.Calculate(&stock, 5.5, 9.0)

		if result != expectedResult {
			t.Errorf("expected [%v], got [%v]", expectedResult, result)
		}

	})
}
