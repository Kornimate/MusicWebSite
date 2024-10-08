package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	model "example.com/music-app/models"
	service "example.com/music-app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"content": "Hello World",
		})
	})

	group := router.Group("/api/v1")
	{
		group.POST("/music", MusicHandler)
	}

	//scheduled task for clean up with go keyword and time.Sleep(x*time.Minute)
	//implement clean up goroutine in services
	//handle and test the same for playlists

	host := os.Getenv("HOST")

	if host == "" {
		host = "localhost:8080"
	}

	router.Run(host)
}

func MusicHandler(context *gin.Context) {

	var dto model.MusicDTO

	err := context.BindJSON(&dto)

	if err != nil {
		context.String(http.StatusBadRequest, "Error in HTTP request")
		return
	}

	fmt.Println(dto)

	path, fileName, success, err := service.DownloadMusic(dto.Url, dto.ActionType)

	if !success {
		context.String(500, fmt.Sprintf("Error while retrieving data %v", err))
	}

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))
	context.Header("Content-Type", "application/octet-stream")
	context.Header("FileName", fileName)
	context.File(path + string(filepath.Separator) + fileName)
	context.Status(http.StatusOK)

	defer service.CleanUp(path)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Content-Disposition, FileName")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, FileName")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
