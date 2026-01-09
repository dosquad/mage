package loga

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
)

// PrintWarningf prints the passed warning message to stdout in white text with yellow chevron.
//
//nolint:forbidigo // printing output
func PrintWarningf(format string, v ...any) {
	fmt.Printf("%s %s\n", color.YellowString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintInfof prints the passed info message to stdout in white text with green chevron.
//
//nolint:forbidigo // printing output
func PrintInfof(format string, v ...any) {
	fmt.Printf("%s %s\n", color.GreenString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintCommandf prints the passed command line to stdout in white text with magenta chevron
// when Verbose or Debug is enabled.
func PrintCommandf(format string, v ...any) {
	if mg.Verbose() || mg.Debug() {
		PrintCommandAlwaysf(format, v...)
	}
}

// PrintCommandAlwaysf prints the passed command line to stdout in white text with magenta chevron.
//
//nolint:forbidigo // printing output
func PrintCommandAlwaysf(format string, v ...any) {
	fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf(format, v...))
}

// PrintDebugf prints the passed debug message to stdout in white text with blue chevron.
//
//nolint:forbidigo // printing output
func PrintDebugf(format string, v ...any) {
	if mg.Verbose() || mg.Debug() {
		fmt.Printf("%s %s\n", color.BlueString(">"), color.New(color.Bold).Sprintf(format, v...))
	}
}

// PrintFileUpdatef prints the passed message to stdout in white text with magenta chevron.
//
//nolint:forbidigo // printing output
func PrintFileUpdatef(format string, v ...any) {
	fmt.Printf("%s %s\n", color.HiBlueString(">"), color.New(color.Bold).Sprintf(format, v...))
}
