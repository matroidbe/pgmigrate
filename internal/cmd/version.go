package cmd

import (
	"fmt"

	"github.com/matroidbe/pgmigrate/internal/db"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and extension versions",
	Long: `Displays the pgmigrate CLI version and, if connected to a database,
the pg_migrate extension version.`,
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	// CLI version
	fmt.Printf("pgmigrate CLI %s", Version)
	if GitCommit != "" {
		if len(GitCommit) > 7 {
			fmt.Printf(" (%s)", GitCommit[:7])
		} else {
			fmt.Printf(" (%s)", GitCommit)
		}
	}
	fmt.Println()

	// Try to get extension version
	conn, err := db.Connect(getDatabaseURL())
	if err != nil {
		fmt.Println("pg_migrate extension: (unable to connect)")
		return nil
	}
	defer conn.Close(cmd.Context())

	extVersion, err := db.GetExtensionVersion(conn)
	if err != nil {
		fmt.Println("pg_migrate extension: not installed")
	} else {
		fmt.Printf("pg_migrate extension: %s\n", extVersion)
	}

	return nil
}
