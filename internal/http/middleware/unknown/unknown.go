package unknown

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnknownPath(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusNotFound,
		gin.H{
			"message": "Path not found.",
		})
}