package new_event_handler

import (
	"log/slog"
	"net/http"

	"github.com/google/go-github/v57/github"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/models"
	"github.com/jbockle/captivated/server/services"
)

func NewEventHandler(response http.ResponseWriter, request *http.Request) {
	eventType := request.Header.Get("X-GitHub-Event")
	hookId := request.Header.Get("X-GitHub-Hook-ID")

	slog.Debug("Validating payload", "type", eventType, "id", hookId)
	payload, err := github.ValidatePayload(request, config.WebhookSecret)
	if err != nil {
		slog.Error("Error validating payload", "type", eventType, "id", hookId, "err", err)
		return
	}

	slog.Debug("Getting payload event", "type", eventType, "id", hookId)
	event, err := models.GetGitHubEvent(request, payload)
	if err != nil {
		slog.Error("Error getting payload", "type", eventType, "id", hookId, "err", err)
		return
	}

	slog.Debug("Publishing payload event", "type", eventType, "id", hookId)
	err = services.Publisher.Publish(request.Context(), &event)
	if err != nil {
		slog.Error("Error publishing payload", "type", eventType, "id", hookId, "err", err)
	}
}
