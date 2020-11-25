package controllers

import (
	userprofileModel "github.com/nagymarci/stock-user-profile/model"
	"github.com/nagymarci/stock-watchlist/api"
	"github.com/nagymarci/stock-watchlist/model"
	"github.com/nagymarci/stock-watchlist/service"
	"github.com/sirupsen/logrus"

	stockHttp "github.com/nagymarci/stock-commons/http"
)

type StockController struct {
	stockClient       *api.StockClient
	userprofileClient *api.UserprofileClient
	stockService      *service.StockService
}

func NewStockController(sc *api.StockClient, upC *api.UserprofileClient, ss *service.StockService) *StockController {
	return &StockController{
		stockClient:       sc,
		userprofileClient: upC,
		stockService:      ss,
	}
}

func (sc *StockController) GetAllCalculated(log *logrus.Entry, userID string) ([]model.CalculatedStockInfo, error) {
	stocks, err := sc.stockClient.GetAll()

	if err != nil {
		log.Errorln(err)
		return nil, stockHttp.NewFailedDependencyError(err.Error())
	}

	var userprofile userprofileModel.Userprofile
	if userID != "" {
		userprofile, err = sc.userprofileClient.GetUserprofile(userID)

		if err != nil {
			log.Errorln(err)
		}
	}

	if err != nil || userID == "" {
		defaultExpectation := 9.0
		defaultExpectedReturn := 9.0
		userprofile = userprofileModel.Userprofile{DefaultExpectation: &defaultExpectation, ExpectedReturn: &defaultExpectedReturn}
	}

	var stockInfos []model.CalculatedStockInfo
	for _, stock := range stocks {

		expectation := userprofile.GetExpectation(stock.Ticker)

		log.Debugf("Symbol [%s] expectation [%f]\n", stock.Ticker, expectation)

		calculatedStockInfo := sc.stockService.Calculate(&stock, expectation, *userprofile.ExpectedReturn)

		stockInfos = append(stockInfos, calculatedStockInfo)
	}

	return stockInfos, nil
}
