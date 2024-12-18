package libselectclause

import (
	"fmt"
	"strings"

	types "github.com/opentofu/opentofu/internal/hof/lib/mantis/cql/shared"
)

// FormatResults formats the query results into a table-like string
func FormatResults(result types.QueryResult, config types.QueryConfig) string {
	var output strings.Builder

	if len(result.Matches) == 0 {
		return "No matches found in the configurations.\n"
	}

	// Determine fields to display
	fields := determineDisplayFields(result, config.Select)

	// Print header
	printHeader(&output, fields)

	// Print separator
	printSeparator(&output, fields)

	// Print values
	for _, matches := range result.Matches {
		for _, match := range matches {
			FormatMatchAsTable(&output, match, fields)
		}
	}

	return output.String()
}

// FormatMatchAsTable formats a single match as a table row
func FormatMatchAsTable(output *strings.Builder, match types.Match, fields []string) {
	// Print only the requested fields
	for _, field := range fields {
		value := findFieldValue(match, field)
		fmt.Fprintf(output, "%-20s", value)
	}
	output.WriteString("\n")
}

// Helper functions

func determineDisplayFields(result types.QueryResult, selects []string) []string {
	if len(selects) == 1 && selects[0] == "*" {
		fieldSet := make(map[string]bool)
		for _, matches := range result.Matches {
			for _, match := range matches {
				fieldSet[match.Path] = true
				for _, child := range match.Children {
					fieldSet[match.Path+"."+child.Path] = true
				}
			}
		}
		var fields []string
		for field := range fieldSet {
			fields = append(fields, field)
		}
		return fields
	}
	return selects
}

func printHeader(output *strings.Builder, fields []string) {
	// Print only the field names from SELECT
	for _, h := range fields {
		fmt.Fprintf(output, "%-20s", h)
	}
	output.WriteString("\n")
}

func printSeparator(output *strings.Builder, fields []string) {
	// Print separator only for selected fields
	for range fields {
		output.WriteString(strings.Repeat("-", 20))
	}
	output.WriteString("\n")
}

func findFieldValue(match types.Match, field string) string {
	// First check Fields map
	if val, ok := match.Fields[field]; ok {
		switch v := val.(type) {
		case []interface{}:
			return fmt.Sprintf("%v", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	// Then check other fields
	if field == "*" {
		return match.Value
	}

	if match.Path == field {
		return match.Value
	}
	for _, child := range match.Children {
		if child.Path == field {
			return child.Value
		}
	}
	return ""
}
