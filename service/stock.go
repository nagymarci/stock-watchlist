package service

import (
	"math"

	"github.com/nagymarci/stock-watchlist/model"
)

type sP500Client interface {
	GetSP500DivYield() float64
}

type StockService struct {
	sp500Client sP500Client
}

func NewStockService(s sP500Client) *StockService {
	return &StockService{
		sp500Client: s,
	}
}

const (
	defaultDividendPerYear       float64 = 4
	monthlyDividendPerYear       float64 = 12
	lowerDividendYieldGuardScore float64 = 1.5
	maxOptInPeWeight             float64 = 0.5
	minOptInYieldWeight          float64 = 0.4
)

//Calculate returns the dynamically computed data from the latest information
func (ss *StockService) Calculate(stockInfo *model.StockData, expectedRaise float64, expectedReturn float64) model.CalculatedStockInfo {
	var result model.CalculatedStockInfo

	sp500DivYield := ss.sp500Client.GetSP500DivYield()

	minYieldFromExpRaise := expectedReturn - expectedRaise
	if minYieldFromExpRaise <= 0.0 {
		minYieldFromExpRaise = 0.1
	}

	optInYield, minOptInYield := calculateOptInYield(stockInfo.DividendYield5yr.Max, stockInfo.DividendYield5yr.Avg, sp500DivYield, minYieldFromExpRaise)

	optInPe := calculateOptInPe(stockInfo.PeRatio5yr.Min, stockInfo.PeRatio5yr.Avg)

	result.Ticker = stockInfo.Ticker
	result.AnnualDividend = stockInfo.Dividend * defaultDividendPerYear

	//TODO store this in DB with the other info for the given stock
	if result.Ticker == "O" {
		result.AnnualDividend = stockInfo.Dividend * monthlyDividendPerYear
	}
	result.Price = stockInfo.Price
	result.DividendYield = result.AnnualDividend / result.Price * 100
	result.CurrentPe = result.Price / stockInfo.Eps
	if stockInfo.Eps == 0 {
		result.CurrentPe = math.MaxFloat64
	}
	result.OptInYield = optInYield
	result.DividendColor = calculateDividendColor(result.DividendYield, minOptInYield, stockInfo.DividendYield5yr.Avg)
	result.OptInPe = optInPe
	result.PeColor = calculatePeColor(result.CurrentPe, optInPe, stockInfo.PeRatio5yr.Avg)

	optInPrice := calculateOptInPrice(optInYield, result.AnnualDividend, sp500DivYield, minYieldFromExpRaise)

	result.OptInPrice = optInPrice
	result.PriceColor = calculatePriceColor(result.Price, optInPrice)

	return result
}

func calculatePriceColor(price float64, optInPrice float64) string {
	if price < optInPrice {
		return "green"
	}
	if price < optInPrice*1.05 {
		return "yellow"
	}

	return "red"
}

func calculateOptInPrice(optInYield float64, annualDividend float64, sp float64, minYieldFromRaise float64) float64 {
	spOptInPrice := annualDividend / (sp * lowerDividendYieldGuardScore) * 100
	minOptInPrice := annualDividend / optInYield * 100
	expectedRaiseOptInPrice := annualDividend / minYieldFromRaise * 100

	return math.Min(spOptInPrice, math.Min(minOptInPrice, expectedRaiseOptInPrice))
}

func calculatePeColor(currentPe float64, optInPe float64, avg float64) string {
	if currentPe < optInPe {
		return "green"
	}

	if currentPe < avg {
		return "yellow"
	}

	return "blank"
}

func calculateOptInPe(min float64, avg float64) float64 {
	return (avg-min)*maxOptInPeWeight + min
}

func calculateDividendColor(dividendYield float64, minOptInYield float64, avg float64) string {
	if dividendYield > minOptInYield {
		return "green"
	}
	if dividendYield > avg {
		return "yellow"
	}

	return "blank"
}

func calculateOptInYield(max float64, avg float64, sp float64, exp float64) (float64, float64) {
	minOptInYield := calculateMinOptInYield(max, avg)
	return math.Max(minOptInYield, math.Max(sp*lowerDividendYieldGuardScore, exp)), minOptInYield
}

func calculateMinOptInYield(max float64, avg float64) float64 {
	return (max-avg)*minOptInYieldWeight + avg
}
