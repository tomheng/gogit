package file

import (
	"os"
	"path"
)

//check file exist
func IsExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//create dir
func MakeDir(dir string) (err error) {
	if IsExist(dir) {
		return
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	return
}

func Create(name string) (file *os.File, err error) {
	dir := path.Dir(name)
	err = MakeDir(dir)
	if err != nil {
		return
	}
	return os.Create(name)
}

//wrap os OpenFile
func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	dir := path.Dir(name)
	err = MakeDir(dir)
	if err != nil {
		return
	}
	return os.OpenFile(name, flag, perm)
}
