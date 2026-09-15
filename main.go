package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	database "github.com/glenngenre/code-paste-service/database"
	docs "github.com/glenngenre/code-paste-service/docs"
	v1 "github.com/glenngenre/code-paste-service/handlers/v1"
	"github.com/glenngenre/code-paste-service/middleware"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed static/*
var staticFiles embed.FS

// @title Code Paste Service API
// @version 1.0
// @description API for code snippet sharing service
// @host api.apps.skwtr.com
// @BasePath /codelibrary/api
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	middleware.Init(database.DB)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(setupCORS())
	r.Use(middleware.RateLimitMiddleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = ""
	}

	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to load embedded static files:", err)
	}
	r.StaticFS(basePath+"/static", http.FS(assets))

	docs.SwaggerInfo.BasePath = basePath + "/api"
	v1.RegisterRoutes(r, basePath)

	r.GET(basePath+"/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":" + port)
}

func setupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"https://codelibrary.skwtr.com",
			"http://localhost:5173",
			"http://localhost:8080",
			"http://localhost",
			"http://localhost:9080",
		},

		AllowMethods: []string{
			"GET",
			"POST",
		},

		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},

		ExposeHeaders: []string{
			"Content-Length",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})
}
