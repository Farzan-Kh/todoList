package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func logRequest(c *gin.Context) {
	log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next()
	log.Printf("Response: %d", c.Writer.Status())
}

func main() {
	r := gin.Default()
	r.Use(logRequest)

	r.GET("/audit/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Audit Service is running"})
	})

	r.Run(":8083") // Run on port 8083
}
