package services

import (
	"context"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/models"
)

type ServiceBusBroker struct {
	sender *azservicebus.Sender
}

func (broker *ServiceBusBroker) Send(ctx context.Context, event *models.GitHubEvent) (err error) {
	body, err := event.ToBytes()
	if err != nil {
		log.Println("Error converting event to bytes:", event, err)
		return
	}

	contentType := "application/json"

	message := &azservicebus.Message{
		MessageID:     &event.Id,
		Body:          body,
		ContentType:   &contentType,
		CorrelationID: &event.DeliveryId,
		Subject:       &event.Type,
		TimeToLive:    &ttl,
	}

	err = broker.sender.SendMessage(ctx, message, nil)
	if errors.Is(err, azservicebus.ErrMessageTooLarge) {
		return ErrMsgTooLarge
	}
	return
}

func CreateBroker() BrokerService {
	asbclient, err := azservicebus.NewClientFromConnectionString(config.AsbConnectionString, nil)
	if err != nil {
		panic(err)
	}

	sender, err := asbclient.NewSender(config.AsbEntity, nil)
	if err != nil {
		panic(err)
	}

	return &ServiceBusBroker{sender}
}
