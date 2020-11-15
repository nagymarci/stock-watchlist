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
		wlC := controllers.NewWatchlistController(wlDb, stockClient)

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
		wlC := controllers.NewWatchlistController(wlDb, stockClient)

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
		wlC := controllers.NewWatchlistController(wlDb, stockClient)

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

}
