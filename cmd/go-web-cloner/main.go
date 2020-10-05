package main

import (
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

	//config:=cors.DefaultConfig()
	//config.AllowedOrigins=[]string{"http://estimate.appgreenhouse.io"}
	//router.Use(cors.New(config))

	router.Use(CORSMiddleware(),static.Serve("/data",static.LocalFile("./data",true)))
	//router.Use(cors.New(config),static.Serve("/data",static.LocalFile("./data",true)))


	router.GET("/",handler.Index(dispatcher))
	router.GET("/api",handler.Index(dispatcher))
	router.POST("/api/scrape", handler.Scrape(dispatcher))
	router.GET("/api/status",handler.Status(dispatcher))
	router.GET("/api/stop", handler.Stop(dispatcher))
	router.GET("/api/redirect",handler.Redirect)
	router.GET("/report",handler.Report(dispatcher))
	//router.Handle("/static",http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	log.Fatal(router.Run("0.0.0.0:"+PORT))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://estimate.appgreenhouse.io/")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

