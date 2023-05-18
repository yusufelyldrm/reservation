package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/handlers"
	"github.com/yusufelyldrm/reservation/internal/models"
	"github.com/yusufelyldrm/reservation/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting application on port %s\n Press 'Ctrl + C' to stop", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

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
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return nil
}
