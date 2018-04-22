package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"./constant"
	"./helpers"
	"github.com/fatih/color"
)

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
		return errors.New("You must specify the path.(-p [PATH])")
	}

	if help && (decrypt || encrypt || path) {
		return errors.New("-h cannot be used with other parameters")
	}

	if help {
		renamer.MODE = constant.HELP
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
	helpers.Printf(color.FgWhite, "\n    2018-04-21, kispi@naver.com\n\n")
	helpers.Printf(color.FgWhite, "    This encrypts all filenames in specified path using crypto/cipher.(%s)\n", renamer.MyName)
	helpers.Printf(color.FgCyan, "\n    OPTIONS:\n\n    -e: ENCRYPT\n    -d: DECRYPT\n    -h: SHOW HELP\n    -p: PATH\n\n")
	helpers.Printf(color.FgWhite, "    EX:) renamer -e -p ./ (encrypt all filenames in current path)\n")
}
