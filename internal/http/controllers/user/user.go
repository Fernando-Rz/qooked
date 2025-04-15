package user

import (
	"net/http"
	"qooked/internal/documentdb"
	"qooked/internal/managers/user"
	"qooked/internal/models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userManager user.UserManager
}

func NewUserController(userManager user.UserManager) *UserController {
	return &UserController{
		userManager: userManager,
	}
}

func (userController *UserController) GetUsers(ctx *gin.Context) {
	users, err := userController.userManager.GetUsers()

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users."})
		return
	}

	ctx.IndentedJSON(http.StatusOK, *users)
}

func (userController *UserController) GetUser(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username is required."})
		return
	}

	user, err := userController.userManager.GetUser(username)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, *user)
}

func (userController *UserController) PutUser(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username is required."})
		return
	}

	var userData models.User

	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	if userData.Username != username {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username in the URL does not match the Username in the body."})
		return
	}

	if err := userController.userManager.UpsertUser(username, &userData); err != nil {
		if err == user.ErrUsernameExists {
			ctx.IndentedJSON(http.StatusConflict, gin.H{"error": "Username already exists."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user."})
		}

		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "User successfully created or updated."})
}

func (userController *UserController) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username is required."})
		return
	}

	if err := userController.userManager.DeleteUser(username); err != nil {
		if err == documentdb.ErrDocumentNotFound {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found."})
		} else {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user."})
		}
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "User successfully deleted."})
}
