package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	distributionspecsv1 "github.com/opencontainers/distribution-spec/specs-go/v1"
	imagespecv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type OCIDistributionClient struct {
	Server string
	Auth   Auth
}

// OCI Distribution Specification Client
//
// For more information visit below URL
// https://github.com/opencontainers/distribution-spec/blob/main/spec.md#endpoints
func NewOCIDistributionClient(server string, auth Auth) (*OCIDistributionClient, error) {
	return &OCIDistributionClient{Server: server, Auth: auth}, nil
}

// see: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#determining-support
// To check whether or not the registry implements this specification,
// perform a GET request to the following endpoint: /v2/ end-1.
// If the response is 200 OK, then the registry implements this specification.
// It can used to detectd server connection and auth too.
// end-1	GET	/v2/	200	404/401
func (c *OCIDistributionClient) Ping(ctx context.Context) error {
	return c.request(ctx, http.MethodGet, "/v2", nil, nil)
}

// end-3 	GET/HEAD /v2/<name>/manifests/<reference>
func (c *OCIDistributionClient) GetManifest(ctx context.Context, name, reference string) (*imagespecv1.Manifest, error) {
	manifest := &imagespecv1.Manifest{}
	err := c.request(ctx, http.MethodGet, "/v2/"+name+"/manifests/"+reference, nil, manifest)
	return manifest, err
}

// end-8a	GET	/v2/<name>/tags/list
func (c *OCIDistributionClient) ListTags(ctx context.Context, name string) (*distributionspecsv1.TagList, error) {
	tags := &distributionspecsv1.TagList{}
	err := c.request(ctx, http.MethodGet, "/v2/"+name+"/tags/list", nil, tags)
	return tags, err
}

// end-8b	GET	/v2/<name>/tags/list?n=<integer>&last=<integer>
func (c *OCIDistributionClient) ListTagsPaged(ctx context.Context, name string, n, last int) (*distributionspecsv1.TagList, error) {
	path := fmt.Sprintf("/v2/%s/tags/list?n=%d&last=%d", name, n, last)
	tags := &distributionspecsv1.TagList{}
	err := c.request(ctx, http.MethodGet, path, nil, tags)
	return tags, err
}

// end-9	DELETE	/v2/<name>/manifests/<reference>
func (c *OCIDistributionClient) DeleteManifest(ctx context.Context, name, reference string) error {
	return c.request(ctx, http.MethodDelete, "/v2/"+name+"/manifests/"+reference, nil, nil)
}

func (c *OCIDistributionClient) request(ctx context.Context, method string, path string, postbody interface{}, into interface{}) error {
	var body io.Reader
	switch typed := postbody.(type) {
	// convert to bytes
	case []byte:
		body = bytes.NewBuffer(typed)
	// thise type can processed by 'http.NewRequestWithContext(...)'
	case io.Reader:
		body = typed
	case nil:
		// do nothing
	// send json format
	default:
		bts, err := json.Marshal(postbody)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(bts)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.Server+path, body)
	if err != nil {
		return err
	}
	if c.Auth != nil {
		c.Auth(req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errresp := &distributionspecsv1.ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(errresp); err != nil {
			return err
		}
		return errresp
	}
	return json.NewDecoder(resp.Body).Decode(into)
}

// https://github.com/opencontainers/distribution-spec/blob/main/spec.md#error-codes
//
// code-1	BLOB_UNKNOWN			blob unknown to registry
// code-2	BLOB_UPLOAD_INVALID		blob upload invalid
// code-3	BLOB_UPLOAD_UNKNOWN		blob upload unknown to registry
// code-4	DIGEST_INVALID			provided digest did not match uploaded content
// code-5	MANIFEST_BLOB_UNKNOWN	blob unknown to registry
// code-6	MANIFEST_INVALID		manifest invalid
// code-7	MANIFEST_UNKNOWN		manifest unknown
// code-8	NAME_INVALID			invalid repository name
// code-9	NAME_UNKNOWN			repository name not known to registry
// code-10	SIZE_INVALID			provided length did not match content length
// code-12	UNAUTHORIZED			authentication required
// code-13	DENIED					requested access to the resource is denied
// code-14	UNSUPPORTED				the operation is unsupported
// code-15	TOOMANYREQUESTS			too many requests
