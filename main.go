package main

import (
	"fmt"
	"net/http"

	"github.com/CarlosBarbosaFilho/go-jwt-authentication/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DatabaseConnection()
	initializers.ConfigureTables()
}

func main() {
	fmt.Println("Hello World")

	router := gin.Default()
	router.GET("/welcome", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to JWT in golang project"})
	})

	router.Run()
}
