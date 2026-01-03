package cmd

import (
	"fmt"
	"os"

	"github.com/matroidbe/pgmigrate/internal/db"
	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/spf13/cobra"
)

var (
	allowDestructive bool
	autoApprove      bool
)

var applyCmd = &cobra.Command{
	Use:   "apply [file]",
	Short: "Apply schema changes to database",
	Long: `Applies safe changes from schema.yaml to the database.

By default, only safe (additive) changes are applied. Destructive changes
(DROP operations) require the --allow-destructive flag.

Breaking changes (type alterations, adding NOT NULL) cannot be applied
automatically. Use pgmigrate.dba_migrate() in psql for those.

Examples:
  pgmigrate apply                          # Apply safe changes
  pgmigrate apply --allow-destructive      # Include DROP operations
  pgmigrate apply --auto-approve           # Skip confirmation
  pgmigrate apply myschema.yaml            # Use specific file`,
	Args: cobra.MaximumNArgs(1),
	RunE: runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&allowDestructive, "allow-destructive", false,
		"Allow destructive changes (DROP TABLE, DROP COLUMN)")
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false,
		"Skip confirmation prompt")
}

func runApply(cmd *cobra.Command, args []string) error {
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

	// Get plan first to show what will happen
	plan, err := db.Plan(conn, string(yamlContent))
	if err != nil {
		return err
	}

	// Handle empty plan
	if plan.IsEmpty() {
		output.PrintPlanTerraform(plan)
		return nil
	}

	// Check for breaking changes
	if plan.HasBreaking() {
		output.PrintPlanTerraform(plan)
		fmt.Println()
		output.PrintError("Breaking changes detected. These cannot be applied automatically.")
		fmt.Println("Use pgmigrate.dba_migrate() in psql to apply breaking changes manually.")
		return fmt.Errorf("breaking changes require manual intervention")
	}

	// Check for destructive without flag
	if plan.HasDestructive() && !allowDestructive {
		output.PrintPlanTerraform(plan)
		fmt.Println()
		output.PrintWarning("Destructive changes will be skipped.")
		fmt.Println("Use --allow-destructive to include them.")
		fmt.Println()
	} else {
		output.PrintPlanTerraform(plan)
		fmt.Println()
	}

	// Skip if only destructive changes and not allowed
	safeCount := plan.SafeCount()
	if safeCount == 0 && !allowDestructive {
		fmt.Println("No safe changes to apply.")
		return nil
	}

	// Confirm unless auto-approve
	if !autoApprove {
		if !output.ConfirmPrompt("Do you want to apply these changes?") {
			fmt.Println("Apply cancelled.")
			return nil
		}
	}

	// Apply changes
	result, err := db.Apply(conn, string(yamlContent), allowDestructive)
	if err != nil {
		return err
	}

	output.PrintApplyResult(result)
	return nil
}
