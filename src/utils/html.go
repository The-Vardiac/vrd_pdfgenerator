package utils

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func HTMLTemplateParser(html string, templateName ...string) string {
	templateNameName := "default"
	if len(templateName) > 0 {
		templateNameName = templateName[0]
	}
	
	// Create a template object
	tmpl, err := template.New(templateNameName).Parse(html)
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		os.Exit(1)
	}

	// Create a buffer to store the rendered HTML
	htmlContent := new(strings.Builder)

	// Render the template into the buffer
	err = tmpl.Execute(htmlContent, nil)
	if err != nil {
		fmt.Println("Failed to render template:", err)
		os.Exit(1)
	}

	return htmlContent.String()
}