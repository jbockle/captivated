package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbockle/captivated/server/routes"
)

func Serve() {
	router := mux.NewRouter()
	router.
		Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						http.Error(w, fmt.Sprintf("An error occurred: %+v", err), http.StatusInternalServerError)
					}
				}()

				slog.Debug(
					"request received",
					"path", r.URL.Path,
					"method", r.Method,
					"length", r.ContentLength,
				)

				next.ServeHTTP(w, r)
			})
		})
	router.
		HandleFunc("/", routes.HomeHandler).
		Methods(http.MethodGet)
	router.
		HandleFunc("/events", routes.NewEventHandler).
		Methods(http.MethodPost)
	router.
		HandleFunc("/events/{eventId}", routes.EventsHandler).
		Methods(http.MethodGet)

	http.Handle("/", router)

	slog.Info("started http server", "address", "http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
