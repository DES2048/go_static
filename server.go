package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

		// sort by modTime
		sort.Sort(
			sort.Reverse(Mp4ByModTime(videos)),
		)

		type Video struct {
			Url   string
			Title string
			Size  string
			Time  string
		}

		videosTmpl := make([]Video, 0, len(videos))

		for _, v := range videos {
			videosTmpl = append(videosTmpl,
				Video{
					Url:   path.Join(config.StaticPath, v.Name()),
					Title: v.Name(),
					Size:  humanize.Bytes(uint64(v.Size())),
					Time:  humanize.Time(v.ModTime()),
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

	e.DELETE("/d/:file", func(c echo.Context) error {
		filePar, err := url.PathUnescape(c.Param("file"))

		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid path parameter",
			)
		}

		delPath := filepath.Join(config.StaticFolder, filePar)

		if err := os.Remove(delPath); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
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
