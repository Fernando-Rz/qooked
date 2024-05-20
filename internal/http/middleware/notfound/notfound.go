package notfound

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusNotFound,
		gin.H{
			"message": "Path not found.",
		})
}