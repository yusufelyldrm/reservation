package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/yusufelyldrm/reservation/pkg/config"
	"github.com/yusufelyldrm/reservation/pkg/handlers"
	"github.com/yusufelyldrm/reservation/pkg/render"
	"net/http"
	"time"

	"log"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s\n Press 'Ctrl + C' to stop", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	//err = http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
