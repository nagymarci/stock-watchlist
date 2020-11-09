package controllers

import (
	"errors"

	"github.com/nagymarci/stock-watchlist/api"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nagymarci/stock-watchlist/database"

	"github.com/nagymarci/stock-watchlist/model"
	"github.com/nagymarci/stock-watchlist/service"

	stockHttp "github.com/nagymarci/stock-commons/http"
	userprofileModel "github.com/nagymarci/stock-user-profile/model"
)

type WatchlistController struct {
	watchlists        *database.Watchlists
	stockClient       *api.StockClient
	userprofileClient *api.UserprofileClient
	stockService      *service.StockService
}

func NewWatchlistController(w *database.Watchlists, sc *api.StockClient) *WatchlistController {
	return &WatchlistController{
		watchlists:  w,
		stockClient: sc,
	}
}

//Create creates a new watchlist
func (wl *WatchlistController) Create(log *logrus.Logger, request *model.WatchlistRequest) (*model.Watchlist, error) {
	var addedStocks []string

	for _, symbol := range request.Stocks {
		err := wl.stockClient.RegisterStock(symbol)

		if err != nil {
			log.WithField("symbol", symbol).Warnln(err)
			continue
		}

		addedStocks = append(addedStocks, symbol)
	}

	request.Stocks = addedStocks
	id, err := wl.watchlists.Create(*request)

	if err != nil {
		return nil, stockHttp.NewInternalServerError(err.Error())
	}

	watchlistResponse := model.Watchlist{
		ID:     id,
		Name:   request.Name,
		Stocks: request.Stocks,
		UserID: request.UserID}

	return &watchlistResponse, err
}

//Delete deletes the specified watchlist if that belongs to the authorized user
func (wl *WatchlistController) Delete(log *logrus.Logger, id primitive.ObjectID, userID string) error {
	_, err := wl.getAndValidateUserAuthorization(id, userID)

	if err != nil {
		return stockHttp.NewBadRequestError(err.Error())
	}

	result, err := wl.watchlists.Delete(id)

	if result != 1 {
		return stockHttp.NewInternalServerError("No object were removed from database")
	}

	if err != nil {
		return stockHttp.NewInternalServerError(err.Error())
	}

	return nil
}

func (wl *WatchlistController) Get(log *logrus.Logger, id primitive.ObjectID, userID string) (model.Watchlist, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, userID)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Errorln(message)
		return model.Watchlist{}, stockHttp.NewBadRequestError(message)
	}

	return watchlist, nil
}

func (wl *WatchlistController) GetAll(log *logrus.Logger, userID string) ([]model.Watchlist, error) {
	watchlists, err := wl.watchlists.GetAll(userID)

	if err != nil {
		message := "Unable to list watchlists " + err.Error()
		log.Errorln(message)
		return nil, stockHttp.NewBadRequestError(message)
	}

	return watchlists, nil
}

func (wl *WatchlistController) GetCalculated(log *logrus.Logger, id primitive.ObjectID, userID string) ([]model.CalculatedStockInfo, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, userID)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Errorln(message)
		return nil, stockHttp.NewBadRequestError(message)
	}

	var stockInfos []model.CalculatedStockInfo

	userprofile, err := wl.userprofileClient.GetUserprofile(userID)

	if err != nil {
		log.Errorln(err)
		defaultExpectation := 9.0
		defaultExpectedReturn := 9.0
		userprofile = userprofileModel.Userprofile{DefaultExpectation: &defaultExpectation, ExpectedReturn: &defaultExpectedReturn}
	}

	for _, symbol := range watchlist.Stocks {
		result, err := wl.stockClient.Get(symbol)

		if err != nil {
			log.Warnf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		expectation := userprofile.GetExpectation(symbol)

		log.Debugf("Symbol [%s] expectation [%f]\n", symbol, expectation)

		calculatedStockInfo := wl.stockService.Calculate(&result, expectation, *userprofile.ExpectedReturn)

		stockInfos = append(stockInfos, calculatedStockInfo)
	}

	return stockInfos, nil
}

func (w *WatchlistController) getAndValidateUserAuthorization(id primitive.ObjectID, userID string) (model.Watchlist, error) {
	watchlist, err := w.watchlists.Get(id)
	if err != nil {
		return watchlist, err
	}

	if watchlist.UserID != userID {
		return watchlist, errors.New("Watchlist does not belong to user")
	}

	return watchlist, err
}
