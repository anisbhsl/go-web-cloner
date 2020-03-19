package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	asyncq "go-web-cloner/asynq"
	"go-web-cloner/handler"
	"log"
	"os"
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

	router := gin.Default()

	router.GET("/",handler.Index)
	router.GET("/api",handler.Index)
	router.GET("/api/scrape", handler.Scrape)
	router.GET("/api/status",handler.Status)
	router.GET("/api/stop", handler.Stop)
	router.GET("/report",handler.Report)
	//router.Handle("/static",http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	log.Fatal(router.Run("0.0.0.0:"+PORT))
}
