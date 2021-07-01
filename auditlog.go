package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goharbor/harbor/src/pkg/audit/model"
)

// GET /audit-logs
func (c *Client) ListAuditLogs(ctx context.Context, options CommonListOptions) ([]model.AuditLog, error) {
	path := fmt.Sprintf("/audit-logs?%s", options.toQuery().Encode())
	ret := []model.AuditLog{}
	err := c.doRequest(ctx, http.MethodPut, path, nil, &ret)
	return ret, err
}
