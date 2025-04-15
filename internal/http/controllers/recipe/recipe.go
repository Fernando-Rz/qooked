package recipe

import (
	"net/http"
	"qooked/internal/documentdb"
	"qooked/internal/managers/recipe"
	"qooked/internal/managers/user"
	"qooked/internal/models"

	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	recipeManager recipe.RecipeManager
	userManager   user.UserManager
}

func NewRecipeController(recipeManager recipe.RecipeManager, userManager user.UserManager) *RecipeController {
	return &RecipeController{
		recipeManager: recipeManager,
		userManager:   userManager,
	}
}

func (recipeController *RecipeController) GetRecipes(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username is required."})
		return
	}

	user, err := recipeController.userManager.GetUser(username)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		}

		return
	}

	recipes, err := recipeController.recipeManager.GetRecipes(user.UserId)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipes."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, *recipes)
}

func (recipeController *RecipeController) GetRecipe(ctx *gin.Context) {
	username := ctx.Param("username")
	recipeName := ctx.Param("recipe-name")

	if username == "" || recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username and RecipeName are required."})
		return
	}

	user, err := recipeController.userManager.GetUser(username)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		}

		return
	}

	recipe, err := recipeController.recipeManager.GetRecipe(recipeName, user.UserId)

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
	username := ctx.Param("username")
	recipeName := ctx.Param("recipe-name")

	if username == "" || recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username and RecipeName are required."})
		return
	}

	var recipeData models.Recipe

	if err := ctx.ShouldBindJSON(&recipeData); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	if recipeData.Name != recipeName {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "RecipeName in the URL does not match the RecipeName in the body."})
		return
	}

	user, err := recipeController.userManager.GetUser(username)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		}

		return
	}

	if err := recipeController.recipeManager.UpsertRecipe(recipeName, &recipeData, user.UserId); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update recipe."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully created or updated."})
}

func (recipeController *RecipeController) DeleteRecipe(ctx *gin.Context) {
	username := ctx.Param("username")
	recipeName := ctx.Param("recipe-name")

	if username == "" || recipeName == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username and RecipeName are required."})
		return
	}

	user, err := recipeController.userManager.GetUser(username)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		}

		return
	}

	if err := recipeController.recipeManager.DeleteRecipe(recipeName, user.UserId); err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Recipe not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe."})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Recipe successfully deleted."})
}
