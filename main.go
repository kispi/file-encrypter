package main

import (
	"os"

	"./constant"
	"./helpers"

	"github.com/fatih/color"
)

var encrypter Encrypter
var startPath string

func main() {
	encrypter.MyName = os.Args[0]
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

	if encrypter.MODE == constant.HELP || len(args) < 1 {
		showHelp()
		return
	}

	last := startPath[len(startPath)-1]
	if last != '/' && last != '\\' {
		startPath += string(os.PathSeparator)
	}

	encrypter.ReadFiles()

	success, fail, err := encrypter.Encrypt()
	if err != nil {
		helpers.Error(err)
		return
	}

	if encrypter.MODE == constant.ENCRYPT {
		helpers.Printf(color.FgHiMagenta, "%d / %d filenames are encrypted.", success, success+fail)
	} else if encrypter.MODE == constant.DECRYPT {
		helpers.Printf(color.FgHiMagenta, "%d / %d filenames are decrypted.", success, success+fail)
	}

}
