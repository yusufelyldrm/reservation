package main

import (
	"fmt"
	"github.com/yusufelyldrm/reservation/pkg/config"
	"github.com/yusufelyldrm/reservation/pkg/handlers"
	"github.com/yusufelyldrm/reservation/pkg/render"
	"net/http"

	"log"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig
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
