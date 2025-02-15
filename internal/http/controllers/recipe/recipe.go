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
	userId := ctx.Param("user-id")
	recipeId := ctx.Param("recipe-id")

	if userId == "" || recipeId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID and Recipe ID are required."})
		return
	}

	recipe, err := recipeController.recipeManager.GetRecipe(recipeId, userId)
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
	userId := ctx.Param("user-id")
	recipeId := ctx.Param("recipe-id")

	if userId == "" || recipeId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID and Recipe ID are required."})
		return
	}

	var recipeData models.Recipe
	if err := ctx.ShouldBindJSON(&recipeData); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	recipeData.UserId = userId
	recipeData.RecipeId = recipeId

	if err := recipeController.recipeManager.UpsertRecipe(recipeId, &recipeData, userId); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update recipe."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully created or updated."})
}

func (recipeController *RecipeController) PostRecipe(ctx *gin.Context) {
	userId := ctx.Param("user-id")

	if userId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required."})
		return
	}

	var recipeData models.Recipe
	if err := ctx.ShouldBindJSON(&recipeData); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	// Decide if the uuid should be set to the recipe ID here
	recipeData.UserId = userId

	if err := recipeController.recipeManager.UpsertRecipe(recipeData.RecipeId, &recipeData, userId); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully created."})
}

func (recipeController *RecipeController) DeleteRecipe(ctx *gin.Context) {
	userId := ctx.Param("user-id")
	recipeId := ctx.Param("recipe-id")

	if userId == "" || recipeId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID and Recipe ID are required."})
		return
	}

	if err := recipeController.recipeManager.DeleteRecipe(recipeId, userId); err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Recipe not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe."})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully deleted."})
}
