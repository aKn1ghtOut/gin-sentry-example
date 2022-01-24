package main

import (
	"log"
	"os"

	"ginserver/controllers"

	"github.com/gin-gonic/gin"

	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	sentryDsn := os.Getenv("SENTRY_DSN")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		TracesSampleRate: 1,
		Debug:            true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	r := GetRouter()
	r.Use(gin.Logger())
	if err := r.Run(":8090"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func GetRouter() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")
	corsConfig.AllowMethods = append(corsConfig.AllowMethods, "OPTIONS")

	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(cors.New(corsConfig))

	r.GET("/getPosts", controllers.GetPostsHandler)

	return r
}
