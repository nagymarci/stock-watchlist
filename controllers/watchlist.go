package controllers

import (
	"errors"

	"github.com/nagymarci/stock-watchlist/api"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nagymarci/stock-watchlist/database"

	"github.com/nagymarci/stock-watchlist/model"
	"github.com/nagymarci/stock-watchlist/service"

	userprofileModel "github.com/nagymarci/stock-user-profile/model"
)

type WatchlistController struct {
	watchlists   database.WatchlistCollection
	stockService service.Stock
}

func NewWatchlistController(w database.WatchlistCollection) *WatchlistController {
	return &WatchlistController{
		watchlists: w,
	}
}

//Create creates a new watchlist
func (wl *WatchlistController) Create(log *logrus.Logger, request *model.WatchlistRequest) (*model.Watchlist, error) {
	var addedStocks []string

	for _, symbol := range request.Stocks {
		err := stockService.SaveStock(symbol)

		if err != nil {
			log.Warnln(err)
			continue
		}

		addedStocks = append(addedStocks, symbol)
	}

	request.Stocks = addedStocks
	id, err := wl.watchlists.Create(*request)

	if err != nil {
		return nil, model.NewInternalServerError(err.Error())
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
		return model.NewBadRequestError(err.Error())
	}

	result, err := wl.watchlists.Delete(id)

	if result != 1 {
		return model.NewInternalServerError("No object were removed from database")
	}

	if err != nil {
		return model.NewInternalServerError(err.Error())
	}

	return nil
}

func (wl *WatchlistController) Get(log *logrus.Logger, id primitive.ObjectID, userID string) (model.Watchlist, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, userID)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Errorln(message)
		return model.Watchlist{}, model.NewBadRequestError(message)
	}

	return watchlist, nil
}

func (wl *WatchlistController) GetAll(log *logrus.Logger, userID string) ([]model.Watchlist, error) {
	watchlists, err := wl.watchlists.GetAll(userID)

	if err != nil {
		message := "Unable to list watchlists " + err.Error()
		log.Errorln(message)
		return nil, model.NewBadRequestError(message)
	}

	return watchlists, nil
}

func (wl *WatchlistController) GetCalculated(log *logrus.Logger, id primitive.ObjectID, userID string) ([]model.CalculatedStockInfo, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, userID)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Errorln(message)
		return nil, model.NewBadRequestError(message)
	}

	var stockInfos []model.CalculatedStockInfo

	userprofile, err := api.GetUserprofile(userID)

	if err != nil {
		log.Errorln(err)
		defaultExpectation := 9.0
		userprofile = userprofileModel.Userprofile{DefaultExpectation: &defaultExpectation}
	}

	for _, symbol := range watchlist.Stocks {
		result, err := database.Get(symbol)

		if err != nil {
			log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		expectation := userprofile.GetExpectation(symbol)

		log.Printf("Symbol [%s] expectation [%f]", symbol, expectation)

		calculatedStockInfo := service.Calculate(&result, expectation)

		stockInfos = append(stockInfos, calculatedStockInfo)
	}

	return stockInfos, nil
}

func (w *WatchlistController) getAndValidateUserAuthorization(id primitive.ObjectID, userID string) (model.Watchlist, error) {
	watchlist, err := w.watchlists.Get(id)
	if err != nil {
		return watchlist, err
	}

	if watchlist.userID != userID {
		return watchlist, errors.New("Watchlist does not belong to user")
	}

	return watchlist, err
}

func saveStock(stock string) error {
	_, err := database.Get(stock)

	if err == nil {
		return err
	}

	stockData, err := service.Get(stock)

	if err != nil {
		return err
	}

	err = database.Save(stockData)

	return err
}
