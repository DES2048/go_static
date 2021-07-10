package mp4server

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

type EchoServer struct {
	Config *ServerConfig
	e      *echo.Echo
}

func NewEchoServer(config *ServerConfig) (*EchoServer, error) {
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

		videosTmpl := make([]VideoData, 0, len(videos))

		for _, v := range videos {
			videosTmpl = append(videosTmpl,
				VideoData{
					Url:      path.Join(config.StaticPath, v.Name()),
					Title:    v.Name(),
					Size:     humanize.Bytes(uint64(v.Size())),
					Time:     humanize.Time(v.ModTime()),
					Duration: v.Duration.String(),
				},
			)
		}

		data := &struct {
			Videos []VideoData
		}{
			Videos: videosTmpl,
		}

		err = config.ListTemplate.Execute(c.Response(), data)
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

	return &EchoServer{
		Config: config,
		e:      e,
	}, nil
}

func (s *EchoServer) Start() error {
	return s.e.Start(s.Config.Addr)
}
