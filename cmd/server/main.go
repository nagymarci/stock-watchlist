package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nagymarci/stock-watchlist/api"
	"github.com/nagymarci/stock-watchlist/controllers"
	"github.com/nagymarci/stock-watchlist/routes"
	"github.com/nagymarci/stock-watchlist/service"
	"github.com/robfig/cron/v3"

	"github.com/nagymarci/stock-watchlist/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	db := database.New(os.Getenv("DB_CONNECTION_URI"))
	rDb := database.NewRecommendations(db)
	wDb := database.NewWatchlists(db)

	sC := api.NewStockClient(os.Getenv("STOCK_SCREENER_URL"))
	upC := api.NewUserprofileClient(os.Getenv("USERPROFILE_URL"))

	sS := service.NewStockService(sC)

	wC := controllers.NewWatchlistController(wDb, sC, upC, sS)

	router := routes.Route(wC)

	mC := service.NewMail()
	c := cron.New()
	n := service.NewNotifier(rDb, wDb, sC, sS, upC, mC)
	_, err := c.AddFunc("CRON_TZ=America/New_York 0 8-18 * * MON-FRI", n.NotifyChanges)
	if err != nil {
		log.Errorln(err)
	}

	c.Start()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
