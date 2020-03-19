package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func Stop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[[handler/stop]] New HTTP Request")

		response:=make(map[string]interface{})
		response["msg"]="Scrapping Stopped"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}