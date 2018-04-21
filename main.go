package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"./constant"
	"./helpers"

	"github.com/fatih/color"
)

// Option command line option
type Option struct {
	Key   string
	Value string
}

var renamer Renamer
var startPath string

func main() {
	args := os.Args[1:]
	if len(args) >= 1 {
		options, err := parseCommandLineArguments(args)
		if err != nil {
			helpers.Error(err)
			return
		}
		err = validateCommandLineArguments(options)
		if err != nil {
			helpers.Error(err)
			return
		}
	} else {
		showHelp()
		return
	}

	last := startPath[len(startPath)-1]
	if last != '/' && last != '\\' {
		startPath += string(os.PathSeparator)
	}

	renamer.ReadDirRec(startPath)
	success, fail, err := renamer.Rename()
	if err != nil {
		helpers.Error(err)
		return
	}

	helpers.Printf(color.FgHiMagenta, "%d / %d files are renamed.", success, success+fail)
}

func validateCommandLineArguments(options []*Option) error {
	encrypt := false
	decrypt := false
	help := false
	path := false
	for _, o := range options {
		switch o.Key {
		case constant.ENCRYPT:
			encrypt = true
		case constant.DECRYPT:
			decrypt = true
		case constant.PATH:
			path = true
			_, err := ioutil.ReadDir(o.Value)
			if err != nil {
				return errors.New("Cannot access to given path")
			}
			startPath = o.Value
		case constant.HELP:
			help = true
		}
	}

	if !help && !path {
		return errors.New("You must specify path.(-p [PATH])")
	}

	if help && (decrypt || encrypt || path) {
		return errors.New("-h cannot be used with other parameters")
	}

	if encrypt && decrypt {
		return errors.New("-e -d cannot be used together")
	} else if encrypt {
		renamer.MODE = constant.ENCRYPT
	} else if decrypt {
		renamer.MODE = constant.DECRYPT
	}

	return nil
}

func parseCommandLineArguments(args []string) (options []*Option, err error) {
	var curOpt string
	for i := range args {
		option := new(Option)
		opt := strings.ToLower(args[i])
		switch opt {
		case "-e":
			option.Key = constant.ENCRYPT
		case "-d":
			option.Key = constant.DECRYPT
		case "-h":
			option.Key = "HELP"
		case "-p":
			if i+1 < len(args) {
				option.Key = constant.PATH
				option.Value = args[i+1]
			} else {
				err = errors.New("-p must have path as value.(EX: -p ./directoryName)")
			}
		default:
			if curOpt != constant.PATH {
				err = errors.New("non exist argument")
			}
		}

		if err == nil {
			curOpt = option.Key
			options = append(options, option)
		} else {
			return nil, err
		}
	}
	return
}

func showHelp() {
	helpers.Printf(color.FgHiBlue, "2018-04-21, kispi@naver.com\n\n")
	helpers.Printf(color.FgHiWhite, "This encrypts all filename using crypto/cipher except for bat, exe file and directory name.\n")
	helpers.Printf(color.FgCyan, "OPTIONS:\n    -e: ENCRYPT\n    -d: DECRYPT\n    -h: SHOW HELP\n    -p: PATH\n\n")
	helpers.Printf(color.FgGreen, "    EX:) renamer -e -p ./\n")
}
