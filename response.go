package router

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, data interface{}) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	w.Write(jsonBytes)

	return nil
}
