package server

import (
	"main/core"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (api *API) stream(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sendError(api.log, w, http.StatusInternalServerError, err)
		return
	}

	defer c.Close()

	sinceStr := r.URL.Query().Get("since")
	var since time.Time
	if sinceStr != "" {
		sinceInt, err := strconv.ParseInt(sinceStr, 10, 64)
		if err != nil {
			sendError(api.log, w, http.StatusBadRequest, err)
			return
		}
		since = time.UnixMilli(sinceInt)
	}

	readings, err := api.session.GetReadings()
	if err != nil {
		sendError(api.log, w, http.StatusInternalServerError, err)
		return
	}

	for _, reading := range readings {
		if reading.Received.After(since) {
			err = c.WriteJSON(reading)
			if err != nil {
				api.log.Error(err)
				return
			}
		}
	}

	readings = nil

	updates := make(chan core.Reading, 1024)

	listener := func(_ core.Session, r core.Reading) {
		updates <- r
	}

	i := api.session.AddListener(listener)
	defer api.session.RemoveListener(i)

	for reading := range updates {
		err = c.WriteJSON(reading)
		if err != nil {
			api.log.Error(err)
			return
		}
	}
}
