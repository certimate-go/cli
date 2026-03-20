package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/certimate-go/cli/internal/api"
	"github.com/certimate-go/cli/internal/config"
	"github.com/certimate-go/cli/internal/models"
	"github.com/certimate-go/cli/internal/output"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage certificate workflows",
	Long:  `List, view, and execute certificate workflows.`,
}

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflows",
	RunE:  runWorkflowList,
}

var workflowGetCmd = &cobra.Command{
	Use:   "get WORKFLOW_ID",
	Short: "Get workflow details",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowGet,
}

var workflowRunCmd = &cobra.Command{
	Use:   "run WORKFLOW_ID",
	Short: "Execute a workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowRun,
}

var workflowCancelCmd = &cobra.Command{
	Use:   "cancel WORKFLOW_ID RUN_ID",
	Short: "Cancel a running workflow",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkflowCancel,
}

var workflowRunsCmd = &cobra.Command{
	Use:   "runs WORKFLOW_ID",
	Short: "List workflow execution history",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowRuns,
}

var workflowCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workflow",
	RunE:  runWorkflowCreate,
}

var workflowEditCmd = &cobra.Command{
	Use:   "edit WORKFLOW_ID",
	Short: "Edit an existing workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowEdit,
}

var workflowDeleteCmd = &cobra.Command{
	Use:   "delete WORKFLOW_ID",
	Short: "Delete a workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowDelete,
}

var workflowEnableCmd = &cobra.Command{
	Use:   "enable WORKFLOW_ID",
	Short: "Enable a workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowEnable,
}

var workflowDisableCmd = &cobra.Command{
	Use:   "disable WORKFLOW_ID",
	Short: "Disable a workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowDisable,
}

var (
	workflowWait    bool
	workflowTimeout int
	workflowFilter  string
	workflowLimit   int

	// Create/Edit flags
	workflowName        string
	workflowDescription string
	workflowGraphDraft  string
	workflowHasDraft    bool
	workflowForce       bool
)

func init() {
	rootCmd.AddCommand(workflowCmd)
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowGetCmd)
	workflowCmd.AddCommand(workflowRunCmd)
	workflowCmd.AddCommand(workflowCancelCmd)
	workflowCmd.AddCommand(workflowRunsCmd)
	workflowCmd.AddCommand(workflowCreateCmd)
	workflowCmd.AddCommand(workflowEditCmd)
	workflowCmd.AddCommand(workflowDeleteCmd)
	workflowCmd.AddCommand(workflowEnableCmd)
	workflowCmd.AddCommand(workflowDisableCmd)

	workflowListCmd.Flags().StringVar(&workflowFilter, "filter", "", "PocketBase filter expression")
	workflowListCmd.Flags().IntVar(&workflowLimit, "limit", 100, "Maximum number of results")

	workflowRunCmd.Flags().BoolVar(&workflowWait, "wait", false, "Wait for workflow to complete")
	workflowRunCmd.Flags().IntVar(&workflowTimeout, "timeout", 300, "Timeout in seconds when waiting")

	workflowRunsCmd.Flags().IntVar(&workflowLimit, "limit", 30, "Maximum number of results")

	// Create command flags
	workflowCreateCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name")
	workflowCreateCmd.Flags().StringVar(&workflowDescription, "description", "", "Workflow description")
	workflowCreateCmd.Flags().StringVar(&workflowGraphDraft, "graph-draft", "", "Workflow graph structure (JSON string or @path/to/file.json)")
	workflowCreateCmd.Flags().BoolVar(&workflowHasDraft, "has-draft", true, "Mark workflow as having draft")
	workflowCreateCmd.MarkFlagRequired("name")
	workflowCreateCmd.MarkFlagRequired("graph-draft")

	// Edit command flags
	workflowEditCmd.Flags().StringVar(&workflowName, "name", "", "New workflow name")
	workflowEditCmd.Flags().StringVar(&workflowDescription, "description", "", "New workflow description")
	workflowEditCmd.Flags().StringVar(&workflowGraphDraft, "graph-draft", "", "New workflow graph structure (JSON string or @path/to/file.json)")

	// Delete command flags
	workflowDeleteCmd.Flags().BoolVar(&workflowForce, "force", false, "Skip confirmation prompt")
}

func getAPIClient() (*api.Client, error) {
	cfg, err := config.Load(profile)
	if err != nil {
		return nil, err
	}
	return api.NewClient(cfg), nil
}

func runWorkflowList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	resp, err := client.ListWorkflows(context.Background(), workflowFilter, "", 1, workflowLimit)
	if err != nil {
		return err
	}

	return output.Print(resp, outputFmt)
}

func runWorkflowGet(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflow, err := client.GetWorkflow(context.Background(), args[0])
	if err != nil {
		return err
	}

	return output.Print(workflow, outputFmt)
}

