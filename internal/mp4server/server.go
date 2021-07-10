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

type VideoData struct {
	Url      string
	Title    string
	Size     string
	Time     string
	Duration string
}

type Server interface {
	Start() error
}
