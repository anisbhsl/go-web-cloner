package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[[handler/index]] New HTTP Request")

		response:=make(map[string]interface{})
		response["msg"]="Website Cloner"
		response["status"]="Under Development : WIP"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}