func runWorkflowRun(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]
	runID, err := client.RunWorkflow(context.Background(), workflowID)
	if err != nil {
		return err
	}

	debug := viper.GetBool("debug")
	if debug {
		fmt.Fprintf(os.Stderr, "Debug: started workflow run, runId=%s\n", runID)
	}

	result := map[string]interface{}{
		"status":      "started",
		"run_id":      runID,
		"workflow_id": workflowID,
		"message":     "Workflow execution started",
	}

	if workflowWait {
		// Poll for completion by checking workflow's lastRunStatus
		timeout := time.Duration(workflowTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Wait a moment for the workflow run to start
		time.Sleep(1 * time.Second)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				result["status"] = "timeout"
				result["message"] = fmt.Sprintf("Timed out waiting for workflow after %v", timeout)
				return output.Print(result, outputFmt)

			case <-ticker.C:
				// Poll the workflow to get lastRunStatus
				workflow, err := client.GetWorkflow(ctx, workflowID)
				if err != nil {
					if debug {
						fmt.Fprintf(os.Stderr, "Debug: error getting workflow status: %v\n", err)
					}
					continue
				}

				if debug {
					fmt.Fprintf(os.Stderr, "Debug: workflow lastRunStatus = %s\n", workflow.LastRunStatus)
				}

				switch workflow.LastRunStatus {
				case models.WorkflowRunStatusSucceeded:
					result["status"] = "succeeded"
					result["message"] = "Workflow completed successfully"
					return output.Print(result, outputFmt)

				case models.WorkflowRunStatusFailed:
					result["status"] = "failed"
					result["message"] = "Workflow execution failed"
					return output.Print(result, outputFmt)

				case models.WorkflowRunStatusCanceled:
					result["status"] = "canceled"
					result["message"] = "Workflow was canceled"
					return output.Print(result, outputFmt)
				}
			}
		}
	}

	return output.Print(result, outputFmt)
}

func runWorkflowCancel(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]
	runID := args[1]

	if err := client.CancelWorkflowRun(context.Background(), workflowID, runID); err != nil {
		return err
	}

	result := map[string]interface{}{
		"status":      "canceled",
		"workflow_id": workflowID,
		"run_id":      runID,
		"message":     "Workflow run canceled",
	}

	return output.Print(result, outputFmt)
}

func runWorkflowRuns(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	resp, err := client.ListWorkflowRuns(context.Background(), args[0], 1, workflowLimit)
	if err != nil {
		return err
	}

	return output.Print(resp, outputFmt)
}

func runWorkflowCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Parse graphDraft JSON
	graphDraft, err := parseGraphDraftInput(workflowGraphDraft)
	if err != nil {
		return err
	}

	workflow := &models.Workflow{
		Name:        workflowName,
		Description: workflowDescription,
		GraphDraft:  graphDraft,
		HasDraft:    workflowHasDraft,
		Enabled:     false, // New workflows start disabled
	}

	created, err := client.CreateWorkflow(context.Background(), workflow)
	if err != nil {
		return err
	}

	return output.Print(created, outputFmt)
}

func runWorkflowEdit(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]

	// Fetch existing workflow
	existing, err := client.GetWorkflow(context.Background(), workflowID)
	if err != nil {
		return err
	}

	// Apply updates (only if flag provided)
	if workflowName != "" {
		existing.Name = workflowName
	}
	if cmd.Flags().Changed("description") {
		existing.Description = workflowDescription
	}
	if workflowGraphDraft != "" {
		graphDraft, err := parseGraphDraftInput(workflowGraphDraft)
		if err != nil {
			return err
		}
		existing.GraphDraft = graphDraft
	}

	updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
	if err != nil {
		return err
	}

	return output.Print(updated, outputFmt)
}

func runWorkflowDelete(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]

	// Fetch workflow to show what will be deleted (unless --force)
	if !workflowForce {
		existing, err := client.GetWorkflow(context.Background(), workflowID)
		if err != nil {
			return err
		}
		fmt.Printf("Will delete workflow: %s (%s)\n", existing.Name, workflowID)
		fmt.Print("Confirm? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	if err := client.DeleteWorkflow(context.Background(), workflowID); err != nil {
		return err
	}

	fmt.Printf("Workflow %s deleted\n", workflowID)
	return nil
}

func runWorkflowEnable(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]

	// Fetch existing workflow
	existing, err := client.GetWorkflow(context.Background(), workflowID)
	if err != nil {
		return err
	}

	if existing.Enabled {
		fmt.Printf("Workflow %s is already enabled\n", workflowID)
		return nil
	}

	existing.Enabled = true

	updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
	if err != nil {
		return err
	}

	return output.Print(updated, outputFmt)
}

func runWorkflowDisable(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	workflowID := args[0]

	// Fetch existing workflow
	existing, err := client.GetWorkflow(context.Background(), workflowID)
	if err != nil {
		return err
	}

	if !existing.Enabled {
		fmt.Printf("Workflow %s is already disabled\n", workflowID)
		return nil
	}

	existing.Enabled = false

	updated, err := client.UpdateWorkflow(context.Background(), workflowID, existing)
	if err != nil {
		return err
	}

	return output.Print(updated, outputFmt)
}

// parseGraphDraftInput parses graphDraft from JSON string or file path
func parseGraphDraftInput(input string) (models.WorkflowGraph, error) {
	var data []byte
	var err error

	if filePath, ok := strings.CutPrefix(input, "@"); ok {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return models.WorkflowGraph{}, fmt.Errorf("failed to read graph-draft file: %w", err)
		}
	} else {
		data = []byte(input)
	}

	var graphDraft models.WorkflowGraph
	if err := json.Unmarshal(data, &graphDraft); err != nil {
		return models.WorkflowGraph{}, fmt.Errorf("invalid graph-draft JSON format: %w", err)
	}

	return graphDraft, nil
}
