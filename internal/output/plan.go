package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/matroidbe/pgmigrate/internal/db"
)

// PrintPlanTerraform outputs a Terraform-like migration plan
func PrintPlanTerraform(plan *db.PlanResult) {
	if plan.IsEmpty() {
		fmt.Println(Green("No changes.") + " Your schema matches the database.")
		return
	}

	fmt.Println(Bold("pgmigrate will perform the following actions:"))
	fmt.Println()

	// Safe changes (green +)
	for _, change := range plan.Safe {
		printChange(PlusSymbol, Green, change)
	}

	// Destructive changes (red -)
	for _, change := range plan.Destructive {
		printChange(MinusSymbol, Red, change)
	}

	// Breaking changes (yellow !)
	for _, change := range plan.Breaking {
		printChange(BangSymbol, Yellow, change)
	}

	fmt.Println()
	printSummary(plan)
}

func printChange(symbol string, colorFn func(...interface{}) string, change db.Change) {
	switch change.ChangeType {
	case "CreateSchema":
		fmt.Printf("  %s %s\n", symbol, colorFn(fmt.Sprintf("CREATE SCHEMA %s", change.Schema)))

	case "DropSchema":
		fmt.Printf("  %s %s\n", symbol, colorFn(fmt.Sprintf("DROP SCHEMA %s", change.Schema)))

	case "CreateTable":
		fmt.Printf("  %s %s\n", symbol, colorFn(fmt.Sprintf("CREATE TABLE %s.%s", change.Schema, change.Table)))

	case "DropTable":
		fmt.Printf("  %s %s\n", symbol, colorFn(fmt.Sprintf("DROP TABLE %s.%s", change.Schema, change.Table)))

	case "AddColumn":
		fmt.Printf("  %s %s.%s.%s %s\n", symbol,
			change.Schema, change.Table, colorFn(change.Column), Faint("(add column)"))

	case "DropColumn":
		fmt.Printf("  %s %s.%s.%s %s\n", symbol,
			change.Schema, change.Table, colorFn(change.Column), Faint("(drop column)"))

	case "AlterColumnType":
		typeInfo := ""
		if change.OldType != nil && change.NewType != nil {
			typeInfo = fmt.Sprintf(" %s -> %s", *change.OldType, *change.NewType)
		}
		fmt.Printf("  %s %s.%s.%s %s%s\n", symbol,
			change.Schema, change.Table, colorFn(change.Column), Faint("(alter type)"), typeInfo)

	case "AlterColumnNullable":
		fmt.Printf("  %s %s.%s.%s %s\n", symbol,
			change.Schema, change.Table, colorFn(change.Column), Faint("(alter nullable)"))

	case "AlterColumnDefault":
		fmt.Printf("  %s %s.%s.%s %s\n", symbol,
			change.Schema, change.Table, colorFn(change.Column), Faint("(alter default)"))

	case "CreateIndex":
		fmt.Printf("  %s %s %s\n", symbol, colorFn(fmt.Sprintf("CREATE INDEX %s", change.Index)), Faint(fmt.Sprintf("ON %s.%s", change.Schema, change.Table)))

	case "DropIndex":
		fmt.Printf("  %s %s %s\n", symbol, colorFn(fmt.Sprintf("DROP INDEX %s", change.Index)), Faint(fmt.Sprintf("ON %s.%s", change.Schema, change.Table)))

	default:
		// Fallback for unknown change types
		if change.Description != "" {
			fmt.Printf("  %s %s\n", symbol, colorFn(change.Description))
		} else {
			fmt.Printf("  %s %s\n", symbol, colorFn(change.ChangeType))
		}
	}
}

func printSummary(plan *db.PlanResult) {
	parts := []string{}

	if plan.SafeCount() > 0 {
		parts = append(parts, Green(fmt.Sprintf("%d to add", plan.SafeCount())))
	}
	if plan.DestructiveCount() > 0 {
		parts = append(parts, Red(fmt.Sprintf("%d to destroy", plan.DestructiveCount())))
	}
	if plan.BreakingCount() > 0 {
		parts = append(parts, Yellow(fmt.Sprintf("%d breaking", plan.BreakingCount())))
	}

	fmt.Printf("Plan: %s.\n", strings.Join(parts, ", "))

	if plan.BreakingCount() > 0 {
		fmt.Println()
		fmt.Println(Yellow("Warning:") + " Breaking changes require manual intervention.")
		fmt.Println("Run SQL directly or use pgmigrate.dba_migrate() in psql.")
	}
}

// PrintPlanJSON outputs the plan as JSON
func PrintPlanJSON(plan *db.PlanResult) error {
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// ConfirmPrompt asks user for confirmation
func ConfirmPrompt(message string) bool {
	fmt.Printf("%s [y/N]: ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// PrintApplyResult shows the result of apply
func PrintApplyResult(result *db.ApplyResult) {
	if len(result.Applied) == 0 {
		fmt.Println(Yellow("No changes applied."))
		return
	}

	fmt.Println()
	fmt.Printf("%s Applied %d change(s) in %dms.\n",
		Green("Apply complete!"), len(result.Applied), result.DurationMs)

	if len(result.Skipped) > 0 {
		fmt.Printf("%s Use --allow-destructive to include.\n",
			Yellow(fmt.Sprintf("Skipped %d destructive change(s).", len(result.Skipped))))
	}
}

// PrintError prints an error message in red
func PrintError(message string) {
	fmt.Println(Red("Error:") + " " + message)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(message string) {
	fmt.Println(Yellow("Warning:") + " " + message)
}

// PrintSuccess prints a success message in green
func PrintSuccess(message string) {
	fmt.Println(Green("Success:") + " " + message)
}

// PrintHistoryTable prints migration history in a table format
func PrintHistoryTable(entries []db.HistoryEntry) {
	if len(entries) == 0 {
		fmt.Println("No migration history found.")
		return
	}

	fmt.Println(Bold("Migration History"))
	fmt.Println()
	fmt.Printf("%-6s %-12s %-12s %-24s %-20s %s\n",
		"ID", "TYPE", "HASH", "APPLIED AT", "BY", "DURATION")
	fmt.Println(strings.Repeat("-", 90))

	for _, e := range entries {
		typeColor := Green
		if e.MigrationType == "dba_migrate" {
			typeColor = Yellow
		}

		hash := e.YAMLHash
		if len(hash) > 10 {
			hash = hash[:10] + "..."
		}

		appliedAt := e.AppliedAt
		if len(appliedAt) > 22 {
			appliedAt = appliedAt[:22]
		}

		fmt.Printf("%-6d %s %-12s %-24s %-20s %dms\n",
			e.ID,
			typeColor(fmt.Sprintf("%-12s", e.MigrationType)),
			hash,
			appliedAt,
			e.AppliedBy,
			e.DurationMs)
	}
}
