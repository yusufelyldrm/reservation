package render

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// RenderTemplateTest renders templates using html/template
func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	// parse the template files
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.gohtml")
	// execute the template
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}
}

var tc = make(map[string]*template.Template)

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	//check if template have  already in our cache
	_, inMap := tc[t]
	if !inMap {
		//we need to create the
		log.Println("Creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		//we use the template from the cache
		log.Println("Using template from cache")
	}

	tmpl = tc[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

// createTemplateCache creates a template cache as a map
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.gohtml",
	}
	//parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	//add the template to the cache
	tc[t] = tmpl
	return nil

}
