package main

import (
	"errors"
	"io/ioutil"
	"log"
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

// Renamer Renamer
type Renamer struct {
	FilePaths []*FilePath
	MODE      string
	MyName    string
}

func (r Renamer) Len() int           { return len(r.FilePaths) }
func (r Renamer) Less(i, j int) bool { return r.FilePaths[i].Depth > r.FilePaths[j].Depth }
func (r Renamer) Swap(i, j int)      { r.FilePaths[i], r.FilePaths[j] = r.FilePaths[j], r.FilePaths[i] }

// ReadFiles wrapper for ReadDirRec and Sort
func (r *Renamer) ReadFiles() {
	r.ReadDirRec(startPath, 0)
	/* NOTE:
	Sort sorts acquired paths by it's depth in DESC order.
	Without this, it's impossible to change directory name.
	When try to change child path name, it may not exist if it's parent directory name was changed before.
	So make sure rename children files first, and then change it's directory name then.
	*/
	renamer.Sort() // <- Must be executed before rename.
}

// Sort sorts the filepaths by it's depth.
func (r *Renamer) Sort() {
	sort.Sort(r)
}

// ReadDirRec reads files recursively.
func (r *Renamer) ReadDirRec(prefix string, depth int) {
	files, err := ioutil.ReadDir(prefix)
	if err != nil {
		log.Fatal(err)
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
}

// Rename Rename
func (r *Renamer) Rename() (success int, fail int, err error) {
	errCount := 0
	for _, path := range r.FilePaths {

		if r.exclude(path.FileName) {
			continue
		}

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

func (r *Renamer) exclude(target string) bool {
	if target == r.MyName ||
		target == r.MyName+".exe" ||
		target == r.MyName+".bat" {
		return true
	}
	return false
}
