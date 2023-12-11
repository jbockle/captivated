package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/models"
)

type BlobStorage struct {
	container *container.Client
}

func (storage *BlobStorage) Save(ctx context.Context, event *models.GitHubEvent) (err error) {
	bytes, err := event.ToBytes()
	if err != nil {
		log.Println("Error converting event to bytes:", event, err)
		return
	}

	blobClient := storage.getBlobClient(event.Id)

	if _, err = blobClient.UploadBuffer(ctx, bytes, &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{BlobContentType: &contentType},
	}); err != nil {
		log.Println("Error uploading blob:", err)
	}

	return
}

func (storage *BlobStorage) Stream(ctx context.Context, eventId string, to io.Writer) (err error) {
	blobClient := storage.getBlobClient(eventId)

	response, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		if stgErr, ok := err.(*azcore.ResponseError); ok && stgErr.StatusCode == http.StatusNotFound {
			err = ErrFileNotFound
		}

		return
	}

	body := response.Body
	_, err = io.Copy(to, body)

	return
}

func (storage *BlobStorage) DeleteExpired() error {
	now := time.Now().UTC()
	maxResults := int32(100)

	// TODO use container.FilterBlobs
	pager := storage.container.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include:    container.ListBlobsInclude{},
		MaxResults: &maxResults,
	})

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			return err
		}

		for _, blob := range response.Segment.BlobItems {
			if isExpired(blob.Properties.CreationTime, &now) {
				_, err = storage.container.
					NewBlobClient(*blob.Name).
					Delete(context.Background(), nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (storage *BlobStorage) getBlobClient(eventId string) *blockblob.Client {
	return storage.container.NewBlockBlobClient(eventId + ".json")
}

func CreateStorage() StorageService {
	blobclient, err := azblob.NewClientFromConnectionString(config.BlobConnectionString, nil)
	if err != nil {
		panic(err)
	}

	container := blobclient.ServiceClient().NewContainerClient(config.BlobContainer)

	return &BlobStorage{container}
}
