package mp4server

import (
	"html/template"
)

type ServerConfig struct {
	Addr         string
	StaticPath   string
	StaticFolder string
	ListTemplate *template.Template
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:       ":8500",
		StaticPath: "/videos",
	}
}

type Server interface {
	Start() error
}
