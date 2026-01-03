package main

import (
	"os"

	"github.com/matroidbe/pgmigrate/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
