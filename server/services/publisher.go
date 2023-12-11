package services

import (
	"context"
	"errors"

	"github.com/jbockle/captivated/server/models"
)

type PublisherImpl struct{}

func (publisher *PublisherImpl) Publish(ctx context.Context, event *models.GitHubEvent) (err error) {
	if err = Storage.Save(ctx, event); err != nil {
		return
	}

	err = Broker.Send(ctx, event)
	if errors.Is(err, ErrMsgTooLarge) {
		if err = publishAsReference(ctx, event); err != nil {
			return
		}
	}

	return
}

func publishAsReference(ctx context.Context, event *models.GitHubEvent) error {
	referenceEvent := event.ToReference()

	return Broker.Send(ctx, referenceEvent)
}
