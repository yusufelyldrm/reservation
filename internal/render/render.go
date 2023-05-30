package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer set the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders template using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	//create a buffer to write to
	buf := new(bytes.Buffer)

	//add default data to template
	td = AddDefaultData(td, r)

	//execute the template
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	//render the template
	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil

	/*ParsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.html")
	err = ParsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template", err)
		return
	}*/

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.gohtml
	for _, page := range pages {

		//get the file name
		name := filepath.Base(page)

		//ParseFiles returns a *Template associated with the name
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		//get all the files named *.layout.gohtml from ./templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		//range through all files ending with *.layout.gohtml
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		//add the template to the cache
		myCache[name] = ts
	}
	//return the cache
	return myCache, nil
}

/*
//This is the simple one

var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	var gohtml *template.Template
	var err error

	//check to see if we already have the template in our cache
	_, inMap := tc[t]
	if !inMap {
		log.Println("creating template and adding to cache ")
		//need to create the template
		err = createTemplateCache(t)

		if err != nil {
			log.Println(err)
		}

	} else {
		//we have the template in cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}

}
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.html",
	}

	//parse the template
	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		return err
	}

	//add template to cache(map)
	tc[t] = tmpl

	return nil

}
*/
