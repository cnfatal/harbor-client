package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/pkg/label/model"
)

// LabelColor must be '#xxxxxx' format
type LabelColor string

const (
	LabelColorRed    = "#C92100"
	LabelColorGreen  = "#00AB9A"
	LabelColorWhite  = "#FFFFFF"
	LabelColorYellow = "#FFDC0B"
)

// Label holds information used for a label
type Label struct {
	ID           int64     `orm:"pk;auto;column(id)" json:"id"`
	Name         string    `orm:"column(name)" json:"name"`
	Description  string    `orm:"column(description)" json:"description"`
	Color        string    `orm:"column(color)" json:"color"`
	Level        string    `orm:"column(level)" json:"-"`
	Scope        string    `orm:"column(scope)" json:"scope"`
	ProjectID    int64     `orm:"column(project_id)" json:"project_id"`
	CreationTime time.Time `orm:"column(creation_time);auto_now_add" json:"creation_time"`
	UpdateTime   time.Time `orm:"column(update_time);auto_now" json:"update_time"`
	Deleted      bool      `orm:"column(deleted)" json:"deleted"`
}

// POST /labels
func (c *Client) CreateGlobalLabel(ctx context.Context, key, desc string, color LabelColor) error {
	label := model.Label{
		Name:        key,
		Description: desc,
		Color:       string(color),
		Scope:       common.LabelScopeGlobal,
	}
	return c.doRequest(ctx, http.MethodPost, "/labels", label, nil)
}

// POST /labels
func (c *Client) CreateProjectLabel(ctx context.Context, projectID int64, key, desc string, color LabelColor) error {
	label := model.Label{
		Name:        key,
		Description: desc,
		Color:       string(color),
		Scope:       common.LabelScopeProject,
		ProjectID:   int64(projectID),
	}
	return c.doRequest(ctx, http.MethodPost, "/labels", label, nil)
}

// GET /labels?scope=g
func (c *Client) ListGlobalLabels(ctx context.Context) ([]model.Label, error) {
	labels := []model.Label{}
	if err := c.doRequest(ctx, http.MethodGet, "/labels?scope=g", nil, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// GET /labels?scope=p&project_id={id}
func (c *Client) ListProjectLabels(ctx context.Context, projectID int) ([]model.Label, error) {
	labels := []model.Label{}
	path := fmt.Sprintf("/labels?scope=p&project_id=%d", projectID)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// GET /labels/{id}
func (c *Client) GetLabel(ctx context.Context, id int) (model.Label, error) {
	path := fmt.Sprintf("/labels/%d", id)
	label := model.Label{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &label); err != nil {
		return label, err
	}
	return label, nil
}

// DELETE /labels/{id}
func (c *Client) DeleteLabel(ctx context.Context, id int) error {
	path := fmt.Sprintf("/labels/%d", id)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// PUT /labels/{id}
func (c *Client) UpdateLabel(ctx context.Context, label model.Label) error {
	path := fmt.Sprintf("/labels/%d", label.ID)
	return c.doRequest(ctx, http.MethodDelete, path, label, nil)
}

// POST /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/labels
// body:  {"id":2}
func (c *Client) AttachArtifactLabel(ctx context.Context, project, repository, reference string, labelID int64) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/labels", project, repository, reference)
	return c.doRequest(ctx, http.MethodPost, path, model.Label{ID: labelID}, nil)
}

// DELETE /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/labels/{label_id}
func (c *Client) DettachArtifactLabel(ctx context.Context, project, repository, reference string, labelID int) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/labels/%d", project, repository, reference, labelID)
	return c.doRequest(ctx, http.MethodDelete, path, model.Label{ID: int64(labelID)}, nil)
}

// GET /chartrepo/{repo}/charts/{name}/{version}/labels
func (c *Client) ListChartLabels(ctx context.Context, repo, name, version string) ([]model.Label, error) {
	path := fmt.Sprintf("/chartrepo/%s/charts/%s/%s/labels", repo, name, version)
	labels := []model.Label{}
	err := c.doRequest(ctx, http.MethodPost, path, &labels, nil)
	return labels, err
}

// POST /chartrepo/{repo}/charts/{name}/{version}/labels
// body:  {"id":2}
func (c *Client) AttachChartLabel(ctx context.Context, repo, name, version string, labelID int64) error {
	path := fmt.Sprintf("/chartrepo/%s/charts/%s/%s/labels", repo, name, version)
	return c.doRequest(ctx, http.MethodPost, path, model.Label{ID: labelID}, nil)
}

// Delete /chartrepo/{repo}/charts/{name}/{version}/labels
func (c *Client) DettachChartLabel(ctx context.Context, repo, name, version string, labelID int64) ([]model.Label, error) {
	path := fmt.Sprintf("/chartrepo/%s/charts/%s/%s/labels", repo, name, version)
	labels := []model.Label{}
	err := c.doRequest(ctx, http.MethodPost, path, &labels, nil)
	return labels, err
}

// AddArtifactLabelFromKey adds a label to artifact by the key. it find key name from all keys
// if key not exist,create and do it again
func (c *Client) AddArtifactLabelFromKey(ctx context.Context, project, repository, reference string, key, desc string, color LabelColor) error {
	labels, err := c.ListGlobalLabels(ctx)
	if err != nil {
		return nil
	}
	for _, label := range labels {
		if label.Name == key {
			return c.AttachArtifactLabel(ctx, project, repository, reference, label.ID)
		}
	}
	// create the label
	if err := c.CreateGlobalLabel(ctx, key, desc, color); err != nil {
		return err
	}
	// try again
	labels, err = c.ListGlobalLabels(ctx)
	if err != nil {
		return nil
	}
	for _, label := range labels {
		if label.Name == key {
			return c.AttachArtifactLabel(ctx, project, repository, reference, label.ID)
		}
	}
	// impossible
	return errors.New("unknow err,plaese try again")
}

func (c *Client) DettachArtifactLabelFromKey(ctx context.Context, project, repository, reference string, key string) error {
	labels, err := c.ListGlobalLabels(ctx)
	if err != nil {
		return nil
	}
	for _, label := range labels {
		if label.Name == key {
			return c.DettachArtifactLabel(ctx, project, repository, reference, int(label.ID))
		}
	}
	return errors.New("unknow err,plaese try again")
}
