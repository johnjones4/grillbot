package server

import (
	"main/ui"
	"net/http"
)

func (api *API) getChart(w http.ResponseWriter, r *http.Request) {
	bytes, err := ui.GenerateChart(api.session)
	if err != nil {
		sendError(api.log, w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
