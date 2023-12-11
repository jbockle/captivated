package events_handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbockle/captivated/server/services"
)

func EventsHandler(response http.ResponseWriter, request *http.Request) {
	eventId := mux.Vars(request)["eventId"]

	err := services.Storage.Stream(request.Context(), eventId, response)
	if errors.Is(err, services.ErrFileNotFound) {
		http.Error(response, "Not found, it may have expired", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
}
