package main

import (
	"qooked/internal/http/controllers/health"
	"qooked/internal/http/controllers/recipe"
	"qooked/internal/http/middleware/notfound"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
    router := gin.Default()

	// health check routes
	router.GET("/health-check", health.HealthCheck)

	// recipe scope routes
    router.GET("/recipes", recipe.GetRecipes)
	router.GET("/recipes/:recipe-name", recipe.GetRecipe)
	router.PUT("/recipes/:recipe-name", recipe.PutRecipe)
	router.DELETE("/recipes/:recipe-name", recipe.DeleteRecipe)

	// route not found
	router.Use(notfound.NotFound)

	// run server
	router.Run()
}
