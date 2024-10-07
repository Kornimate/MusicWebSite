package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"content": "Hello World",
		})
	})

	router.Run("localhost:8080")
}
