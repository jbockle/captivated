package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v57/github"
)

// represents a github webhook event
type GitHubEvent struct {
	// the unique identifier for this event
	Id string `json:"id"`
	// the name of the event
	Type string `json:"type"`
	// the GUID to identify this delivery
	DeliveryId string `json:"deliveryId"`
	// the type of resource where the webhook was created
	TargetType string `json:"targetType"`
	// the unique identifier of the resource where the webhook was created
	TargetId string `json:"targetId"`
	// the event data
	Data interface{} `json:"data,omitempty"`
	// indicates this event was added to blob storage as its size was too large to include as an event, event.data is nil
	IsReference bool `json:"isReference,omitempty"` // TODO should change this to an address
}

type ValidationFailed struct {
	error
	errors []string
}

// decodes the request body from a GitHubEvent type
func GetGitHubEvent(request *http.Request, payload []byte) (event GitHubEvent, err error) {
	event.Id = request.Header.Get("X-GitHub-Hook-ID")
	event.Type = request.Header.Get("X-GitHub-Event")
	event.DeliveryId = request.Header.Get("X-GitHub-Delivery")
	event.TargetType = request.Header.Get("X-GitHub-Hook-Installation-Target-Type")
	event.TargetId = request.Header.Get("X-GitHub-Hook-Installation-Target-ID")

	if err = validate(&event); err != nil {
		return
	}

	if event.Data, err = github.ParseWebHook(github.WebHookType(request), payload); err != nil {
		err = fmt.Errorf("Error parsing payload %v %v %w", event.Type, event.Id, err)
		return
	}

	return
}

func validate(event *GitHubEvent) error {
	validationFailure := new(ValidationFailed)

	if len(event.Id) == 0 {
		validationFailure.errors = append(validationFailure.errors, "event has an empty id")
	}

	if len(event.Type) == 0 {
		validationFailure.errors = append(validationFailure.errors, "event has an empty type")
	}

	if len(validationFailure.errors) != 0 {
		return validationFailure
	}

	return nil
}

func (event *GitHubEvent) ToBytes() ([]byte, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (event GitHubEvent) ToReference() *GitHubEvent {
	reference := event
	reference.Data = nil
	reference.IsReference = true

	return &reference
}
