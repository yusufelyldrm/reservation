package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// renderTemplate renders templates using html/template
func renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}
}
