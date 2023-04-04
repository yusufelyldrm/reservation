package main

import (
	"github.com/bmizerany/pat"
	"github.com/yusufelyldrm/reservation/pkg/config"
	"github.com/yusufelyldrm/reservation/pkg/handlers"
	"net/http"
)

func routes(config config.AppConfig) http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	return mux
}
