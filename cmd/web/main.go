package main

import (
	"encoding/gob"
	"fmt"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/handlers"
	"github.com/yusufelyldrm/reservation/internal/models"
	"github.com/yusufelyldrm/reservation/internal/render"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"log"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	//what am I going to put in the session
	gob.Register(models.Reservation{})

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

	fmt.Printf(fmt.Sprintf("Starting application on port %s\n Press 'Ctrl + C' to stop", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
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
