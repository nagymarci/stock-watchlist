package service

import (
	"testing"

	"github.com/nagymarci/stock-watchlist/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	userprofileModel "github.com/nagymarci/stock-user-profile/model"
	"github.com/nagymarci/stock-watchlist/service/mocks"

	"github.com/golang/mock/gomock"
)

func TestNotification(t *testing.T) {
	t.Run("no email if nothing changed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recommendations := mocks.NewMockrecommendationProvider(ctrl)
		watchlists := mocks.NewMockwatchlistList(ctrl)
		stockClient := mocks.NewMockstockGetter(ctrl)
		stockService := mocks.NewMockstockRecommendator(ctrl)
		userprofileClient := mocks.NewMockuserprofileGetter(ctrl)
		emailClient := mocks.NewMockemailSender(ctrl)

		notifier := NewNotifier(recommendations, watchlists, stockClient, stockService, userprofileClient, emailClient)

		watchlistID := primitive.NewObjectID()
		expectedWatchlist := model.Watchlist{ID: watchlistID, Name: "watchlist", Stocks: []string{"INTC"}, UserID: "userId"}

		stock := model.StockData{}
		stock.Ticker = "INTC"
		stock.Dividend = 0.33
		stock.Eps = 5.43
		stock.Price = 49.28
		stock.DividendYield5yr.Avg = 2.62
		stock.DividendYield5yr.Max = 3.65
		stock.PeRatio5yr.Avg = 14.89
		stock.PeRatio5yr.Min = 8.79

		expectedReturn := 9.0
		expectedRaise := 5.5
		userprofile := userprofileModel.Userprofile{Email: "alice@example.com", ExpectedReturn: &expectedReturn, Expectations: []userprofileModel.Expectation{userprofileModel.Expectation{Stock: "INTC", ExpectedRaise: &expectedRaise}}}

		watchlists.EXPECT().List().Return([]model.Watchlist{expectedWatchlist}, nil)
		recommendations.EXPECT().Get(watchlistID).Return([]string{}, nil)
		stockClient.EXPECT().Get("INTC").Return(stock, nil)
		userprofileClient.EXPECT().GetUserprofile("userId").Return(userprofile, nil)
		stockService.EXPECT().GetAllRecommendedStock(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.CalculatedStockInfo{})
		emailClient.EXPECT().SendNotification(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		notifier.NotifyChanges()
	})
	t.Run("email when stock removed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recommendations := mocks.NewMockrecommendationProvider(ctrl)
		watchlists := mocks.NewMockwatchlistList(ctrl)
		stockClient := mocks.NewMockstockGetter(ctrl)
		stockService := mocks.NewMockstockRecommendator(ctrl)
		userprofileClient := mocks.NewMockuserprofileGetter(ctrl)
		emailClient := mocks.NewMockemailSender(ctrl)

		notifier := NewNotifier(recommendations, watchlists, stockClient, stockService, userprofileClient, emailClient)

		watchlistID := primitive.NewObjectID()
		expectedWatchlist := model.Watchlist{ID: watchlistID, Name: "watchlist", Stocks: []string{"INTC"}, UserID: "userId"}

		stock := model.StockData{}
		stock.Ticker = "INTC"
		stock.Dividend = 0.33
		stock.Eps = 5.43
		stock.Price = 49.28
		stock.DividendYield5yr.Avg = 2.62
		stock.DividendYield5yr.Max = 3.65
		stock.PeRatio5yr.Avg = 14.89
		stock.PeRatio5yr.Min = 8.79

		expectedReturn := 9.0
		expectedRaise := 5.5
		userprofile := userprofileModel.Userprofile{Email: "alice@example.com", ExpectedReturn: &expectedReturn, Expectations: []userprofileModel.Expectation{userprofileModel.Expectation{Stock: "INTC", ExpectedRaise: &expectedRaise}}}

		var empty []string

		watchlists.EXPECT().List().Return([]model.Watchlist{expectedWatchlist}, nil)
		recommendations.EXPECT().Get(watchlistID).Return([]string{"INTC"}, nil)
		recommendations.EXPECT().Update(gomock.Any(), watchlistID, empty).Return(nil)
		stockClient.EXPECT().Get("INTC").Return(stock, nil)
		userprofileClient.EXPECT().GetUserprofile("userId").Return(userprofile, nil)
		stockService.EXPECT().GetAllRecommendedStock(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.CalculatedStockInfo{})
		emailClient.EXPECT().SendNotification(expectedWatchlist.Name, []string{"INTC"}, empty, empty, userprofile.Email).Times(1)

		notifier.NotifyChanges()
	})
	t.Run("email when stock added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recommendations := mocks.NewMockrecommendationProvider(ctrl)
		watchlists := mocks.NewMockwatchlistList(ctrl)
		stockClient := mocks.NewMockstockGetter(ctrl)
		stockService := mocks.NewMockstockRecommendator(ctrl)
		userprofileClient := mocks.NewMockuserprofileGetter(ctrl)
		emailClient := mocks.NewMockemailSender(ctrl)

		notifier := NewNotifier(recommendations, watchlists, stockClient, stockService, userprofileClient, emailClient)

		watchlistID := primitive.NewObjectID()
		expectedWatchlist := model.Watchlist{ID: watchlistID, Name: "watchlist", Stocks: []string{"INTC"}, UserID: "userId"}

		stock := model.StockData{}
		stock.Ticker = "INTC"
		stock.Dividend = 0.33
		stock.Eps = 5.43
		stock.Price = 37
		stock.DividendYield5yr.Avg = 2.62
		stock.DividendYield5yr.Max = 3.65
		stock.PeRatio5yr.Avg = 14.89
		stock.PeRatio5yr.Min = 8.79

		expectedReturn := 9.0
		expectedRaise := 5.5
		userprofile := userprofileModel.Userprofile{Email: "alice@example.com", ExpectedReturn: &expectedReturn, Expectations: []userprofileModel.Expectation{userprofileModel.Expectation{Stock: "INTC", ExpectedRaise: &expectedRaise}}}

		calculatedStockInfo := model.CalculatedStockInfo{}
		calculatedStockInfo.Ticker = "INTC"
		calculatedStockInfo.PriceColor = "green"

		var empty []string

		watchlists.EXPECT().List().Return([]model.Watchlist{expectedWatchlist}, nil)
		recommendations.EXPECT().Get(watchlistID).Return(empty, nil)
		recommendations.EXPECT().Update(gomock.Any(), watchlistID, []string{"INTC"}).Return(nil)
		stockClient.EXPECT().Get("INTC").Return(stock, nil)
		userprofileClient.EXPECT().GetUserprofile("userId").Return(userprofile, nil)
		stockService.EXPECT().GetAllRecommendedStock(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.CalculatedStockInfo{calculatedStockInfo})
		emailClient.EXPECT().SendNotification(expectedWatchlist.Name, empty, []string{"INTC"}, []string{"INTC"}, userprofile.Email).Times(1)

		notifier.NotifyChanges()
	})
}
