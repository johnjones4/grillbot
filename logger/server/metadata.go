package server

import (
	"main/core"
	"net/http"
)

func (api *API) getMetadata(w http.ResponseWriter, r *http.Request) {
	m, err := api.session.GetMetadata()
	if err != nil {
		sendError(api.log, w, http.StatusInternalServerError, err)
		return
	}
	sendJSON(api.log, w, http.StatusOK, m)
}

func (api *API) putMetadata(w http.ResponseWriter, r *http.Request) {
	var md core.Metadata
	err := readJSON(r, &md)
	if err != nil {
		sendError(api.log, w, http.StatusBadRequest, err)
		return
	}

	err = api.session.SetMetadata(md)
	if err != nil {
		sendError(api.log, w, http.StatusInternalServerError, err)
		return
	}

	sendJSON(api.log, w, http.StatusOK, md)
}
