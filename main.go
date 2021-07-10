package main

import (
	"embed"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func ListMp4(path string) ([]os.FileInfo, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	globPath := filepath.Join(absPath, "*.mp4")

	paths, err := filepath.Glob(globPath)

	if err != nil {
		return nil, err
	}

	infos := make([]os.FileInfo, 0, len(paths))

	// gather file info's
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

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

	config := DefaultServerConfig()
	config.StaticFolder = os.Args[1]

	server, err := NewServer(config)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())

}
