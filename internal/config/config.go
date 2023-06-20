package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/yusufelyldrm/reservation/internal/models"
	"html/template"
	"log"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template // TemplateCache is a map of cached templates
	UseCache      bool                          // UseCache is true when app is in production mode
	InProduction  bool                          // InProduction is true when app is in production mode
	InfoLog       *log.Logger                   // InfoLog is a logger dedicated to logging info messages
	ErrorLog      *log.Logger                   // ErrorLog is a logger dedicated to logging error messages
	Session       *scs.SessionManager           // Session is a session manager
	MailChan      chan models.MailData
}
