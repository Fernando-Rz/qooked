package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// recipe definition
type Recipe struct {
    Name         string       `json:"name"`
    Description  string       `json:"description"`
    Time         RecipeTime   `json:"time"`
    Servings     int          `json:"servings"`
    Ingredients  []Ingredient `json:"ingredients"`
    Instructions []string     `json:"instructions"`
}

// ingredient definition
type Ingredient struct {
    Name   string `json:"name"`
    Amount string `json:"amount"`
}

// time it takes to complete recipe
type RecipeTime struct {
	Prep  string `json:"prep"`
	Cook  string `json:"cook"`
	Total string `json:"total"`
}

var recipes = []Recipe{
	{
		Name: "",
		Description: "",
		Time: RecipeTime{
			Prep: "2",
			Cook: "2",
			Total: "2",
		},        
		Servings: 1,
		Ingredients: []Ingredient{
			{
				Name: "bread",
				Amount: "4 slices",
			},
		},
		Instructions: []string{
			"2",
		},
	},
}

// getRecipes responds with the list of all recipes as JSON.
func getRecipes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, recipes)
}

func main() {
    router := gin.Default()

	router.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Server is up and running.",
		})
	})

    router.GET("/recipes", getRecipes)

    // Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local development
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
