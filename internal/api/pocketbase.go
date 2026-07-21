package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/certimate-go/cli/internal/models"
)

// ListResponse is the standard PocketBase response for list queries
type ListResponse[T any] struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Items      []T `json:"items"`
}

// buildQueryParams builds PocketBase query parameters
func buildQueryParams(filter, sort string, page, perPage int) string {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	if sort != "" {
		params.Set("sort", sort)
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		params.Set("perPage", fmt.Sprintf("%d", perPage))
	}
	return params.Encode()
}

// ListWorkflows fetches all workflows with optional filtering
func (c *Client) ListWorkflows(ctx context.Context, filter, sort string, page, perPage int) (*ListResponse[models.Workflow], error) {
	if sort == "" {
		sort = "-created"
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	path := "/api/collections/workflow/records?" + buildQueryParams(filter, sort, page, perPage)

	var result ListResponse[models.Workflow]
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	return &result, nil
}

// GetWorkflow fetches a single workflow by ID
func (c *Client) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	path := fmt.Sprintf("/api/collections/workflow/records/%s", id)

	var result models.Workflow
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}
	return &result, nil
}

// RunWorkflow executes a workflow via custom API
func (c *Client) RunWorkflow(ctx context.Context, workflowID string) (string, error) {
	path := fmt.Sprintf("/api/workflows/%s/runs", workflowID)
	body := map[string]string{"trigger": "manual"}

	var result struct {
		RunId string `json:"runId"`
	}
	if err := c.Do(ctx, "POST", path, body, &result); err != nil {
		return "", fmt.Errorf("run workflow: %w", err)
	}

	return result.RunId, nil
}

// CancelWorkflowRun cancels a running workflow
func (c *Client) CancelWorkflowRun(ctx context.Context, workflowID, runID string) error {
	path := fmt.Sprintf("/api/workflows/%s/runs/%s/cancel", workflowID, runID)
	return c.Do(ctx, "POST", path, nil, nil)
}

// ListWorkflowRuns fetches workflow execution history
func (c *Client) ListWorkflowRuns(ctx context.Context, workflowID string, page, perPage int) (*ListResponse[models.WorkflowRun], error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 30
	}

	// Use workflowRef which is the PocketBase field name
	filter := fmt.Sprintf("workflowRef='%s'", workflowID)
	path := "/api/collections/workflow_run/records?" + buildQueryParams(filter, "-created", page, perPage)

	var result ListResponse[models.WorkflowRun]
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("list workflow runs: %w", err)
	}
	return &result, nil
}

// GetWorkflowRun fetches a single workflow run by ID
func (c *Client) GetWorkflowRun(ctx context.Context, runID string) (*models.WorkflowRun, error) {
	path := fmt.Sprintf("/api/collections/workflow_run/records/%s", runID)

	var result models.WorkflowRun
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get workflow run: %w", err)
	}
	return &result, nil
}

// ListCertificates fetches all certificates with optional filtering
func (c *Client) ListCertificates(ctx context.Context, filter, sort string, page, perPage int) (*ListResponse[models.Certificate], error) {
	if sort == "" {
		sort = "-created"
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	path := "/api/collections/certificate/records?" + buildQueryParams(filter, sort, page, perPage)

	var result ListResponse[models.Certificate]
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("list certificates: %w", err)
	}
	return &result, nil
}

// GetCertificate fetches a single certificate by ID
func (c *Client) GetCertificate(ctx context.Context, id string) (*models.Certificate, error) {
	path := fmt.Sprintf("/api/collections/certificate/records/%s", id)

	var result models.Certificate
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get certificate: %w", err)
	}
	return &result, nil
}

// DownloadCertificate downloads a certificate in the specified format.
//
// The server responds with resp.Ok's envelope
// {"code":0,"msg":"success","data":{"zipBytes":"<base64 zip>"}}, where zipBytes
// is a base64-encoded ZIP archive containing the certificate/key files. Go's
// encoding/json base64-decodes []byte for us. Note that the server returns
// HTTP 200 even on failure, so a non-zero code must be treated as an error.
func (c *Client) DownloadCertificate(ctx context.Context, certificateID, format string) ([]byte, error) {
	path := fmt.Sprintf("/api/certificates/%s/download", certificateID)
	body := map[string]string{"fileFormat": format}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ZipBytes []byte `json:"zipBytes"`
		} `json:"data"`
	}
	if err := c.Do(ctx, "POST", path, body, &result); err != nil {
		return nil, fmt.Errorf("download certificate: %w", err)
	}

	if result.Code != 0 {
		msg := result.Msg
		if msg == "" {
			msg = "unknown error"
		}
		return nil, fmt.Errorf("download certificate: %s", msg)
	}

	return result.Data.ZipBytes, nil
}

// ListAccess fetches all access credentials with optional filtering
func (c *Client) ListAccess(ctx context.Context, filter, sort string, page, perPage int) (*ListResponse[models.Access], error) {
	if sort == "" {
		sort = "-created"
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	path := "/api/collections/access/records?" + buildQueryParams(filter, sort, page, perPage)

	var result ListResponse[models.Access]
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("list access: %w", err)
	}
	return &result, nil
}

// GetAccess fetches a single access by ID
func (c *Client) GetAccess(ctx context.Context, id string) (*models.Access, error) {
	path := fmt.Sprintf("/api/collections/access/records/%s", id)

	var result models.Access
	if err := c.Do(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("get access: %w", err)
	}
	return &result, nil
}

// CreateAccess creates a new access credential
func (c *Client) CreateAccess(ctx context.Context, access *models.Access) (*models.Access, error) {
	path := "/api/collections/access/records"

	var result models.Access
	if err := c.Do(ctx, "POST", path, access, &result); err != nil {
		return nil, fmt.Errorf("create access: %w", err)
	}
	return &result, nil
}

// UpdateAccess updates an existing access credential
func (c *Client) UpdateAccess(ctx context.Context, id string, access *models.Access) (*models.Access, error) {
	path := fmt.Sprintf("/api/collections/access/records/%s", id)

	var result models.Access
	if err := c.Do(ctx, "PATCH", path, access, &result); err != nil {
		return nil, fmt.Errorf("update access: %w", err)
	}
	return &result, nil
}

// DeleteAccess deletes an access credential
func (c *Client) DeleteAccess(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/collections/access/records/%s", id)
	return c.Do(ctx, "DELETE", path, nil, nil)
}

// CreateWorkflow creates a new workflow
func (c *Client) CreateWorkflow(ctx context.Context, workflow *models.Workflow) (*models.Workflow, error) {
	path := "/api/collections/workflow/records"

	var result models.Workflow
	if err := c.Do(ctx, "POST", path, workflow, &result); err != nil {
		return nil, fmt.Errorf("create workflow: %w", err)
	}
	return &result, nil
}

// UpdateWorkflow updates an existing workflow
func (c *Client) UpdateWorkflow(ctx context.Context, id string, workflow *models.Workflow) (*models.Workflow, error) {
	path := fmt.Sprintf("/api/collections/workflow/records/%s", id)

	var result models.Workflow
	if err := c.Do(ctx, "PATCH", path, workflow, &result); err != nil {
		return nil, fmt.Errorf("update workflow: %w", err)
	}
	return &result, nil
}

// DeleteWorkflow deletes a workflow
func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/collections/workflow/records/%s", id)
	return c.Do(ctx, "DELETE", path, nil, nil)
}
