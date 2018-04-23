package main

import (
	"errors"
	"strings"

	"./constant"
)

type Argument struct {
	Encrypt bool
	Decrypt bool
	Help    bool
	Path    bool
	File    bool
}

// Option command line option
type Option struct {
	Key   string
	Value string
}

func (a *Argument) validateCommandLineArguments(options []*Option) error {

	a.Encrypt = false
	a.Decrypt = false
	a.Help = false
	a.Path = false
	a.File = false

	for _, o := range options {
		switch o.Key {
		case constant.ENCRYPT:
			a.Encrypt = true
		case constant.DECRYPT:
			a.Decrypt = true
		case constant.PATH:
			a.Path = true
			startPath = o.Value
		case constant.HELP:
			a.Help = true
		case constant.FILE:
			a.File = true
		}
	}

	err := a.testArgumentValidity()
	if err != nil {
		return err
	}

	if a.Encrypt {
		encrypter.MODE = constant.ENCRYPT
	} else if a.Decrypt {
		encrypter.MODE = constant.DECRYPT
	}

	if a.File {
		encrypter.OnFile = true
	}

	return nil
}

func (a *Argument) testArgumentValidity() error {
	if !a.Help && !a.Path {
		return errors.New("You must specify the path.(-p [PATH])")
	}

	if a.Help && (a.Decrypt || a.Encrypt || a.Path) {
		return errors.New("-h cannot be used with other parameters")
	}

	if a.Help {
		encrypter.MODE = constant.HELP
	}

	if a.Path && (!a.Decrypt && !a.Encrypt) {
		return errors.New("-p should be used with either -e or -d")
	}

	if a.Encrypt && a.Decrypt {
		return errors.New("-e -d cannot be used together")
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
			option.Key = constant.HELP
		case "-f":
			option.Key = constant.FILE
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
