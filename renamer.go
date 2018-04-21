package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"./constant"
	"./helpers"
	"github.com/fatih/color"
)

// FilePath FilePath
type FilePath struct {
	FileName   string
	Prefix     string
	EntirePath string
	IsDir      bool
}

// Renamer Renamer
type Renamer struct {
	FilePaths []*FilePath
	MODE      string
}

// ReadDirRec reads files recursively.
func (r *Renamer) ReadDirRec(prefix string) {
	files, err := ioutil.ReadDir(prefix)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filePath := new(FilePath)
		filePath.FileName = f.Name()
		filePath.Prefix = prefix
		filePath.EntirePath = prefix + f.Name()

		if f.IsDir() {
			filePath.IsDir = true
		} else {
			filePath.IsDir = false
		}

		r.FilePaths = append(r.FilePaths, filePath)

		if f.IsDir() {
			r.ReadDirRec(filePath.EntirePath + "/")
		}
	}
}

// Rename Rename
func (r *Renamer) Rename() (success int, fail int, err error) {
	errCount := 0
	for _, path := range r.FilePaths {
		// Since this doesn't change names in recursive way.
		// Rather, this renamer stores all directories first and change them.
		if path.IsDir {
			continue
		}

		if !strings.Contains(path.FileName, ".bat") && !strings.Contains(path.FileName, ".exe") {
			modifiedName, err := r.getNewName(path.FileName)
			if err != nil {
				errCount++
				continue
			}

			dir := filepath.Dir(path.EntirePath)
			base := filepath.Base(path.EntirePath)
			oldPath := filepath.Join(dir, base)
			newPath := filepath.Join(dir, strings.Replace(base, path.FileName, modifiedName, 1))

			oldAbs, _ := filepath.Abs(oldPath)
			newAbs, _ := filepath.Abs(newPath)

			err = os.Rename(oldAbs, newAbs)
			if err != nil {
				helpers.Error(err)
				fail++
			} else {
				success++
			}

			helpers.Printf(color.FgHiBlue, "OLD: %s\n", oldAbs)
			helpers.Printf(color.FgHiGreen, "NEW: %s\n\n", newAbs)
		}
	}

	return success, fail, nil
}

func (r *Renamer) getNewName(filename string) (modifiedName string, err error) {
	key := []byte("You can't hack since I made this")
	if r.MODE == constant.ENCRYPT {
		modifiedName, err = helpers.Encrypt(key, filename)
		if err != nil {
			return "", err
		}
	} else if r.MODE == constant.DECRYPT {
		modifiedName, err = helpers.Decrypt(key, filename)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("Rename mode is not set")
	}
	return modifiedName, nil
}
