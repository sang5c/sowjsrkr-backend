package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	authorize = "https://discord.com/api/oauth2/authorize"    // Base authorization URL
	token     = "https://discord.com/api/oauth2/token"        // Token URL
	revoke    = "https://discord.com/api/oauth2/token/revoke" // Token Revocation URL
)

func main() {
	router := gin.Default()
	router.GET("/", helloWorld)
	router.Run()
}

func helloWorld(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"HELLO": "WORLD",
	})
}
