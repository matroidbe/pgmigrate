package cmd

import (
	"fmt"
	"os"

	"github.com/matroidbe/pgmigrate/internal/db"
	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/spf13/cobra"
)

var (
	dumpOutput string
)

var dumpCmd = &cobra.Command{
	Use:   "dump <schema> [schemas...]",
	Short: "Export database schema as YAML",
	Long: `Introspects the database and generates a YAML schema definition.

This is useful for:
  - Importing an existing database into pgmigrate management
  - Reviewing the current state before making changes
  - Generating documentation

Examples:
  pgmigrate dump public                    # Dump public schema
  pgmigrate dump public api                # Dump multiple schemas
  pgmigrate dump public -o schema.yaml     # Write to file`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDump,
}

func init() {
	dumpCmd.Flags().StringVarP(&dumpOutput, "output", "o", "-",
		"Output file (- for stdout)")
}

func runDump(cmd *cobra.Command, args []string) error {
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

	// Dump schemas
	yaml, err := db.Dump(conn, args)
	if err != nil {
		return err
	}

	// Output
	if dumpOutput == "-" {
		fmt.Print(yaml)
		return nil
	}

	if err := os.WriteFile(dumpOutput, []byte(yaml), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", dumpOutput, err)
	}

	output.PrintSuccess(fmt.Sprintf("Schema written to %s", dumpOutput))
	return nil
}
