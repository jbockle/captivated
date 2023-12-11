package services

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/jbockle/captivated/server/models"
)

var (
	Broker    BrokerService
	Storage   StorageService
	Publisher PublisherService

	ttlDays         int           = 7
	ttl             time.Duration = 24 * time.Hour * time.Duration(ttlDays)
	contentType     string        = "application/json"
	ErrMsgTooLarge  error         = errors.New("Message too large to publish to broker")
	ErrFileNotFound error         = errors.New("File not found")
)

func Init() {
	log.Println("Initializing services")
	Broker = CreateBroker()
	Storage = CreateStorage()
	Publisher = &PublisherImpl{}
}

type PublisherService interface {
	Publish(ctx context.Context, event *models.GitHubEvent) error
}

type BrokerService interface {
	Send(ctx context.Context, event *models.GitHubEvent) error
}

type StorageService interface {
	Save(ctx context.Context, event *models.GitHubEvent) error

	Stream(ctx context.Context, eventId string, to io.Writer) error

	DeleteExpired() error
}

func checkExpired(age time.Duration) bool {
	return age > ttl
}

func StartDeleteExpiredTask() {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Starting task Storage.DeleteExpired")
				if err := Storage.DeleteExpired(); err != nil {
					log.Println("Task Storage.DeleteExpired failed:", err)
				}
			}
		}
	}()
}
