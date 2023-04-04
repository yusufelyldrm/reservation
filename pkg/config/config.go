package config

import (
	"github.com/alexedwards/scs/v2"
	"log"
	"text/template"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	InfoLog       *log.Logger
	Session       *scs.SessionManager
}
