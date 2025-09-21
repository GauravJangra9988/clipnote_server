package main

import (
	"clipnote/server/cmd/auth"
	"clipnote/server/cmd/db"
	handle "clipnote/server/cmd/handler.go"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",
			"http://localhost:8081",
			"https://clipnote-frontend.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db.DBstart()

	r.GET("/", handle.Test)
	r.POST("/register", handle.Register)
	r.POST("/login", handle.Login)

	protected := r.Group("/")
	protected.Use(auth.AuthenticateUser())
	protected.POST("/save", handle.Save)
	protected.GET("/clips", handle.GetAllClips)
	protected.DELETE("/delete", handle.DeleteClip)
	protected.POST("/logout", handle.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":"+port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	db.DB.Close()

}
