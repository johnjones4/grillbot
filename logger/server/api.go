package server

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type API struct {
	session core.Session
	device  core.Device
	log     *logrus.Logger
	Mux     *chi.Mux
}

func New(session core.Session, device core.Device, log *logrus.Logger) *API {
	api := &API{
		session: session,
		device:  device,
		log:     log,
	}

	api.Mux = chi.NewRouter()

	api.Mux.Use(middleware.RequestID)
	api.Mux.Use(middleware.RealIP)
	api.Mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log, NoColor: false}))
	api.Mux.Use(middleware.Recoverer)

	api.Mux.Route("/api", func(r chi.Router) {
		r.Route("/metadata", func(r chi.Router) {
			r.Get("/", api.getMetadata)
			r.Put("/", api.putMetadata)
		})

		r.Get("/chart.png", api.getChart)
		r.Get("/stream", api.stream)
	})

	fs := http.FileServer(http.Dir("./static"))
	api.Mux.Handle("/*", fs)

	return api
}
