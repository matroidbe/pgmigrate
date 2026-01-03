package cmd

import (
	"github.com/matroidbe/pgmigrate/internal/db"
	"github.com/matroidbe/pgmigrate/internal/output"
	"github.com/spf13/cobra"
)

var (
	historyLimit int
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show migration history",
	Long: `Displays the migration history from the database, showing all
previous migrations applied via pgmigrate.apply() and dba_migrate().

Examples:
  pgmigrate history           # Show last 20 migrations
  pgmigrate history -n 50     # Show last 50 migrations`,
	RunE: runHistory,
}

func init() {
	historyCmd.Flags().IntVarP(&historyLimit, "limit", "n", 20,
		"Number of entries to show")
}

func runHistory(cmd *cobra.Command, args []string) error {
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

	// Get history
	entries, err := db.GetHistory(conn, historyLimit)
	if err != nil {
		return err
	}

	output.PrintHistoryTable(entries)
	return nil
}
