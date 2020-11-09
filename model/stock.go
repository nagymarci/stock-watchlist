package model

type pERatioInfo struct {
	Avg float64 `json:"avg"`
	Min float64 `json:"min"`
}

type dividendYieldInfo struct {
	Avg float64 `json:"avg"`
	Max float64 `json:"max"`
}

//StockData holds the information for one stock
type StockData struct {
	Ticker           string            `json:"ticker"`
	Price            float64           `json:"price"`
	Eps              float64           `json:"eps"`
	Dividend         float64           `json:"dividend"`
	PeRatio5yr       pERatioInfo       `json:"peRatio5yr"`
	DividendYield5yr dividendYieldInfo `json:"dividendYield5yr"`
}

//CalculatedStockInfo holds the data calculated for investment suggestions
type CalculatedStockInfo struct {
	Ticker         string  `json:"ticker"`
	Price          float64 `json:"price"`
	OptInPrice     float64 `json:"optInPrice"`
	PriceColor     string  `json:"priceColor"`
	AnnualDividend float64 `json:"dividend"`
	DividendYield  float64 `json:"dividendYield"`
	OptInYield     float64 `json:"optInYield"`
	DividendColor  string  `json:"dividendColor"`
	CurrentPe      float64 `json:"currentPe"`
	OptInPe        float64 `json:"optInPe"`
	PeColor        string  `json:"pecolor"`
}
