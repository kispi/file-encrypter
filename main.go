package main

import (
	"io/ioutil"
	"os"

	"./constant"
	"./helpers"

	"github.com/fatih/color"
)

var encrypter Encrypter
var startPath string

func init() {
	encrypter.MyName = os.Args[0]
	encrypter.Key = []byte("You can't hack since I made this")
}

func main() {
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

	// If this succeed, that means it's done on file. So should be terminated.
	err := tryAsFile(startPath)
	if err == nil {
		return
	}

	last := startPath[len(startPath)-1]
	if last != '/' && last != '\\' {
		startPath += string(os.PathSeparator)
	}

	err = encrypter.ReadFiles()
	if err != nil {
		helpers.Error(err)
		return
	}

	success, fail, err := encrypter.Encrypt(encrypter.OnFile)
	if err != nil {
		helpers.Error(err)
		return
	}

	if encrypter.MODE == constant.ENCRYPT {
		helpers.Printf(color.FgHiMagenta, "%d / %d files are encrypted.", success, success+fail)
	} else if encrypter.MODE == constant.DECRYPT {
		helpers.Printf(color.FgHiMagenta, "%d / %d files are decrypted.", success, success+fail)
	}

}

func tryAsFile(start string) error {
	// If given path refers to file, not a directory, then do on that file only.
	_, err := ioutil.ReadFile(start)
	if err == nil {
		encrypter.ModifyFile(start)
		if encrypter.MODE == constant.ENCRYPT {
			helpers.Printf(color.FgHiMagenta, "%s is encrypted.", startPath)
		} else if encrypter.MODE == constant.DECRYPT {
			helpers.Printf(color.FgHiMagenta, "%s is decrypted.", startPath)
		}
		return nil
	}
	return err
}

func showHelp() {
	helpers.Printf(color.FgWhite, "\n    2018-04-21, kispi@naver.com\n\n")
	helpers.Printf(color.FgWhite, "    This program encrypts all filenames(or even files) in specified path using AES.\n")
	helpers.Printf(color.FgWhite, "    This doesn't change the content of file, but just name as default.\n")
	helpers.Printf(color.FgWhite, "    That is, even after the filename has been changed, someone can execute it with proper app.\n")
	helpers.Printf(color.FgWhite, "    If you don't want it to happen, encrypt with -f option which also encrypts the content of file.\n\n")
	helpers.Printf(color.FgRed, "    CAUTION: Do not decrypt(-d) before you encrypt(-e) by this program\n")
	helpers.Printf(color.FgRed, "    since there can be some filenames that are already encrypted by AES.\n\n")
	helpers.Printf(color.FgCyan, "    OPTIONS:\n")
	helpers.Printf(color.FgCyan, "    -e: Encrypt [REQUIRED]\n")
	helpers.Printf(color.FgCyan, "    -d: Decrypt [REQUIRED]\n")
	helpers.Printf(color.FgCyan, "    -p: Path [REQUIRED]\n")
	helpers.Printf(color.FgCyan, "    -f: File. This will also encrypt the contents of the file. They can't be opened before decryption.\n")
	helpers.Printf(color.FgCyan, "        Usually, I don't think you have to use this option since only changing name is enough in many cases.\n")
	helpers.Printf(color.FgCyan, "    -k: Custom key.(Will be truncated as proper length and shown after execution.) [NOT_IMPLEMENTED]\n")
	helpers.Printf(color.FgCyan, "    [-h]: Show help\n\n")
	helpers.Printf(color.FgYellow, "    EX:) [BINARY] -e -p ./ (encrypt all filenames in current path)\n")
	helpers.Printf(color.FgYellow, "    EX:) [BINARY] -d -p ./ (decrypt all filenames in current path)\n")
	helpers.Printf(color.FgRed, "    EX:) [BINARY] -e -p ./ -f (encrypt all files in current path)\n")
	helpers.Printf(color.FgRed, "    EX:) [BINARY] -d -p ./ -f (decrypt all files in current path)\n\n")
}
