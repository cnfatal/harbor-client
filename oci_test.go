package client_test

import (
	"context"
	"fmt"
	"log"

	client "github.com/cnfatal/harbor-client"
)

func ExampleNewOCIDistributionClient() {
	ocicli, err := client.NewOCIDistributionClient("registry.example.com", client.BasicAuth("user", "password"))
	if err != nil {
		log.Fatal(err)
	}

	if err := ocicli.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	tags, err := ocicli.ListTags(context.Background(), "library/nginx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tags: %s", tags.Tags)
}
