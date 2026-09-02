package generator

import (
	"fmt"
)

func getMainFileTemplate(modulePath string) string {
	return fmt.Sprintf(`
package main

import (
	"github.com/gin-gonic/gin"
	"%s/config"
	"%s/internal/routes"
)

func init() {
	config.LoadEnv()
}

func main() {
	app := gin.Default()
	routes.ExampleRoute(app)
	app.Run(":8080")
}`, modulePath, modulePath)
}

func getConfigFileTemplate() string {
	return `
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file, error: %v", err)
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}
`
}

func getHandlerExampleTemplate() string {
	return `
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Server is running smoothly 🚀",
	})
}
`
}

func getRouteExampleTemplate(modulePath string) string {
	return fmt.Sprintf(`
package routes

import (
	"github.com/gin-gonic/gin"
	"%s/internal/handlers"
)

func ExampleRoute(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.Healthcheck)
	}
}
	`, modulePath)
}

func getMakefileTemplate() string {
	return `
run:
	go run ./cmd/api/
	`
}
