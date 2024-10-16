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
	service.InitializePath()

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"content": "Hello World",
		})
	})

	//@todo: implement for playlists too

	group := router.Group("/api/v1")
	{
		group.POST("/music", MusicHandler)
		group.GET("/version", VersionHandler)
	}

	go service.ScheduledCleanUp()

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

	path, fileName, success, output, err := service.DownloadMusic(dto.Url)

	if !success {
		context.IndentedJSON(500, gin.H{
			"error":  fmt.Sprintf("Error while retrieving data %v", err),
			"output": output,
		})

		return
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

func VersionHandler(context *gin.Context) {
	version, err := service.GetAppVersion()

	if err != nil {
		version = "0.0 (def)"
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"version": version,
	})
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
