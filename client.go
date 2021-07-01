package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	harborerrors "github.com/goharbor/harbor/src/lib/errors"
)

const (
	apiVerisonPrefix  = "/api/v2.0"
	csrfTokenHeader   = "X-Harbor-CSRF-Token"
	xTotalCountHeader = "X-Total-Count"
	linkHeader        = "Link"
)

func NewClient(addr string, options ...Option) (*Client, error) {
	cli := &Client{
		httpclient: &http.Client{},
		Server:     addr + apiVerisonPrefix,
	}
	for _, opt := range options {
		opt(cli)
	}
	return cli, nil
}

type Auth func(req *http.Request)

type Client struct {
	Server     string
	Auth       Auth
	httpclient *http.Client
	csrftoken  string
}

func WithBasicAuth(username, password string) Option {
	return func(opts *Client) { opts.Auth = BasicAuth(username, password) }
}

func BasicAuth(username, password string) Auth {
	return func(req *http.Request) { req.SetBasicAuth(username, password) }
}

func WithTokenAuth(token string) Option {
	return func(opts *Client) { opts.Auth = TokenAuth(token) }
}

func TokenAuth(token string) Auth {
	return func(req *http.Request) { req.Header.Add("Authorization", "Bearer "+token) }
}

type Option func(opts *Client)

func (c *Client) doRequest(ctx context.Context, method string, path string, data interface{}, decodeinto interface{}) error {
	_, err := c.doRequestWithResponse(ctx, method, path, data, decodeinto)
	return err
}

func (c *Client) doRequestWithResponse(ctx context.Context, method string, path string, data interface{}, decodeinto interface{}) (*http.Response, error) {
	var body io.Reader
	switch typed := data.(type) {
	case io.Reader:
		body = typed
	case []byte:
		body = bytes.NewReader(typed)
	case string:
		body = bytes.NewBufferString(typed)
	case nil:
	default:
		bts, err := json.Marshal(typed)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bts)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.Server+path, body)
	if err != nil {
		return nil, err
	}
	if method != http.MethodGet {
		// add csrftoken header
		if c.csrftoken == "" {
			if _, err := c.SystemInfo(ctx); err != nil {
				return nil, fmt.Errorf("error in harbor when get csrt token %w", err)
			}
		}
		req.Header.Add(csrfTokenHeader, c.csrftoken)
		// always add json content header Content-Type: application/json
		req.Header.Add("Content-Type", "application/json")
	}
	if c.Auth != nil {
		c.Auth(req)
	}

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		errobj := &harborerrors.Error{}
		if err = json.NewDecoder(resp.Body).Decode(errobj); err != nil {
			return resp, err
		}
		return resp, errobj
	}
	// update csrftoken if exist
	if method == http.MethodGet {
		if csrftoken := resp.Header.Get(csrfTokenHeader); csrftoken != "" {
			c.csrftoken = csrftoken
		}
	}

	// resp into writer
	switch into := decodeinto.(type) {
	case io.Writer:
		_, err := io.Copy(into, resp.Body)
		return resp, err
	case []byte:
		_, err := io.Copy(bytes.NewBuffer(into), resp.Body)
		return resp, err
	case *[]byte:
		_, err := io.Copy(bytes.NewBuffer(*into), resp.Body)
		return resp, err
	case nil:
	default:
		return resp, json.NewDecoder(resp.Body).Decode(decodeinto)
	}
	return resp, nil
}
