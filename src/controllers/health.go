package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	ctx.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": "Server is up and running!",
		})
}