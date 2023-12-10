package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbockle/captivated/server/routes"
)

func Serve() error {
	router := mux.NewRouter()
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

	fmt.Println("listening http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	return err
}
