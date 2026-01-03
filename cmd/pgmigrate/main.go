package main

import (
	"fmt"
	"os"

	"github.com/matroidbe/pgmigrate/internal/cmd"
	"github.com/matroidbe/pgmigrate/internal/output"
)

func main() {
	if err := cmd.Execute(); err != nil {
		output.PrintError(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}
