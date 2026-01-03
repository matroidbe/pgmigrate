package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Change represents a single schema change
type Change struct {
	ChangeType  string  `json:"change_type"`
	Schema      string  `json:"schema,omitempty"`
	Table       string  `json:"table,omitempty"`
	Column      string  `json:"column,omitempty"`
	Index       string  `json:"index,omitempty"`
	Description string  `json:"description,omitempty"`
	SQL         string  `json:"sql,omitempty"`
	Safety      string  `json:"safety"`
	OldType     *string `json:"old_type,omitempty"`
	NewType     *string `json:"new_type,omitempty"`
}

// PlanResult represents the output of pgmigrate.plan()
type PlanResult struct {
	Safe        []Change `json:"safe"`
	Destructive []Change `json:"destructive"`
	Breaking    []Change `json:"breaking"`
}

// HasBreaking returns true if there are breaking changes
func (p *PlanResult) HasBreaking() bool {
	return len(p.Breaking) > 0
}

// HasDestructive returns true if there are destructive changes
func (p *PlanResult) HasDestructive() bool {
	return len(p.Destructive) > 0
}

// IsEmpty returns true if there are no changes
func (p *PlanResult) IsEmpty() bool {
	return len(p.Safe) == 0 && len(p.Destructive) == 0 && len(p.Breaking) == 0
}

// SafeCount returns the number of safe changes
func (p *PlanResult) SafeCount() int {
	return len(p.Safe)
}

// DestructiveCount returns the number of destructive changes
func (p *PlanResult) DestructiveCount() int {
	return len(p.Destructive)
}

// BreakingCount returns the number of breaking changes
func (p *PlanResult) BreakingCount() int {
	return len(p.Breaking)
}

// ApplyResult represents the output of pgmigrate.apply()
type ApplyResult struct {
	Applied    []Change `json:"applied"`
	Skipped    []Change `json:"skipped"`
	DurationMs int      `json:"duration_ms"`
}

// HistoryEntry represents a row from pgmigrate.get_history()
type HistoryEntry struct {
	ID            int    `json:"id"`
	MigrationType string `json:"migration_type"`
	YAMLHash      string `json:"yaml_hash"`
	ExecutedSQL   string `json:"executed_sql"`
	AppliedAt     string `json:"applied_at"`
	AppliedBy     string `json:"applied_by"`
	DurationMs    int    `json:"duration_ms"`
}

// Plan loads YAML and returns the migration plan
func Plan(conn *pgx.Conn, yamlContent string) (*PlanResult, error) {
	ctx := context.Background()

	// Load YAML into session
	var loaded bool
	err := conn.QueryRow(ctx, "SELECT pgmigrate.load($1)", yamlContent).Scan(&loaded)
	if err != nil {
		return nil, fmt.Errorf("load failed: %w", err)
	}

	// Get plan
	var planJSON []byte
	err = conn.QueryRow(ctx, "SELECT pgmigrate.plan()::text").Scan(&planJSON)
	if err != nil {
		return nil, fmt.Errorf("plan failed: %w", err)
	}

	var result PlanResult
	if err := json.Unmarshal(planJSON, &result); err != nil {
		return nil, fmt.Errorf("parse plan failed: %w", err)
	}

	return &result, nil
}

// Apply executes the migration
func Apply(conn *pgx.Conn, yamlContent string, allowDestructive bool) (*ApplyResult, error) {
	ctx := context.Background()

	// Load YAML
	var loaded bool
	if err := conn.QueryRow(ctx, "SELECT pgmigrate.load($1)", yamlContent).Scan(&loaded); err != nil {
		return nil, fmt.Errorf("load failed: %w", err)
	}

	// Apply
	var resultJSON []byte
	err := conn.QueryRow(ctx,
		"SELECT pgmigrate.apply($1)::text", allowDestructive).Scan(&resultJSON)
	if err != nil {
		return nil, fmt.Errorf("apply failed: %w", err)
	}

	var result ApplyResult
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return nil, fmt.Errorf("parse result failed: %w", err)
	}

	return &result, nil
}

// Dump exports schema as YAML
func Dump(conn *pgx.Conn, schemas []string) (string, error) {
	ctx := context.Background()

	var yaml string
	err := conn.QueryRow(ctx, "SELECT pgmigrate.dump($1)", schemas).Scan(&yaml)
	if err != nil {
		return "", fmt.Errorf("dump failed: %w", err)
	}

	return yaml, nil
}

// GetHistory returns migration history
func GetHistory(conn *pgx.Conn, limit int) ([]HistoryEntry, error) {
	ctx := context.Background()

	rows, err := conn.Query(ctx, `
		SELECT id, migration_type, yaml_hash, executed_sql,
		       applied_at::text, applied_by, duration_ms
		FROM pgmigrate.get_history($1)
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("get history failed: %w", err)
	}
	defer rows.Close()

	var entries []HistoryEntry
	for rows.Next() {
		var e HistoryEntry
		if err := rows.Scan(&e.ID, &e.MigrationType, &e.YAMLHash, &e.ExecutedSQL,
			&e.AppliedAt, &e.AppliedBy, &e.DurationMs); err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

// Clear clears the pending session state
func Clear(conn *pgx.Conn) error {
	ctx := context.Background()

	var cleared bool
	if err := conn.QueryRow(ctx, "SELECT pgmigrate.clear()").Scan(&cleared); err != nil {
		return fmt.Errorf("clear failed: %w", err)
	}

	return nil
}
