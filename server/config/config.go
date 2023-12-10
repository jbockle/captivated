package config

import (
	"fmt"
	"log"
	"os"
)

var (
	env_acs_key = "ASB_CONNECTION_STRING"
	env_bcs_key = "BLOB_CONNECTION_STRING"
	env_ws_key  = "WEBHOOK_SECRET"

	webhookSecretRaw     string = os.Getenv(env_ws_key)
	WebhookSecret        []byte = []byte(webhookSecretRaw)
	AsbConnectionString  string = os.Getenv(env_acs_key)
	AsbEntity            string = os.Getenv("ASB_ENTITY")
	BlobConnectionString string = os.Getenv(env_bcs_key)
	BlobContainer        string = os.Getenv("BLOB_CONTAINER")
)

func Init() {
	log.Println("Initializing config")

	assertNotEmpty(AsbConnectionString, env_acs_key)
	assertNotEmpty(BlobConnectionString, env_bcs_key)

	if len(AsbEntity) == 0 {
		AsbEntity = "events"
	}

	if len(BlobContainer) == 0 {
		BlobContainer = "events"
	}
}

func assertNotEmpty(item string, envKey string) {
	if len(item) == 0 {
		log.Fatalln(fmt.Sprintf("Environment variable '%v' is undefined", envKey))
	}
}
