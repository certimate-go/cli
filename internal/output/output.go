package output

import (
	"encoding/json"
	"io"
	"os"
)

// Formatter defines the output formatter interface
type Formatter interface {
	Format(data interface{}) error
}

// Print outputs data in the specified format
func Print(data interface{}, format string) error {
	var formatter Formatter

	switch format {
	case "table":
		formatter = &TableFormatter{Writer: os.Stdout}
	default:
		formatter = &JSONFormatter{Writer: os.Stdout}
	}

	return formatter.Format(data)
}

// PrintTo outputs data to the specified writer
func PrintTo(data interface{}, format string, w io.Writer) error {
	var formatter Formatter

	switch format {
	case "table":
		formatter = &TableFormatter{Writer: w}
	default:
		formatter = &JSONFormatter{Writer: w}
	}

	return formatter.Format(data)
}

// PrintError outputs an error in JSON format
func PrintError(err error) error {
	result := map[string]interface{}{
		"error":   true,
		"message": err.Error(),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
