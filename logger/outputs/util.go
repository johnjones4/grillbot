package outputs

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, status int, info any) {
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func errorResponse(w http.ResponseWriter, status int, err error) {
	msg := map[string]any{
		"ok":      false,
		"message": err.Error(),
	}
	jsonResponse(w, status, msg)
}
