package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	stockHttp "github.com/nagymarci/stock-commons/http"
	"github.com/nagymarci/stock-watchlist/controllers"
	"github.com/urfave/negroni"

	"github.com/sirupsen/logrus"
)

func StockGetAllCalculatedHandler(mux *mux.Router, stockController *controllers.StockController) {
	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		log := logrus.WithField("userId", "")

		result, err := stockController.GetAllCalculated(log, "")

		if err != nil {
			stockHttp.HandleError(err, w)
			log.Errorln(err)
			return
		}

		stockHttp.HandleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func StockGetAllCalculatedForUserHandler(router *mux.Router, auth *negroni.Negroni, stockController *controllers.StockController, extractUserIDFromToken func(*http.Request) string) {
	router.Handle("/{userId}", auth.With(negroni.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserIDFromToken(r)
		id := mux.Vars(r)["userId"]

		if userID != id {
			message := "UserID in request doesn't match userID in token"
			stockHttp.HandleErrorResponse(message, w, http.StatusUnauthorized)
			logrus.WithFields(logrus.Fields{"userId": userID, "request_userId": id}).Error("Unauthorized")
			return
		}

		log := logrus.WithField("userId", userID)

		result, err := stockController.GetAllCalculated(log, userID)

		if err != nil {
			stockHttp.HandleError(err, w)
			log.Errorln(err)
			return
		}

		stockHttp.HandleJSONResponse(result, w, http.StatusOK)
	}))).Methods(http.MethodGet)
}
