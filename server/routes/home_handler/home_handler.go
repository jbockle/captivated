package home_handler

import (
	"fmt"
	"net/http"
)

func HomeHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "captivated - github webhooks outbox producer")
}
