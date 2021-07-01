# harbor-client

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

```

OCI distribution client:

```go
ocicli, err := client.NewOCIDistributionClient("registry.example.com", client.BasicAuth("user", "password"))
if err != nil {
    log.Fatal(err)
}

if err := ocicli.Ping(context.Background()); err != nil {
    return
}

tags, err := ocicli.ListTags(context.Background(), "library/nginx")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("tags: %s", tags.Tags)
```

## Documents

See [Go Doc](https://pkg.go.dev/github.com/fatalc/harbor-client)
