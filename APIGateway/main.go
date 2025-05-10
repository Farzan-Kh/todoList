package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Any("/users/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, "http://userservice:8081")
	})

	r.Any("/tasks/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, "http://taskservice:8082")
	})

	r.Any("/audit/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, "http://auditservice:8083")
	})

	r.Run(":8080") // Run on port 8080
}

func proxyRequest(c *gin.Context, target string) {
	url := target + c.Request.URL.Path[0:] // Use the full request path
	fmt.Println("Forwarding request to:", url)

	req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
