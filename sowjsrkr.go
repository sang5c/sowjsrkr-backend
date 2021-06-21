package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const (
	authorize = "https://discord.com/api/oauth2/authorize"    // Base authorization URL
	token     = "https://discord.com/api/oauth2/token"        // Token URL
	revoke    = "https://discord.com/api/oauth2/token/revoke" // Token Revocation URL
)

func main() {
	router := gin.Default()
	router.GET("/", helloWorld)
	router.GET("/test", welcome)
	router.GET("/auth/callback", callback)
	router.Run()
}

func callback(context *gin.Context) {
	fmt.Println(context.Request.URL.Query())
}

func welcome(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"HI": "WELCOME",
	})
}

func helloWorld(context *gin.Context) {
	conf := os.Getenv("CONF")
	context.JSON(http.StatusOK, gin.H{
		"HELLO": "WORLD",
		"CONF":  conf,
	})
}
