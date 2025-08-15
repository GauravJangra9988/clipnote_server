package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ClipboardRequest struct {
	Data string `json:"data"`
}

type JSONformatClipboardData struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Tag   string `json:"tag"`
}

type RegisterRequest struct {
	User_name string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Save handles the /save POST request
func Save(c *gin.Context) {
	var req ClipboardRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	llmOutputData := ProcessWithLLM(req.Data)
	cleanllmData := CleanLLMOutput(llmOutputData)

	JSONformatedData, err := ParseLLMoutput(cleanllmData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse LLM output"})
		return
	}

	fmt.Println("Tittle: ", JSONformatedData.Title)
	fmt.Println("Text: ", JSONformatedData.Text)
	fmt.Println("Tag: ", JSONformatedData.Tag)

	c.JSON(http.StatusOK, JSONformatedData)
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	User_name := req.User_name
	Email := req.Email
	Password := req.Password

	_, err := DB.Exec(context.Background(), "INSERT INTO users (user_name,email,password) VALUES ($1,$2,$3)", User_name, Email, Password)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registeration successful"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var savedPassword string

	err := DB.QueryRow(context.Background(), "SELECT password FROM users WHERE email = $1", req.Email).Scan(&savedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Databse error"})
		return
	}

	if savedPassword != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfuly"})

}
