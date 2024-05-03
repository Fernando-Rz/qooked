package main

import (
	"qooked/controllers"
	"qooked/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize server
    router := gin.Default()

	// health check routes
	//router.GET("/health-check", controllers.HealthCheck)

	// recipe scope routes
    router.GET("/recipes", controllers.GetRecipes)
	router.GET("/recipes/:recipe-name", controllers.GetRecipe)
	router.PUT("/recipes/:recipe-name", controllers.PutRecipe)
	router.DELETE("/recipes/:recipe-name", controllers.DeleteRecipe)

	// route not found
	router.Use(middleware.NotFound)

	// run server
	router.Run()
}
