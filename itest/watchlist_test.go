package itest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	userprofileModel "github.com/nagymarci/stock-user-profile/model"
	"github.com/nagymarci/stock-watchlist/service"

	"github.com/nagymarci/stock-watchlist/model"

	"github.com/nagymarci/stock-watchlist/handlers"

	"github.com/golang/mock/gomock"
	"github.com/nagymarci/stock-watchlist/controllers"
	"github.com/nagymarci/stock-watchlist/database"
	"github.com/nagymarci/stock-watchlist/itest/mocks"

	"github.com/gorilla/mux"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockSp500Client struct{}

func (sp *mockSp500Client) GetSP500DivYield() float64 {
	return 1.0
}

var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(time.Minute * 2),
		Env:          map[string]string{},
	}
	req.Env["MONGO_INITDB_ROOT_USERNAME"] = "mongodb"
	req.Env["MONGO_INITDB_ROOT_PASSWORD"] = "mongodb"
	req.Env["MONGO_INITDB_DATABASE"] = "stock-screener"

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer mongoC.Terminate(ctx)
	ip, err := mongoC.Host(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		log.Fatalln(err)
	}

	dbConnectionURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		"mongodb",
		"mongodb",
		ip,
		port.Int())

	db = database.New(dbConnectionURI)

	code := m.Run()

	os.Exit(code)
}

func TestWatchlistCreateHandler(t *testing.T) {
	t.Run("sends 201Created with data saved to db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		stockClient := mocks.NewMockstockClient(ctrl)
		userprofileClient := mocks.NewMockuserprofileClient(ctrl)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistCreateHandler(router, wlC, func(r *http.Request) string { return "userId" })

		stockClient.EXPECT().RegisterStock("INTC").Return(nil)
		watchlistRequest := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC"}, UserID: "userId"}

		body, _ := json.Marshal(watchlistRequest)

		req := httptest.NewRequest(http.MethodPost, "/watchlist", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected [%d], got [%d]", http.StatusCreated, res.StatusCode)
		}

		var result model.Watchlist
		json.NewDecoder(res.Body).Decode(&result)

		if result.UserID != "userId" {
			t.Fatalf("expected [%s], got [%s]", "userId", result.UserID)
		}

		if result.Name != "name" {
			t.Fatalf("expected [%s], got [%s]", "name", result.Name)
		}

		if result.Stocks[0] != "INTC" {
			t.Fatalf("expected [%s], got [%+v]", "[\"INTC\"]", result.Stocks)
		}
	})
	t.Run("saves watchlist to db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		stockClient := mocks.NewMockstockClient(ctrl)
		userprofileClient := mocks.NewMockuserprofileClient(ctrl)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistCreateHandler(router, wlC, func(r *http.Request) string { return "userId" })

		stockClient.EXPECT().RegisterStock("INTC").Return(nil)
		watchlistRequest := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC"}, UserID: "userId"}

		body, _ := json.Marshal(watchlistRequest)

		req := httptest.NewRequest(http.MethodPost, "/watchlist", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		var result model.Watchlist
		json.NewDecoder(res.Body).Decode(&result)

		savedObject, err := wlDb.Get(result.ID)

		if err != nil {
			t.Fatal("watchlist not found in Db ", err)
		}

		if savedObject.UserID != "userId" {
			t.Fatalf("expected [%s], got [%s]", "userId", result.UserID)
		}

		if savedObject.Name != "name" {
			t.Fatalf("expected [%s], got [%s]", "name", result.Name)
		}

		if savedObject.Stocks[0] != "INTC" {
			t.Fatalf("expected [%s], got [%+v]", "[\"INTC\"]", result.Stocks)
		}
	})
}

func TestWatchlistDeleteHandler(t *testing.T) {
	t.Run("deletes watchlist from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		watchlistRequest := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC"}, UserID: "userId"}
		watchlistID, _ := wlDb.Create(watchlistRequest)

		stockClient := mocks.NewMockstockClient(ctrl)
		userprofileClient := mocks.NewMockuserprofileClient(ctrl)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistDeleteHandler(router, wlC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodDelete, "/watchlist/"+watchlistID.Hex(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatalf("expected [%d], got [%d]", http.StatusNoContent, res.StatusCode)
		}

		_, err := wlDb.Get(watchlistID)
		if err.Error() != "mongo: no documents in result" {
			t.Fatal(err)
		}
	})
}

func TestWatchlistGetAllHandler(t *testing.T) {
	t.Run("returns the watchlists of the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		watchlistRequest1 := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC"}, UserID: "userId"}
		watchlistRequest2 := model.WatchlistRequest{Name: "name2", Stocks: []string{"XOM"}, UserID: "userId"}
		watchlistRequest3 := model.WatchlistRequest{Name: "name3", Stocks: []string{"INTC"}, UserID: "userId2"}
		wlDb.Create(watchlistRequest1)
		wlDb.Create(watchlistRequest2)
		wlDb.Create(watchlistRequest3)

		stockClient := mocks.NewMockstockClient(ctrl)
		userprofileClient := mocks.NewMockuserprofileClient(ctrl)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistGetAllHandler(router, wlC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodGet, "/watchlist", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		var result []model.Watchlist
		json.NewDecoder(res.Body).Decode(&result)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		if len(result) != 2 {
			t.Fatalf("expected array of 2, got [%+v]", result)
		}

		if result[0].UserID != "userId" || result[1].UserID != "userId" {
			t.Fatalf("expected watchlists for userID: userId, got [%+v]", result)
		}
	})
}

