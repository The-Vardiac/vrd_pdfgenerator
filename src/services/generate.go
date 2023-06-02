package services

import (
	"path/filepath"
	"runtime"

	"github.com/go-pdf/fpdf"
)

type Generate struct {
	Filename string
	Text string
}

func (generate *Generate) GeneratePDF() error {
	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + "/../../files/pdfs/"

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	_, lineHt := pdf.GetFontSize()
	pdf.Write(lineHt, generate.Text)
	err := pdf.OutputFileAndClose(absolute_path + generate.Filename + ".pdf")
	if err != nil {
		return err
	}
	
	return nil
}