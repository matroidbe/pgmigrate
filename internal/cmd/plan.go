package cmd

import (
	"fmt"
	"os"

	"github.com/matroidbe/pgmigrate/internal/db"
	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/spf13/cobra"
)

var (
	planOutput string
)

var planCmd = &cobra.Command{
	Use:   "plan [file]",
	Short: "Show migration plan without applying",
	Long: `Compares the schema.yaml file against the live database and shows
what changes would be made. No changes are applied.

The plan output categorizes changes by safety level:
  + Safe:        Additive changes (auto-applied)
  - Destructive: Data loss possible (requires --allow-destructive)
  ! Breaking:    May fail or corrupt (requires manual dba_migrate)

Examples:
  pgmigrate plan                    # Use schema.yaml in current directory
  pgmigrate plan myschema.yaml      # Use specific file
  pgmigrate plan -o json            # Output as JSON`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPlan,
}

func init() {
	planCmd.Flags().StringVarP(&planOutput, "output", "o", "text",
		"Output format: text, json")
}

func runPlan(cmd *cobra.Command, args []string) error {
	// Determine schema file
	schemaFile := "schema.yaml"
	if len(args) > 0 {
		schemaFile = args[0]
	}

	// Read YAML content
	yamlContent, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", schemaFile, err)
	}

	// Connect to database
	conn, err := db.Connect(getDatabaseURL())
	if err != nil {
		return err
	}
	defer conn.Close(cmd.Context())

	// Check extension
	if err := db.CheckExtension(conn); err != nil {
		return err
	}

	// Get plan
	plan, err := db.Plan(conn, string(yamlContent))
	if err != nil {
		return err
	}

	// Output based on format
	if planOutput == "json" {
		return output.PrintPlanJSON(plan)
	}

	output.PrintPlanTerraform(plan)
	return nil
}
