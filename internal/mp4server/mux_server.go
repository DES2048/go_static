package mp4server

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type MuxServer struct {
	config *ServerConfig
	router *mux.Router
}

func (s *MuxServer) Start() error {
	return http.ListenAndServe(s.config.Addr, s.router)
}

func NewMuxServer(config *ServerConfig) (*MuxServer, error) {

	server := &MuxServer{
		config: config,
		router: mux.NewRouter(),
	}

	// setup handlers
	server.router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		RenderMp4ListTemplate(
			config,
			rw,
			"",
		)
	}).Methods("GET")

	server.router.HandleFunc("/d/:file", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		filePar, err := url.PathUnescape(vars["id"])

		if err != nil {
			return
		}

		delPath := filepath.Join(config.StaticFolder, filePar)

		if err := os.Remove(delPath); err != nil {
			return
		}

	}).Methods("DELETE")

	server.router.PathPrefix(config.StaticPath).
		Handler(http.StripPrefix(
			config.StaticPath,
			http.FileServer(http.Dir(config.StaticFolder))))

	return server, nil
}
