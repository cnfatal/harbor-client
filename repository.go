package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tomnomnom/linkheader"
)

type RepositoriesListOptions struct {
	CommonListOptions
}

// https://github.com/goharbor/harbor/blob/4e1f6633afb824cd16341044a0e82f4f1f230cd2/src/pkg/repository/model/model.go#L34
type Repository struct {
	RepositoryID int64     `json:"repositoryID,omitempty"`
	Name         string    `json:"name,omitempty"`
	ProjectID    int64     `json:"projectID,omitempty"`
	Description  string    `json:"description,omitempty"`
	PullCount    int64     `json:"pullCount,omitempty"`
	StarCount    int64     `json:"starCount,omitempty"`
	CreationTime time.Time `json:"creationTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`
}

type RepositoryList struct {
	Total int          `json:"total,omitempty"`
	Items []Repository `json:"items,omitempty"`
	Next  string       `json:"next,omitempty"`
}

// ListRepositories
// GET /projects/{project_name}/repositories
func (c *Client) ListRepositories(ctx context.Context, project string, options RepositoriesListOptions) (RepositoryList, error) {
	path := fmt.Sprintf("/projects/%s/repositories", project)
	ret := RepositoryList{}
	resp, err := c.doRequestWithResponse(ctx, http.MethodGet, path, nil, &ret)
	if err != nil {
		return ret, err
	}
	if total, err := strconv.Atoi(resp.Header.Get(xTotalCountHeader)); err != nil {
		ret.Total = total
	}
	// Link: </api/v2.0/projects/library/repositories?page=2&page_size=10>; rel="next"
	for _, link := range linkheader.Parse(resp.Header.Get(linkHeader)) {
		if link.Rel == "next" {
			ret.Next = link.URL
		}
	}
	return ret, nil
}

// GetRepositories
// GET /projects/{project_name}/repositories/{repository_name}
func (c *Client) GetRepository(ctx context.Context, project, repository string) (Repository, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s", project, repository)
	ret := Repository{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

// PutRepository
// PUT /projects/{project_name}/repositories/{repository_name}
func (c *Client) UpdateRepository(ctx context.Context, project string, repository Repository) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s", project, repository.Name)
	return c.doRequest(ctx, http.MethodPut, path, repository, nil)
}

// DeleteRepository
// DELETE /projects/{project_name}/repositories/{repository_name}
func (c *Client) DeleteRepository(ctx context.Context, project, repository string) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s", project, repository)
	return c.doRequest(ctx, http.MethodDelete, path, repository, nil)
}
