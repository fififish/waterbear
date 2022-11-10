package cryptolib

import (
	"fmt"
	"log"
	"os"
	"path"
)

var homepath string

func SetHomeDir() {
	exepath, err := os.Executable()
	if err != nil {
		log.Fatalf("[Cryptolib] Cannot find absolute path of executable")
	}

	p1 := path.Dir(exepath)
	homepath = path.Dir(p1)
}

func GenPath(id int64) string {
	return fmt.Sprintf(homepath+"/etc/key/%d/", id)
}


//Create Dir. Copying from storage to make sdcsm2 library independent
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

//Determines whether folders exist. Copying from storage to make smcgo library independent
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

