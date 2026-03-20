package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/certimate-go/cli/internal/models"
)

// TableFormatter formats output as a table
type TableFormatter struct {
	Writer io.Writer
}

// Format outputs data as a table
func (f *TableFormatter) Format(data interface{}) error {
	if f.Writer == nil {
		f.Writer = os.Stdout
	}

	w := tabwriter.NewWriter(f.Writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Handle different data types
	switch v := data.(type) {
	case *models.Workflow:
		return f.formatSingle(w, v)
	case *models.Certificate:
		return f.formatSingle(w, v)
	case *models.Access:
		return f.formatSingle(w, v)
	case *models.WorkflowRun:
		return f.formatSingle(w, v)
	case []models.Workflow:
		return f.formatWorkflows(w, v)
	case []models.Certificate:
		return f.formatCertificates(w, v)
	case []models.Access:
		return f.formatAccessList(w, v)
	case []models.WorkflowRun:
		return f.formatWorkflowRuns(w, v)
	default:
		// Handle list responses
		rv := reflect.ValueOf(data)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		if rv.Kind() == reflect.Struct {
			// Check if it's a ListResponse
			if items := rv.FieldByName("Items"); items.IsValid() {
				switch items.Interface().(type) {
				case []models.Workflow:
					return f.formatWorkflows(w, items.Interface().([]models.Workflow))
				case []models.Certificate:
					return f.formatCertificates(w, items.Interface().([]models.Certificate))
				case []models.Access:
					return f.formatAccessList(w, items.Interface().([]models.Access))
				case []models.WorkflowRun:
					return f.formatWorkflowRuns(w, items.Interface().([]models.WorkflowRun))
				}
			}
		}

		// Fallback to single item formatting
		return f.formatSingle(w, data)
	}
}

func (f *TableFormatter) formatSingle(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

func (f *TableFormatter) formatWorkflows(w io.Writer, workflows []models.Workflow) error {
	fmt.Fprintln(w, "ID\tNAME\tTRIGGER\tENABLED\tCREATED")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, wf := range workflows {
		enabled := "false"
		if wf.Enabled {
			enabled = "true"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			wf.ID,
			truncate(wf.Name, 30),
			wf.Trigger,
			enabled,
			wf.Created.Format("2006-01-02"),
		)
	}

	return nil
}

func (f *TableFormatter) formatCertificates(w io.Writer, certs []models.Certificate) error {
	fmt.Fprintln(w, "ID\tSUBJECT\tEXPIRES\tSTATUS\tCREATED")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, cert := range certs {
		status := "unknown"
		expires := "-"

		if !cert.NotAfter.IsZero() {
			expires = cert.NotAfter.Format("2006-01-02")
			if cert.IsExpired() {
				status = "EXPIRED"
			} else if cert.DaysUntilExpiry() <= 30 {
				status = "EXPIRING"
			} else {
				status = "valid"
			}
		}

		subject := strings.Join(cert.SubjectAltNames, ", ")
		if len(subject) > 30 {
			subject = subject[:27] + "..."
		}

		created := "-"
		if !cert.Created.IsZero() {
			created = cert.Created.Format("2006-01-02")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			cert.ID,
			subject,
			expires,
			status,
			created,
		)
	}

	return nil
}

func (f *TableFormatter) formatAccessList(w io.Writer, accessList []models.Access) error {
	fmt.Fprintln(w, "ID\tNAME\tPROVIDER\tCREATED")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, acc := range accessList {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			acc.ID,
			truncate(acc.Name, 30),
			acc.Provider,
			acc.Created.Format("2006-01-02"),
		)
	}

	return nil
}

func (f *TableFormatter) formatWorkflowRuns(w io.Writer, runs []models.WorkflowRun) error {
	fmt.Fprintln(w, "ID\tSTATUS\tTRIGGER\tSTARTED\tENDED")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, run := range runs {
		started := "-"
		if !run.StartedAt.IsZero() {
			started = run.StartedAt.Format("2006-01-02 15:04")
		}
		ended := "-"
		if !run.EndedAt.IsZero() {
			ended = run.EndedAt.Format("2006-01-02 15:04")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			run.ID,
			run.Status,
			run.Trigger,
			started,
			ended,
		)
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