func TestWatchlistGetHandler(t *testing.T) {
	t.Run("returns the given watchlist of the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		watchlistRequest1 := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC"}, UserID: "userId"}
		watchlistRequest2 := model.WatchlistRequest{Name: "name2", Stocks: []string{"XOM"}, UserID: "userId"}
		watchlistRequest3 := model.WatchlistRequest{Name: "name3", Stocks: []string{"INTC"}, UserID: "userId2"}
		wlDb.Create(watchlistRequest1)
		watchlistID2, _ := wlDb.Create(watchlistRequest2)
		wlDb.Create(watchlistRequest3)

		stockClient := mocks.NewMockstockClient(ctrl)
		userprofileClient := mocks.NewMockuserprofileClient(ctrl)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistGetHandler(router, wlC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodGet, "/watchlist/"+watchlistID2.Hex(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		var result model.Watchlist
		json.NewDecoder(res.Body).Decode(&result)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		if result.UserID != "userId" {
			t.Fatalf("expected watchlists for userID: userId, got [%v]", result)
		}

		if result.ID != watchlistID2 {
			t.Fatalf("expected watchlist with ID: [%v], got [%v]", watchlistID2, result)
		}

		if len(result.Stocks) != 1 || result.Stocks[0] != "XOM" {
			t.Fatalf("expected watchlist with stocks: [\"XOM\"], got [%+v]", result.Stocks)
		}
	})
}

func TestWatchlistGetCalculatedHandler(t *testing.T) {
	t.Run("returns the given calculated watchlist of the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wlDb := database.NewWatchlists(db)
		watchlistRequest1 := model.WatchlistRequest{Name: "name", Stocks: []string{"INTC", "XOM"}, UserID: "userId"}
		watchlistID1, _ := wlDb.Create(watchlistRequest1)

		stockINTC := model.StockData{}
		stockINTC.Ticker = "INTC"
		stockINTC.Dividend = 0.33
		stockINTC.Eps = 5.43
		stockINTC.Price = 49.28
		stockINTC.DividendYield5yr.Avg = 2.62
		stockINTC.DividendYield5yr.Max = 3.65
		stockINTC.PeRatio5yr.Avg = 14.89
		stockINTC.PeRatio5yr.Min = 8.79

		expectedResultINTC := model.CalculatedStockInfo{}
		expectedResultINTC.Ticker = stockINTC.Ticker
		expectedResultINTC.AnnualDividend = 1.32
		expectedResultINTC.CurrentPe = 9.075506445672191
		expectedResultINTC.OptInPe = 11.84
		expectedResultINTC.PeColor = "green"
		expectedResultINTC.Price = stockINTC.Price
		expectedResultINTC.OptInPrice = 37.714285714285715
		expectedResultINTC.PriceColor = "red"
		expectedResultINTC.DividendYield = 2.678571428571429
		expectedResultINTC.OptInYield = 3.5
		expectedResultINTC.DividendColor = "yellow"

		stockClient := mocks.NewMockstockClient(ctrl)

		stockClient.EXPECT().Get(gomock.Any()).Return(stockINTC, nil).Times(2)

		userprofileClient := mocks.NewMockuserprofileClient(ctrl)

		expectedReturn := 9.0
		expectedRaise := 5.5
		defaultExpectation := 5.5
		userprofile := userprofileModel.Userprofile{Email: "alice@example.com", ExpectedReturn: &expectedReturn, Expectations: []userprofileModel.Expectation{userprofileModel.Expectation{Stock: "INTC", ExpectedRaise: &expectedRaise}}, DefaultExpectation: &defaultExpectation}

		userprofileClient.EXPECT().GetUserprofile("userId").Return(userprofile, nil)
		sp500Client := mockSp500Client{}
		stockService := service.NewStockService(&sp500Client)
		wlC := controllers.NewWatchlistController(wlDb, stockClient, userprofileClient, stockService)

		router := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
		handlers.WatchlistGetCalculatedHandler(router, wlC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodGet, "/watchlist/"+watchlistID1.Hex()+"/calculated", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		var result []model.CalculatedStockInfo
		json.NewDecoder(res.Body).Decode(&result)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		if result[0] != expectedResultINTC {
			t.Fatalf("expected [%v], got [%v]", expectedResultINTC, result)
		}

		if result[1] != expectedResultINTC {
			t.Fatalf("expected [%v], got [%v]", expectedResultINTC, result)
		}
	})
}
