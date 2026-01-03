package cmd

import (
	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"
	// GitCommit is set at build time
	GitCommit = ""

	// Global flags
	databaseURL string
	verbose     bool
	noColor     bool
)

var rootCmd = &cobra.Command{
	Use:   "pgmigrate",
	Short: "Declarative PostgreSQL schema management",
	Long: `pgmigrate is a Terraform-like tool for managing PostgreSQL schema.

Define your desired schema state in YAML, and pgmigrate computes the diff
against the live database. No internal state file - pg_catalog IS the state.

Requires the pg_migrate extension to be installed in your database.
Connection is configured via DATABASE_URL environment variable.

Example:
  export DATABASE_URL="postgres://user:pass@localhost:5432/mydb"
  pgmigrate init           # Create template schema.yaml
  pgmigrate plan           # Show what would change
  pgmigrate apply          # Apply safe changes
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor {
			output.DisableColors()
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "",
		"PostgreSQL connection URL (overrides DATABASE_URL env)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"Disable colored output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(historyCmd)
}

// getDatabaseURL returns the database URL from flag or environment
func getDatabaseURL() string {
	return databaseURL
}

// isVerbose returns true if verbose mode is enabled
func isVerbose() bool {
	return verbose
}
