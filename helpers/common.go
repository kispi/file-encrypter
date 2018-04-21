package helpers

import (
	"fmt"

	"github.com/fatih/color"
)

// Printf does fmt.Printf with color
func Printf(c color.Attribute, format string, arg ...interface{}) {
	color.Set(c)
	fmt.Printf(format, arg...)
	color.Unset()
}

// Error prints error message with red font color
func Error(err error) {
	Printf(color.FgHiRed, "%s\n", err.Error())
}
