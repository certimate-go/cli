package models

// WorkflowTriggerType defines the trigger type for workflows
type WorkflowTriggerType string

const (
	WorkflowTriggerManual    WorkflowTriggerType = "manual"
	WorkflowTriggerScheduled WorkflowTriggerType = "scheduled"
)

// Workflow represents a certificate workflow
type Workflow struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Trigger       WorkflowTriggerType `json:"trigger"`
	Enabled       bool                `json:"enabled"`
	Graph         WorkflowGraph       `json:"graph"`
	GraphDraft    WorkflowGraph       `json:"graphDraft"`
	HasDraft      bool                `json:"hasDraft"`
	LastRunStatus WorkflowRunStatus   `json:"lastRunStatus"`
	LastRunTime   CustomTime          `json:"lastRunTime"`
	Created       CustomTime          `json:"created"`
	Updated       CustomTime          `json:"updated"`
}

// WorkflowGraph represents the workflow execution graph
type WorkflowGraph struct {
	Nodes []WorkflowNode `json:"nodes"`
	Edges []WorkflowEdge `json:"edges"`
}

// WorkflowNode represents a node in the workflow graph
type WorkflowNode struct {
	ID     string           `json:"id"`
	Type   string           `json:"type"`
	Data   map[string]any   `json:"data,omitempty"`
	Blocks []WorkflowNode   `json:"blocks,omitempty"`
}

// WorkflowEdge represents an edge between nodes
type WorkflowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}
