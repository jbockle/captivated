package events_handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/gorilla/mux"
	"github.com/jbockle/captivated/server/publisher"
)

func EventsHandler(response http.ResponseWriter, request *http.Request) {
	eventId := mux.Vars(request)["eventId"]

	reader, err := publisher.GetPublishedBlobFromContainer(request.Context(), eventId+".json")
	if err != nil {
		statusCode := http.StatusInternalServerError
		help := ""
		if stgErr, ok := err.(*azcore.ResponseError); ok && stgErr.StatusCode == http.StatusNotFound {
			statusCode = http.StatusNotFound
			help = ", it may have expired"
		}

		log.Println("Error retrieving event with id", eventId, err)

		http.Error(response, fmt.Sprintf("%v: Error retrieving event with id '%v'%v", statusCode, eventId, help), statusCode)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if _, err = io.Copy(response, reader); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

}
