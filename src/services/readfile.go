package services

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Readfile struct {
	Filename string
}

func (readfile *Readfile) ReadFile() string {
	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + "/../../files/text/"

	content, err := os.ReadFile(absolute_path + readfile.Filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(content)
}