package pkg

import (
	"log"
	"os"

	"github.com/henry11996/fugle-golang/fugle"
)

func InitFugle() fugle.Client {
	client, err := fugle.NewFugleClient(fugle.ClientOption{
		ApiToken: os.Getenv("FUGLE_API_TOKEN"),
		Version:  "v0.3",
	})
	if err != nil {
		log.Fatal("failed to init fugle api client, " + err.Error())
	}
	return client
}
