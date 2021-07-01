package client

import (
	"context"
	"fmt"
	"net/http"
)

// https://github.com/goharbor/harbor/blob/4e1f6633afb824cd16341044a0e82f4f1f230cd2/src/controller/icon/controller.go#L51
type Icon struct {
	ContentType string `json:"content_type,omitempty"`
	Content     string `json:"content,omitempty"` // base64 encoded
}

// GET /icons/{digest}
func (c *Client) GetIcon(ctx context.Context, digest string) (Icon, error) {
	ret := Icon{}
	if err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/icons/%s", digest), nil, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}
