package main

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
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
	//gin.SetMode(gin.ReleaseMode)  //use this in release


	dispatcher:=asyncq.NewDispatcher()

	//html:=template.Must(template.ParseFiles("public/templates/index.html"))

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.Static("/data","./data")

	router.Use(cors.Default(),static.Serve("/data",static.LocalFile("./data",true)))
	router.GET("/",handler.Index(dispatcher))
	router.GET("/api",handler.Index(dispatcher))
	router.POST("/api/scrape", handler.Scrape(dispatcher))
	router.GET("/api/status",handler.Status(dispatcher))
	router.GET("/api/stop", handler.Stop(dispatcher))
	router.GET("/api/redirect",handler.Redirect)
	router.GET("/report",handler.Report)
	//router.Handle("/static",http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	log.Fatal(router.Run("0.0.0.0:"+PORT))
}


