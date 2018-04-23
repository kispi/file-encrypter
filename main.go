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

func showHelp() {
	helpers.Printf(color.FgWhite, "\n    2018-04-21, kispi@naver.com\n\n")
	helpers.Printf(color.FgWhite, "    This program encrypts all filenames(or files) in specified path using crypto/cipher.\n")
	helpers.Printf(color.FgWhite, "    This doesn't change the content of file, but just name.\n")
	helpers.Printf(color.FgWhite, "    That is, even after the filename has been changed,\n")
	helpers.Printf(color.FgWhite, "    it can be executed if opened with proper application.\n\n")
	helpers.Printf(color.FgRed, "    CAUTION: Do not decrypt(-d) before you encrypt(-e) by this program\n")
	helpers.Printf(color.FgRed, "    since there can be some filenames that are already encrypted by AES.\n\n")
	helpers.Printf(color.FgCyan, "    OPTIONS:\n\n")
	helpers.Printf(color.FgCyan, "    -e: Encrypt\n")
	helpers.Printf(color.FgCyan, "    -d: Decrypt\n")
	helpers.Printf(color.FgCyan, "    -p: Path\n")
	helpers.Printf(color.FgCyan, "    -k: Custom key.(Will be truncated as proper length and shown after execution.) [NOT_IMPLEMENTED]\n")
	helpers.Printf(color.FgCyan, "    -f: File. This will encrypt the contents of the file. They can't be opened before decryption. [NOT_IMPLEMENTED]\n")
	helpers.Printf(color.FgCyan, "    [-h]: Show help\n\n")
	helpers.Printf(color.FgWhite, "    EX:) [BINARY] -e -p ./ (encrypt all filenames in current path)\n")
}
