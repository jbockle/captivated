package services

import (
	"context"
	"errors"

	"github.com/jbockle/captivated/server/models"
)

type PublisherImpl struct{}

func (publisher *PublisherImpl) Publish(ctx context.Context, event *models.GitHubEvent) (err error) {
	err = Broker.Send(ctx, event)
	if errors.Is(err, ErrMsgTooLarge) {
		if err = publishAsReference(ctx, event); err != nil {
			return
		}
	}

	return
}

func publishAsReference(ctx context.Context, event *models.GitHubEvent) (err error) {
	referenceEvent := event.ToReference()

	if err = Storage.Save(ctx, referenceEvent); err != nil {
		return
	}

	err = Broker.Send(ctx, referenceEvent)

	return
}
