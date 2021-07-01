package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/goharbor/harbor/src/pkg/artifact"
	"github.com/goharbor/harbor/src/pkg/scan/vuln"
	"helm.sh/helm/v3/pkg/chart"
)

/*
 * Redefined  AdditionLink,Tag,Label
 * To avoid import lots of useless dependencies from harbor like beego,borm etc..
 */
type Artifact struct {
	artifact.Artifact
	Tags          []Tag                               `json:"tags"`           // the list of tags that attached to the artifact
	AdditionLinks map[string]AdditionLink             `json:"addition_links"` // the resource link for build history(image), values.yaml(chart), dependency(chart), etc
	Labels        []Label                             `json:"labels"`
	ScanOverview  map[string]vuln.NativeReportSummary `json:"scan_overview"`
}

// AdditionLink is a link via that the addition can be fetched
type AdditionLink struct {
	HREF     string `json:"href"`
	Absolute bool   `json:"absolute"` // specify the href is an absolute URL or not
}

// Tag is the overall view of tag
type Tag struct {
	ID           int64     `orm:"pk;auto;column(id)" json:"id"`
	RepositoryID int64     `orm:"column(repository_id)" json:"repository_id"` // tags are the resources of repository, one repository only contains one same name tag
	ArtifactID   int64     `orm:"column(artifact_id)" json:"artifact_id"`     // the artifact ID that the tag attaches to, it changes when pushing a same name but different digest artifact
	Name         string    `orm:"column(name)" json:"name"`
	PushTime     time.Time `orm:"column(push_time)" json:"push_time"`
	PullTime     time.Time `orm:"column(pull_time)" json:"pull_time"`
	Immutable    bool      `json:"immutable"`
	Signed       bool      `json:"signed"`
}

type Vulnerabilities map[string]vuln.Report

type GetArtifactOptions struct {
	WithTag             bool
	WithScanOverview    bool
	WithLabel           bool
	WithImmutableStatus bool
	WithSignature       bool
}

func (o GetArtifactOptions) toQuery() url.Values {
	return url.Values{
		"with_tag":              []string{strconv.FormatBool(o.WithTag)},
		"with_scan_overview":    []string{strconv.FormatBool(o.WithScanOverview)},
		"with_label":            []string{strconv.FormatBool(o.WithLabel)},
		"with_signature":        []string{strconv.FormatBool(o.WithSignature)},
		"with_immutable_status": []string{strconv.FormatBool(o.WithImmutableStatus)},
	}
}

type Addition string

const (
	AdditionBuildHistory    = "build_history"
	AdditionHelmValues      = "values.yaml"
	AdditionHelmReadme      = "readme.md"
	AdditionDependencies    = "dependencies"
	AdditionVulnerabilities = "vulnerabilities"
)

// GET /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/additions/{addition}
func (c *Client) GetArtifactadditions(ctx context.Context, project, repository, reference string, addition Addition) ([]byte, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/additions/%s", project, repository, reference, addition)
	ret := []byte{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) GetArtifactadditionVulnerabilities(ctx context.Context, project, repository, reference string) (Vulnerabilities, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/additions/%s", project, repository, reference, AdditionVulnerabilities)
	// https://github.com/goharbor/harbor/blob/c39345da96d887acb47d2b1e7cf1adafca5db1bb/src/server/v2.0/handler/artifact.go#L346
	ret := Vulnerabilities{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) GetArtifactadditionDependencies(ctx context.Context, project, repository, reference string) ([]chart.Dependency, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/additions/%s", project, repository, reference, AdditionDependencies)
	ret := []chart.Dependency{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type CopyArtifactResponse struct {
	Location string `json:"location"` // The location of the resource
}

// POST /projects/{project_name}/repositories/{repository_name}/artifacts
// from: The artifact from which the new artifact is copied from, the format should be "project/repository:tag" or "project/repository@digest".
func (c *Client) CopyArtifact(ctx context.Context, project, repository, from string) (CopyArtifactResponse, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts", project, repository)
	ret := CopyArtifactResponse{}
	resp, err := c.doRequestWithResponse(ctx, http.MethodPost, path, nil, &ret)
	if err != nil {
		return ret, err
	}
	ret.Location = resp.Header.Get("Location")
	return ret, nil
}

type ListArtifactsOptions struct {
	CommonListOptions
	GetArtifactOptions
}

func (o *ListArtifactsOptions) toQuery() url.Values {
	values := o.GetArtifactOptions.toQuery()
	for k := range o.CommonListOptions.toQuery() {
		values.Set(k, values.Get(k))
	}
	return values
}

// GET /projects/{project_name}/repositories/{repository_name}/artifacts
func (c *Client) ListArtifacts(ctx context.Context, project, repository string, options ListArtifactsOptions) ([]Artifact, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts?%s", project, repository, options.toQuery().Encode())
	ret := []Artifact{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// DELETE /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags/{tag_name}
func (c *Client) DeleteArtifactTag(ctx context.Context, project, repository, reference, tag string) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/tags/%s", project, repository, reference, tag)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// GET /repositories/{{repository_name}}/artifacts/{{reference}}?
func (c *Client) GetArtifact(ctx context.Context, project, repository, reference string, options GetArtifactOptions) (*Artifact, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s?%s", project, repository, reference, options.toQuery().Encode())
	ret := &Artifact{}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// DELETE /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}
func (c *Client) DeleteArtifact(ctx context.Context, project, repository, reference string) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s", project, repository, reference)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// POST /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags
func (c *Client) CreateArtifactTag(ctx context.Context, project, repository, reference string, tag Tag) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/tags", project, repository, reference)
	return c.doRequest(ctx, http.MethodPost, path, tag, nil)
}

type ListTagsOptions struct {
	CommonListOptions
	WithImmutableStatus bool
	WithSignature       bool
}

func (o *ListTagsOptions) toQuery() url.Values {
	values := o.CommonListOptions.toQuery()
	values.Set("with_signature", strconv.FormatBool(o.WithSignature))
	values.Set("with_immutable_status", strconv.FormatBool(o.WithImmutableStatus))
	return values
}

// GET /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags
func (c *Client) ListArtifactTags(ctx context.Context, project, repository, reference string, options ListTagsOptions) ([]Tag, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/tags?%s", project, repository, reference, options.toQuery().Encode())
	ret := []Tag{}
	if err := c.doRequest(ctx, http.MethodPost, path, ret, nil); err != nil {
		return ret, err
	}
	return ret, nil
}
