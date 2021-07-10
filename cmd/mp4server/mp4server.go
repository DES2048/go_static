package main

import (
	"embed"
	"go_static/internal/mp4server"
	"html/template"
	"log"
	"os"
)

//go:embed templates
var fs embed.FS

var ListTemplate *template.Template

func PrepareTemplates(_ string) error {
	var err error

	/*
			ListTemplate, err = template.ParseFiles(
				filepath.Join("templates", "videos-list.html"),
			)
		if _, err := fs.ReadDir("templates"); err != nil {
			log.Fatal(err)
		} */

	ListTemplate, err = template.ParseFS(fs, "templates/videos-list.html")

	if err != nil {
		return err
	}

	return nil
}

func main() {

	if err := PrepareTemplates("."); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("You must specify videos folder")
	}

	staticFolder := os.Args[1]

	if _, err := os.Stat(staticFolder); err != nil {
		log.Fatal(err)
	}

	config := mp4server.DefaultServerConfig()
	config.StaticFolder = os.Args[1]
	config.ListTemplate = ListTemplate

	server, err := mp4server.NewEchoServer(config)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())

}
