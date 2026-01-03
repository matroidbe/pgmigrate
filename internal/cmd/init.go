package cmd

import (
	"fmt"
	"os"

	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/matroidbe/pgmigrate/internal/schema"
	"github.com/spf13/cobra"
)

var (
	initForce bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a template schema.yaml file",
	Long: `Creates a template schema.yaml file in the current directory.

If a schema already exists in the database, consider using 'pgmigrate dump'
to generate a starting point instead.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false,
		"Overwrite existing schema.yaml")
}

func runInit(cmd *cobra.Command, args []string) error {
	filename := "schema.yaml"

	// Check if file exists
	if _, err := os.Stat(filename); err == nil && !initForce {
		return fmt.Errorf("schema.yaml already exists. Use --force to overwrite")
	}

	// Write template
	if err := os.WriteFile(filename, []byte(schema.DefaultTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create schema.yaml: %w", err)
	}

	output.PrintSuccess("Created schema.yaml")
	fmt.Println()
	fmt.Println("Edit this file to define your desired schema, then run:")
	fmt.Println("  pgmigrate plan    # Preview changes")
	fmt.Println("  pgmigrate apply   # Apply changes")

	return nil
}
