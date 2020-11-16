package routes

import (
	"net/http"
	"os"

	"github.com/urfave/negroni"

	"github.com/nagymarci/stock-commons/authorization"
	"github.com/nagymarci/stock-watchlist/handlers"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-watchlist/controllers"
)

func Route(watchlistController *controllers.WatchlistController) http.Handler {
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	watchlist := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
	handlers.WatchlistCreateHandler(watchlist, watchlistController, authorization.DefaultExtractUserID)
	handlers.WatchlistDeleteHandler(watchlist, watchlistController, authorization.DefaultExtractUserID)
	handlers.WatchlistGetAllHandler(watchlist, watchlistController, authorization.DefaultExtractUserID)
	handlers.WatchlistGetHandler(watchlist, watchlistController, authorization.DefaultExtractUserID)
	handlers.WatchlistGetCalculatedHandler(watchlist, watchlistController, authorization.DefaultExtractUserID)

	audience := os.Getenv("WATCHLIST_AUDIENCE")
	authServer := os.Getenv("AUTHORIZATION_SERVER")
	watchlistScope := os.Getenv("WATCHLIST_SCOPE")

	auth := negroni.New(
		negroni.HandlerFunc(authorization.CreateAuthorizationMiddleware(audience, authServer).HandlerWithNext),
		negroni.HandlerFunc(authorization.CreateScopeMiddleware(watchlistScope, authServer, audience)))

	router.PathPrefix("/watchlist").Handler(auth.With(negroni.Wrap(watchlist)))

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)

	return n
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
