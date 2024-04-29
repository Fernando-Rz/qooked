package main

import (
	"net/http"

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
    router.GET("/recipes", getRecipes)

    router.Run("localhost:8080")
}
