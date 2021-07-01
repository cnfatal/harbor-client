package client

import (
	"context"
	"net/http"

	"github.com/goharbor/harbor/src/controller/systeminfo"
)

// SystemInfo
// GET /systeminfo
func (c *Client) SystemInfo(ctx context.Context) (systeminfo.Data, error) {
	info := systeminfo.Data{}
	if err := c.doRequest(ctx, http.MethodGet, "/systeminfo", nil, &info); err != nil {
		return info, err
	}
	return info, nil
}

type OIDCPing struct {
	Url        string `json:"url,omitempty"`
	VerifyCert string `json:"verify_cert,omitempty"`
}

// GET /system/oidc/ping
func (c *Client) OIDCPing(ctx context.Context) (OIDCPing, error) {
	info := OIDCPing{}
	if err := c.doRequest(ctx, http.MethodGet, "/systeminfo", nil, info); err != nil {
		return info, err
	}
	return info, nil
}
