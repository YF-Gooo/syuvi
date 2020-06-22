package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"syuvi/event"
	"time"
)

var ServerInBuf chan event.Event

func CommandHandler(c *gin.Context) {
	target := c.PostForm("target")
	content := c.PostForm("content")
	ServerInBuf <- event.Event{target, "yf123",content }
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func QuitHandler(c *gin.Context) {
	ServerInBuf <- event.Event{"", "", "quit"}
	c.JSON(http.StatusOK, gin.H{
		"message": "quit success",
	})
}

func ServerBuild() (*http.Server){
	ServerInBuf = make(chan event.Event, 3)
	r := gin.Default()
	r.POST("/command", CommandHandler)
	r.GET("/quit", QuitHandler)
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}

