package new_event_handler

import (
	"log"
	"net/http"

	"github.com/google/go-github/v57/github"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/models"
	"github.com/jbockle/captivated/server/publisher"
)

func NewEventHandler(response http.ResponseWriter, request *http.Request) {
	eventType := request.Header.Get("X-GitHub-Event")
	hookId := request.Header.Get("X-GitHub-Hook-ID")

	log.Println("Validating payload", eventType, hookId)
	payload, err := github.ValidatePayload(request, config.WebhookSecret)
	if err != nil {
		log.Println("Error validating payload", eventType, hookId, err)
		return
	}

	log.Println("Getting payload event", eventType, hookId)
	event, err := models.GetGitHubEvent(request, payload)
	if err != nil {
		log.Println("Error getting payload", eventType, hookId, err)
		return
	}

	publisher.Publish(request.Context(), &event)
}
