package output

import (
	"github.com/fatih/color"
)

var (
	// Color functions
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
	Faint  = color.New(color.Faint).SprintFunc()

	// Prefix symbols
	PlusSymbol  = Green("+")
	MinusSymbol = Red("-")
	BangSymbol  = Yellow("!")
)

// DisableColors disables all colored output
func DisableColors() {
	color.NoColor = true
}

// EnableColors enables colored output
func EnableColors() {
	color.NoColor = false
}
