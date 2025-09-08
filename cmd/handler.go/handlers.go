package handle

import (
	db "clipnote/server/cmd/db"
	"clipnote/server/cmd/token"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server running"})
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

	_, err := db.DB.Exec(context.Background(), "INSERT INTO users (user_name,email,password) VALUES ($1,$2,$3)", User_name, Email, Password)
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

	fmt.Println(req)

	var savedPassword string
	var user_name string

	err := db.DB.QueryRow(context.Background(), "SELECT password,user_name FROM users WHERE email = $1", req.Email).Scan(&savedPassword, &user_name)
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

	token, err := token.CreateToken(user_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfuly", "token": token})

}

// Save handles the /save POST request
func Save(c *gin.Context) {

	username := c.GetString("username")

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

	_, dberr := db.DB.Exec(context.Background(), "INSERT INTO clip_data (user_name,title,text,tag) VALUES ($1,$2,$3,$4)", username, JSONformatedData.Title, JSONformatedData.Text, JSONformatedData.Tag)
	if dberr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dberr.Error()})
		return
	}

	c.JSON(http.StatusOK, JSONformatedData)
}

type ClipnoteDataStruct struct {
	Id int64 `json:"id"`
	JSONformatClipboardData
	Time time.Time `json:"added_at"`
}

func GetAllClips(c *gin.Context) {

	username := c.GetString("username")

	rows, err := db.DB.Query(context.Background(), "SELECT id, title, text, tag, added_at FROM clip_data WHERE user_name = $1", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var ClipnoteData []ClipnoteDataStruct

	for rows.Next() {

		var d ClipnoteDataStruct
		err := rows.Scan(&d.Id, &d.Title, &d.Text, &d.Tag, &d.Time)
		if err != nil {
			fmt.Println("err:", err)
		}
		ClipnoteData = append(ClipnoteData, d)

	}

	c.JSON(http.StatusOK, ClipnoteData)

}

type ClipId struct {
	Id int32 `json:"id"`
}

func DeleteClip(c *gin.Context) {

	var req ClipId

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Clipnote id not provided"})
		return
	}

	cmdTag, err := db.DB.Exec(context.Background(), "DELETE FROM clip_data where id = $1", req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to delete clip data"})
		return
	}
	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"err": "Clipdata not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clipnote deleted"})

}

func Logout(c *gin.Context) {

	c.SetCookie("token", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Logged Out successfuly"})
}




