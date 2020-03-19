package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[[handler/status]] New HTTP Request")

		response:=make(map[string]interface{})
		response["job_id"]="12345678"
		response["msg"]="Scrapper Running"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}
