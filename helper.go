package client

import (
	"errors"
	"strings"

	containerdreference "github.com/containerd/containerd/reference"
	"github.com/docker/distribution/reference"
)

var ErrNotHarborImage = errors.New("not a harbor suit image")

// https://github.com/containerd/containerd/blob/0396089f79f241df4d8724a0cd31cf58523756ff/reference/reference.go#L84
func SplitImageNameTag(image string) (string, string) {
	spec, err := containerdreference.Parse(image)
	if err != nil {
		// backoff
		spls := strings.Split(image, ":")
		if len(spls) > 1 {
			return spls[0], spls[1]
		}
		return spls[0], ""
	}
	return spec.Locator, spec.Object
}

// ParseHarborImage parse a image and return harbor project,repository,reference
func ParseHarborSuitImage(image string) (project, repository, reference string, err error) {
	_, path, name, tag, err := ParseImag(image)
	if err != nil {
		return "", "", "", err
	}
	if path == "" || name == "" || tag == "" {
		return "", "", "", ErrNotHarborImage
	}
	return path, name, tag, nil
}

// ParseImag
// barbor.foo.com/project/artifact:tag -> barbor.foo.com,project,artifact,tag
// barbor.foo.com/project/foo/artifact:tag -> barbor.foo.com,project,foo/artifact,tag
// barbor.foo.com/artifact:tag -> barbor.foo.com,library,artifact,tag
// project/artifact:tag -> docker.io,project,artifact,tag
func ParseImag(image string) (domain, path, name, tag string, err error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return
	}
	domain = reference.Domain(named)

	fullpath := reference.Path(named)
	splits := strings.SplitN(fullpath, "/", 2)
	if len(splits) > 1 {
		path = splits[0]
		name = splits[1]
	} else {
		path = "library"
		name = splits[0]
	}

	if tagged, ok := named.(reference.Tagged); ok {
		tag = tagged.Tag()
	}
	if tagged, ok := named.(reference.Digested); ok {
		tag = tagged.Digest().String()
	}
	if tag == "" {
		tag = "latest"
	}
	return
}
