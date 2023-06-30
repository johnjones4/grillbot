package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type message struct {
	Message string `json:"message"`
	Ok      bool   `json:"bool"`
}

func sendJSON(log *logrus.Logger, w http.ResponseWriter, status int, info any) {
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func sendError(log *logrus.Logger, w http.ResponseWriter, status int, err error) {
	sendJSON(log, w, status, message{
		Message: err.Error(),
		Ok:      false,
	})
}

func readJSON(r *http.Request, target any) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, target)
	if err != nil {
		return err
	}

	return nil
}
