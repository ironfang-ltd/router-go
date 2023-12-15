package router

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, data interface{}) {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, _ = w.Write(jsonBytes)
}
