package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	stockHttp "github.com/nagymarci/stock-commons/http"
	"github.com/nagymarci/stock-commons/reqid"
	"github.com/nagymarci/stock-watchlist/controllers"
	"github.com/nagymarci/stock-watchlist/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WatchlistCreateHandler(mux *mux.Router, watchlist *controllers.WatchlistController, extractUserID func(*http.Request) string) {
	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)

		log := logrus.WithFields(logrus.Fields("userId": userID, "requestId": GetRequestId(r)))

		var watchlistRequest *model.WatchlistRequest

		err := json.NewDecoder(r.Body).Decode(&watchlistRequest)

		if err != nil {
			message := "Failed to deserialize payload: " + err.Error()
			stockHttp.HandleErrorResponse(message, w, http.StatusBadRequest)
			log.Errorln(message)
			return
		}

		if watchlistRequest.Stocks == nil || len(watchlistRequest.Stocks) < 1 {
			message := "Required value 'stocks' is missing"
			stockHttp.HandleErrorResponse(message, w, http.StatusBadRequest)
			log.Errorln(message)
			return
		}

		if len(watchlistRequest.Name) < 1 || watchlistRequest.Name == " " {
			message := "Required value 'name' is missing"
			stockHttp.HandleErrorResponse(message, w, http.StatusBadRequest)
			log.Errorln(message)
			return
		}

		watchlistRequest.UserID = userID

		result, err := watchlist.Create(log, watchlistRequest)

		if err != nil {
			message := "Watchlist creation failed: " + err.Error()
			stockHttp.HandleErrorResponse(message, w, http.StatusInternalServerError)
			log.Errorln(message)
			return
		}

		stockHttp.HandleJSONResponse(result, w, http.StatusCreated)

	}).Methods(http.MethodPost, http.MethodOptions)
}

func WatchlistDeleteHandler(router *mux.Router, watchlist *controllers.WatchlistController, extractUserID func(*http.Request) string) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		watchlistID, err := extractWatchlistID(r)

		log := logrus.WithFields(logrus.Fields("userId": userID, "requestId": GetRequestId(r), "watchlistId": watchlistID))

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		err = watchlist.Delete(log, id, userID)

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete, http.MethodOptions)
}

func WatchlistGetAllHandler(router *mux.Router, watchlist *controllers.WatchlistController, extractUserID func(*http.Request) string) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)

		log := logrus.WithFields(logrus.Fields("userId": userID, "requestId": GetRequestId(r)))

		result, err := watchlist.GetAll(log, userID)

		if err != nil {
			log.Errorln(err)
			handleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func WatchlistGetHandler(router *mux.Router, watchlist *controllers.WatchlistController, extractUserID func(*http.Request) string) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		watchlistID, err := extractWatchlistID(r)

		log := logrus.WithFields(logrus.Fields("userId": userID, "requestId": GetRequestId(r), "watchlistId": watchlistID))

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		result, err := watchlist.Get(log, id, userID)

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func WatchlistGetCalculatedHandler(router *mux.Router, watchlist *controllers.WatchlistController, extractUserID func(*http.Request) string) {
	router.HandleFunc("/{id}/calculated", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		watchlistID, err := extractWatchlistID(r)

		log := logrus.WithFields(logrus.Fields("userId": userID, "requestId": GetRequestId(r), "watchlistId": watchlistID))

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		result, err := watchlist.GetCalculated(log, id, userID)

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func extractWatchlistID(r *http.Request) (primitive.ObjectID, error) {
	id := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		message := "Invalid watchlist id: " + err.Error()
		return primitive.NilObjectID, model.NewBadRequestError(message)
	}

	return objectID, nil
}
