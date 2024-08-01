package recipe

import (
	"net/http"
	"qooked/internal/managers/recipe"

	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	recipeManager recipe.RecipeManager
}

func NewRecipeController(recipeManager recipe.RecipeManager) *RecipeController {
	return &RecipeController{
		recipeManager: recipeManager,
	}
}

func (recipeController *RecipeController) GetRecipes(ctx *gin.Context) {
	recipes, err := recipeController.recipeManager.GetRecipes()

	if err != nil {
		ctx.IndentedJSON(
			http.StatusInternalServerError,
			err)
	}

	ctx.IndentedJSON(
		http.StatusOK,
		*recipes) // might need to update
}

func (recipeController *RecipeController) GetRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "GetRecipe called.",
		})
}

func (recipeController *RecipeController) PutRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "PutRecipe called.",
		})
}

func (recipeController *RecipeController) DeleteRecipe(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "DeleteRecipe called.",
		})
}
