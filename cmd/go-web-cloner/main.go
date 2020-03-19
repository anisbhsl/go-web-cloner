package main

import (
	"net/http"
	"log"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	asyncq "go-web-cloner/asynq"
	"go-web-cloner/handler"
)

func init(){
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found!")
	}
}

func main(){
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	dispatcher:=asyncq.NewDispatcher(1)
	dispatcher.Run()

	router := mux.NewRouter()

	router.HandleFunc("/",handler.Index())
	router.HandleFunc("/api",handler.Index())
	router.HandleFunc("/api/scrape", handler.Scrape())
	router.HandleFunc("/api/status",handler.Status())
	router.HandleFunc("/api/stop", handler.Stop())
	router.HandleFunc("/report",handler.Report())

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + PORT,
	}

	log.Println("[[main]] Server listening on: localhost:", PORT)
	log.Fatal(srv.ListenAndServe())
}
