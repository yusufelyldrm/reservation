package render

import (
	"bytes"
	"fmt"
	"github.com/yusufelyldrm/reservation/pkg/config"
	"net/http"
	"path/filepath"
	"text/template"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	var tc map[string]*template.Template
	if app.UseCache {
		//get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		fmt.Println("Could'nt get template from template cache")
		return
	}

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, nil)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser: ", err)
		return
	}

	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.gohtml")
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all the files named *.page.gohtml from ./templates
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		fmt.Println("Error getting pages: ", err)
		return myCache, err
	}

	//range through the all files ending with *.page.gohtml
	for _, page := range pages {
		//get the file namae
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			fmt.Println("Error parsing template: ", err)
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			fmt.Println("Error getting layout files: ", err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				fmt.Println("Error parsing layout files: ", err)
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

/*
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
*/
