package outputs

import (
	"context"
	"encoding/json"
	"main/core"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Server struct {
	sess          core.Session
	handler       *chi.Mux
	latestReading core.Reading
	readingMutex  sync.RWMutex
	host          string
}

func NewServer(sess core.Session, host string) *Server {
	server := &Server{
		sess,
		chi.NewRouter(),
		core.Reading{},
		sync.RWMutex{},
		host,
	}

	server.handler.Use(middleware.RequestID)
	server.handler.Use(middleware.RealIP)
	server.handler.Use(middleware.Logger)
	server.handler.Use(middleware.Recoverer)

	server.handler.Get("/", server.indexHandler)

	server.handler.Route("/api", func(r chi.Router) {
		r.Get("/readings/stream", server.streamReadingsHandler)
		r.Get("/readings", server.getReadingsHandler)
		r.Get("/metadata", server.getMetadataHandler)
	})

	return server
}

func (s *Server) receiveUpdates(_ core.Session, r core.Reading) {
	s.readingMutex.Lock()
	s.latestReading = r
	s.readingMutex.Unlock()
}

func (s *Server) Listener() core.Listener {
	return s.receiveUpdates
}

func (s *Server) Start(ctx context.Context) error {
	return http.ListenAndServe(s.host, s.handler) //TODO stop
}

func (s *Server) getReadingsHandler(w http.ResponseWriter, r *http.Request) {
	readings, err := s.sess.GetReadings()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"items": readings,
	})
}

func (s *Server) getMetadataHandler(w http.ResponseWriter, r *http.Request) {
	metadata, err := s.sess.GetMetadata()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusOK, metadata)
}

func (s *Server) streamReadingsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	since := time.Now()

	sinceStr := r.URL.Query().Get("since")
	if sinceStr != "" {
		since, err = time.Parse(time.RFC3339Nano, sinceStr)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}
	}

	upgrader := websocket.Upgrader{} //
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	defer conn.Close()

	var lastSentReading core.Reading
	for {
		var latestReading core.Reading
		s.readingMutex.RLock()
		latestReading = s.latestReading
		s.readingMutex.RUnlock()
		if !latestReading.Received.Equal(lastSentReading.Received) && latestReading.Received.After(since) {
			lastSentReading = latestReading
			bytes, err := json.Marshal(latestReading)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				errorResponse(w, http.StatusInternalServerError, err)
				return
			}
		}
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("./outputs/index.html")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
	}
	w.Header().Set("Content-type", "text/html")
	w.Write(bytes)
}
