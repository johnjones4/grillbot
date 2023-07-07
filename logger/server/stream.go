package server

import (
	"main/core"
	"net/http"

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
