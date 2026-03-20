package models

// WorkflowRunStatus defines the status of a workflow run
type WorkflowRunStatus string

const (
	WorkflowRunStatusPending    WorkflowRunStatus = "pending"
	WorkflowRunStatusProcessing WorkflowRunStatus = "processing"
	WorkflowRunStatusSucceeded  WorkflowRunStatus = "succeeded"
	WorkflowRunStatusFailed     WorkflowRunStatus = "failed"
	WorkflowRunStatusCanceled   WorkflowRunStatus = "canceled"
)

// WorkflowRun represents a workflow execution record
type WorkflowRun struct {
	ID         string            `json:"id"`
	WorkflowID string            `json:"workflowRef"`
	Status     WorkflowRunStatus `json:"status"`
	Trigger    string            `json:"trigger"`
	StartedAt  CustomTime        `json:"startedAt"`
	EndedAt    CustomTime        `json:"endedAt"`
	Error      string            `json:"error"`
	Logs       []WorkflowRunLog  `json:"logs"`
	Created    CustomTime        `json:"created"`
	Updated    CustomTime        `json:"updated"`
}

// WorkflowRunLog represents a log entry in a workflow run
type WorkflowRunLog struct {
	NodeID    string     `json:"nodeId"`
	NodeType  string     `json:"nodeType"`
	Message   string     `json:"message"`
	Timestamp CustomTime `json:"timestamp"`
}
