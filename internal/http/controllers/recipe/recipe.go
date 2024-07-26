package recipe

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRecipes(ctx *gin.Context) {
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
