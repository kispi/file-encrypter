package main

import (
	"os"

	"./constant"
	"./helpers"

	"github.com/fatih/color"
)

var renamer Renamer
var startPath string

func main() {
	renamer.MyName = os.Args[0]
	args := os.Args[1:]
	if len(args) >= 1 {
		options, err := parseCommandLineArguments(args)
		if err != nil {
			helpers.Error(err)
			return
		}

		argument := new(Argument)
		err = argument.validateCommandLineArguments(options)
		if err != nil {
			helpers.Error(err)
			return
		}
	}

	if renamer.MODE == constant.HELP || len(args) < 1 {
		showHelp()
		return
	}

	last := startPath[len(startPath)-1]
	if last != '/' && last != '\\' {
		startPath += string(os.PathSeparator)
	}

	renamer.ReadFiles()

	success, fail, err := renamer.Rename()
	if err != nil {
		helpers.Error(err)
		return
	}

	helpers.Printf(color.FgHiMagenta, "%d / %d files are renamed.", success, success+fail)
}
