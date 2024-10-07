package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	model "example.com/music-app/models"
	service "example.com/music-app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"content": "Hello World",
		})
	})

	group := router.Group("/api/v1")
	{
		group.POST("/music", MusicHandler)
	}

	router.Run("localhost:8080")
}

func MusicHandler(context *gin.Context) {

	var dto model.MusicDTO

	err := context.BindJSON(&dto)

	if err != nil {
		context.String(http.StatusBadRequest, "Error in HTTP request")
		return
	}

	path, fileName, success, err := service.DownloadMusic(dto.Url, dto.ActionType)

	if !success {
		context.String(500, fmt.Sprintf("Error while retrieving data %v", err))
	}

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))
	context.Header("Content-Type", "application/octet-stream")
	context.File(fmt.Sprintf("%v%v%v", path, filepath.Separator, fileName))

	defer service.CleanUp(path)
}
