package loga

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
)

// PrintWarning prints the passed warning message to stdout in white text with yellow chevron.
//
//nolint:forbidigo // printing output
func PrintWarning(format string, v ...any) {
	fmt.Printf("%s %s\n", color.YellowString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintInfo prints the passed info message to stdout in white text with green chevron.
//
//nolint:forbidigo // printing output
func PrintInfo(format string, v ...any) {
	fmt.Printf("%s %s\n", color.GreenString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintCommand prints the passed command line to stdout in white text with magenta chevron
// when Verbose or Debug is enabled.
func PrintCommand(format string, v ...any) {
	if mg.Verbose() || mg.Debug() {
		PrintCommandAlways(format, v...)
	}
}

// PrintCommandAlways prints the passed command line to stdout in white text with magenta chevron.
//
//nolint:forbidigo // printing output
func PrintCommandAlways(format string, v ...any) {
	fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintDebug prints the passed debug message to stdout in white text with blue chevron.
//
//nolint:forbidigo // printing output
func PrintDebug(format string, v ...any) {
	if mg.Verbose() || mg.Debug() {
		fmt.Printf("%s %s\n", color.BlueString(">"), color.New(color.Bold).Sprintf(format, v...))
	}
}

// PrintFileUpdate prints the passed message to stdout in white text with magenta chevron.
//
//nolint:forbidigo // printing output
func PrintFileUpdate(format string, v ...any) {
	fmt.Printf("%s %s\n", color.HiBlueString(">"), color.New(color.Bold).Sprintf(format, v...))
}
