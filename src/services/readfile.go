package services

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	textFileAbsolutePath = "/../../files/text/" 
	pdfFileAbsolutePath = "/../../files/pdfs/"
)

type Readfile struct {
	Filename string
}

func (readfile *Readfile) ReadFile() string {
	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + textFileAbsolutePath

	content, err := os.ReadFile(absolute_path + readfile.Filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(content)
}

func (readfile *Readfile) ReadPdfFile() string {
	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + pdfFileAbsolutePath

	content, err := os.ReadFile(absolute_path + readfile.Filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(content)
}

func (readfile *Readfile) GetPdfFilePath() string {
	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + pdfFileAbsolutePath

	file_path := absolute_path + readfile.Filename

	return file_path
}