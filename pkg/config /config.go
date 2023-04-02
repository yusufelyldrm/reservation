package config

import "text/template"

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
}
