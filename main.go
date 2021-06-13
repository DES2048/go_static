package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func ListMp4(path string) ([]string, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	globPath := filepath.Join(absPath, "*.mp4")

	return filepath.Glob(globPath)
}

var ListTemplate *template.Template

func PrepareTemplates(_ string) error {
	var err error

	ListTemplate, err = template.New("videos").Parse(listTmpl)

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

	config := DefaultServerConfig()
	config.StaticFolder = os.Args[1]

	server, err := NewServer(config)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())

}
