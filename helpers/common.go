package helpers

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

// Printf does fmt.Printf with color
func Printf(c color.Attribute, format string, arg ...interface{}) {
	color.Set(c)
	fmt.Printf(format, arg...)
	color.Unset()
}

// Error prints error message with red color
func Error(err error) {
	_, fn, line, _ := runtime.Caller(1)
	Printf(color.FgRed, "[ERROR] %s:%d %v\n", fn, line, err.Error())
}
