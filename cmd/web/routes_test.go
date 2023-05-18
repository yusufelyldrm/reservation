package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yusufelyldrm/reservation/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := Routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux ,but is %T", v))
	}
}
