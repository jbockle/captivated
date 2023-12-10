package publisher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/models"
)

type publisherClients struct {
	message   *azservicebus.Sender
	container *container.Client
}

var (
	clients     publisherClients = publisherClients{}
	ttl         time.Duration    = 24 * time.Hour * 7 // 7 days
	contentType string           = "application/json"
)

func Init() {
	log.Println("Initializing publisher")

	asbclient, err := azservicebus.NewClientFromConnectionString(config.AsbConnectionString, nil)
	if err != nil {
		panic(err)
	}

	value, err := asbclient.NewSender(config.AsbEntity, nil)
	if err != nil {
		panic(err)
	}

	blobclient, err := azblob.NewClientFromConnectionString(config.BlobConnectionString, nil)
	if err != nil {
		panic(err)
	}

	clients.message = value
	clients.container = blobclient.ServiceClient().NewContainerClient(config.BlobContainer)

	go startStorageCleaner()
}

func Publish(context context.Context, event *models.GitHubEvent) {
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

	err = clients.message.SendMessage(context, message, nil)
	if err != nil {
		log.Println("error sending message:", event.Id, err)
		if errors.Is(err, azservicebus.ErrMessageTooLarge) {
			publishToStorage(context, event)
		}
	}
}

func startStorageCleaner() {
	ticker := time.NewTicker(24 * time.Hour)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now().UTC()
				maxResults := int32(100)
				pager := clients.container.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
					Include:    container.ListBlobsInclude{},
					MaxResults: &maxResults,
				})

				for pager.More() {
					response, err := pager.NextPage(context.Background())
					if err != nil {
						log.Println("error paging blobs:", err)
						break
					}

					for _, blob := range response.Segment.BlobItems {
						if isMoreThan7DaysAgo(blob.Properties.CreationTime, &now) {
							_, err = clients.container.NewBlobClient(*blob.Name).Delete(context.Background(), nil)
							if err != nil {
								log.Println("error deleting expired blob:", blob.Name, err)
								continue
							}
						}
					}
				}

			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
}

func publishToStorage(context context.Context, event *models.GitHubEvent) {
	bytes, err := event.ToBytes()
	if err != nil {
		log.Println("Error converting event to bytes:", event, err)
		return
	}

	blobClient := clients.container.NewBlockBlobClient(fmt.Sprintf("%v.json", event.Id))
	_, err = blobClient.UploadBuffer(context, bytes, &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
	})
	if err != nil {
		log.Println("Error uploading blob:", err)
	}

	ref := event.ToReference()

	body, err := ref.ToBytes()
	if err != nil {
		log.Println("Error converting event to bytes:", ref, err)
		return
	}

	contentType := "application/json"

	message := &azservicebus.Message{
		MessageID:     &ref.Id,
		Body:          body,
		ContentType:   &contentType,
		CorrelationID: &ref.DeliveryId,
		Subject:       &ref.Type,
		TimeToLive:    &ttl,
	}

	err = clients.message.SendMessage(context, message, nil)
	if err != nil {
		log.Println("error sending message:", ref.Id, err)
	}
}

func isMoreThan7DaysAgo(t, from *time.Time) bool {
	return t.After(from.AddDate(0, 0, 7))
}

func GetPublishedBlobFromContainer(context context.Context, name string) (io.Reader, error) {
	blobClient := clients.container.NewBlobClient(name)

	response, err := blobClient.DownloadStream(context, nil)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}
