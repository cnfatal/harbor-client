package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/goharbor/harbor/src/pkg/audit/model"
	projectmodels "github.com/goharbor/harbor/src/pkg/project/models"
	"github.com/goharbor/harbor/src/testing/apitests/apilib"
)

// GET /projects/{project_id}/summary
func (c *Client) GetProjectSummary(ctx context.Context, projectID int) (apilib.ProjectSummary, error) {
	path := fmt.Sprintf("/projects/{project_id}/summary")
	ret := apilib.ProjectSummary{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

// PUT /projects/{project_id}
func (c *Client) UpdateProject(ctx context.Context, projectID int, project projectmodels.Project) error {
	path := fmt.Sprintf("/projects/%d", projectID)
	return c.doRequest(ctx, http.MethodPut, path, project, nil)
}

// GET /projects/{project_id}
func (c *Client) GetProject(ctx context.Context, projectID int) (projectmodels.Project, error) {
	path := fmt.Sprintf("/projects/%d", projectID)
	project := projectmodels.Project{}
	err := c.doRequest(ctx, http.MethodGet, path, project, nil)
	return project, err
}

// HEAD /projects/{project_id}
func (c *Client) HeadProject(ctx context.Context, project string) error {
	path := fmt.Sprintf("/projects/%s", project)
	return c.doRequest(ctx, http.MethodHead, path, project, nil)
}

// DELETE /projects/{project_id}
func (c *Client) DeleteProject(ctx context.Context, projectID int) error {
	path := fmt.Sprintf("/projects/%d", projectID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// GET /projects/{project_id}/_deletable
func (c *Client) GetProjectDeletable(ctx context.Context, projectID int) error {
	path := fmt.Sprintf("/projects/%d/_deletable", projectID)
	return c.doRequest(ctx, http.MethodGet, path, nil, nil)
}

// GET /projects/{project_name}/logs
func (c *Client) GetProjectLogs(ctx context.Context, project string, options CommonListOptions) error {
	path := fmt.Sprintf("/projects/%s/logs?%s", project, options.toQuery().Encode())
	ret := []model.AuditLog{}
	return c.doRequest(ctx, http.MethodGet, path, nil, &ret)
}

// POST /projects
func (c *Client) CreateProject(ctx context.Context, project projectmodels.Project) error {
	return c.doRequest(ctx, http.MethodPut, "/projects", project, nil)
}

type ListProjectsOptions struct {
	CommonListOptions
	Name   string `json:"name,omitempty"`
	Public bool   `json:"public,omitempty"`
	Owner  string `json:"owner,omitempty"`
}

func (o *ListProjectsOptions) toQuery() url.Values {
	value := o.CommonListOptions.toQuery()
	value["name"] = []string{o.Name}
	value["owner"] = []string{o.Owner}
	value["public"] = []string{strconv.FormatBool(o.Public)}
	return value
}

// GET /projects
func (c *Client) ListProjects(ctx context.Context, options ListProjectsOptions) error {
	path := fmt.Sprintf("/projects?%s", options.toQuery().Encode())
	ret := []projectmodels.Project{}
	return c.doRequest(ctx, http.MethodPut, path, nil, &ret)
}
