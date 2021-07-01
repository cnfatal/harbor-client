package client

import (
	"context"
	"fmt"
	"net/http"
)

// GET /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/scan/{report_id}/log
func (c *Client) GetScanReportLog(ctx context.Context, project, repository, reference string, reportID int) ([]byte, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/scan/%d/log", project, repository, reference, reportID)
	log := []byte{}
	err := c.doRequest(ctx, http.MethodGet, path, nil, &log)
	return log, err
}

// POST /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/scan
func (c *Client) ScanArtifact(ctx context.Context, project, repository, reference string) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/scan", project, repository, reference)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}
