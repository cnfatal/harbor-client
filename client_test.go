package client_test

import (
	"context"
	"fmt"
	"log"

	client "github.com/fatalc/harbor-client"
)

func ExampleNewClient() {
	cli, err := client.NewClient("harbor.example.com", client.WithBasicAuth("admin", "password"))
	if err != nil {
		log.Fatal(err)
	}
	sysinfo, err := cli.SystemInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sysinfo)

	image := "harbor.example.com/library/nginx:alpine"
	project, repository, reference, err := client.ParseHarborSuitImage(image)
	if err != nil {
		log.Fatal(err)
	}
	artifact, err := cli.GetArtifact(context.Background(), project, repository, reference, client.GetArtifactOptions{
		WithScanOverview: true,
		WithLabel:        true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(artifact)
}
