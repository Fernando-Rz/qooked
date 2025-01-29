package recipe

import (
	"net/http"
	"qooked/internal/documentdb"
	"qooked/internal/managers/recipe"
	"qooked/internal/models"

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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipes."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, *recipes)
}

func (recipeController *RecipeController) GetRecipe(ctx *gin.Context) {
	recipeName := ctx.Param("recipe-name")
	if recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Recipe name is required."})
		return
	}

	recipe, err := recipeController.recipeManager.GetRecipe(recipeName)
	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Recipe not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipe."})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, *recipe)
}

func (recipeController *RecipeController) PutRecipe(ctx *gin.Context) {
	recipeName := ctx.Param("recipe-name")
	if recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Recipe name is required."})
		return
	}

	var recipeData models.Recipe
	if err := ctx.ShouldBindJSON(&recipeData); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	recipeData.Name = recipeName
	recipeData.Id = recipeName
	recipeData.PartitionKey = "recipes"

	if err := recipeController.recipeManager.UpsertRecipe(recipeName, &recipeData); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update recipe."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully created or updated."})
}

func (recipeController *RecipeController) DeleteRecipe(ctx *gin.Context) {
	recipeName := ctx.Param("recipe-name")
	if recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Recipe name is required."})
		return
	}

	if err := recipeController.recipeManager.DeleteRecipe(recipeName); err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Recipe not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe."})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully deleted."})
}
