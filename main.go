package main

import (
	"clipnote/server/cmd"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cmd.DBstart()

	r.POST("/register", cmd.Register)
	r.POST("/login", cmd.Login)
	r.POST("/save", cmd.Save)

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cmd.DB.Close()

}
