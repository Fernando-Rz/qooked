package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRecipes(ctx *gin.Context) {
	// var recipes = []Recipe{
	// 	{
	// 		Name: "",
	// 		Description: "",
	// 		Time: RecipeTime{
	// 			Prep: "2",
	// 			Cook: "2",
	// 			Total: "2",
	// 		},        
	// 		Servings: 1,
	// 		Ingredients: []Ingredient{
	// 			{
	// 				Name: "bread",
	// 				Amount: "4 slices",
	// 			},
	// 		},
	// 		Instructions: []string{
	// 			"2",
	// 		},
	// 	},
	// }

	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "GetRecipes called.",
		})
}

func GetRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "GetRecipe called.",
		})
}

func PutRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "PutRecipe called.",
		})
}

func DeleteRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "DeleteRecipe called.",
		})
}
