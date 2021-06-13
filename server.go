package main

import (
	"path"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var listTmpl string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>List videos</title>
</head>
<body>
    <h2>Videos</h2>
    <ul>
        {{range .Videos}}
            <li>
				<a href="{{.Url}}">{{.Title}}</a>
			</li>
        {{end}}
    </ul> 
</body>
</html>
`

type ServerConfig struct {
	Addr         string
	StaticPath   string
	StaticFolder string
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:       ":8500",
		StaticPath: "/videos",
	}
}

type Server struct {
	Config *ServerConfig
	e      *echo.Echo
}

func NewServer(config *ServerConfig) (*Server, error) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	e.GET("/", func(c echo.Context) error {
		videos, err := ListMp4(config.StaticFolder)

		if err != nil {
			return err
		}
		type Video struct {
			Url   string
			Title string
		}

		videosTmpl := make([]Video, 0, len(videos))

		for _, v := range videos {
			name := filepath.Base(v)
			videosTmpl = append(videosTmpl,
				Video{
					Url:   path.Join(config.StaticPath, name),
					Title: name,
				},
			)
		}

		data := &struct {
			Videos []Video
		}{
			Videos: videosTmpl,
		}

		err = ListTemplate.Execute(c.Response(), data)
		if err != nil {
			return err
		}

		return nil
	})

	e.Static("/videos", config.StaticFolder)

	return &Server{
		Config: config,
		e:      e,
	}, nil
}

func (s *Server) Start() error {
	return s.e.Start(s.Config.Addr)
}
