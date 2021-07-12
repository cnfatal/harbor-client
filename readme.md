# harbor-client

[![Open in Visual Studio Code](https://open.vscode.dev/badges/open-in-vscode.svg)](https://open.vscode.dev/fatalc/harbor-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatalc/harbor-client)](https://goreportcard.com/report/github.com/fatalc/harbor-client)
[![Test Coverage](https://api.codeclimate.com/v1/badges/7b80205b09f94bfb469a/test_coverage)](https://codeclimate.com/github/fatalc/harbor-client/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/7b80205b09f94bfb469a/maintainability)](https://codeclimate.com/github/fatalc/harbor-client/maintainability)

Golang client of [goharbor/harbor](https://github.com/goharbor/harbor).

**NOT STABLE YET**

## Features

- OCI Distribution Client Supported,see [oci.go](oci.go).
- Light && Simple
- Avoid import additional libraries from harbor, like beego etc.
- Compatible with harbor v2
- auto

## Install

```sh
go get github.com/fatalc/harbor-client
```

## Example

harbor client:

```go
import client "github.com/fatalc/harbor-client"

cli, _ := client.NewClient("harbor.example.com", client.WithBasicAuth("admin", "password"))

image := "harbor.example.com/library/nginx:alpine"
project, repository, reference, _ := client.ParseHarborSuitImage(image)

ctx:=context.Background()
artifact, _ := cli.GetArtifact(ctx, project, repository, reference, client.GetArtifactOptions{})

fmt.Println(artifact)
```

> Ommited error handle

OCI distribution client:

```go
ocicli, _ := client.NewOCIDistributionClient("registry.example.com", client.BasicAuth("user", "password"))

if err := ocicli.Ping(context.Background()); err != nil {
    log.Infof("may be auth failed")
    return
}

tags, _ := ocicli.ListTags(context.Background(), "library/nginx")
fmt.Printf("tags: %s", tags.Tags)
```

## Documents

See [Go Doc](https://pkg.go.dev/github.com/fatalc/harbor-client)

## Contributing

Everyone is welcome to contribute.no limit, just creeate a Merge Request.

See [Project](https://github.com/fatalc/harbor-client/projects/1) for more information.
