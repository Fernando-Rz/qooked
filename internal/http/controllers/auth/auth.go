package auth

import (
	"net/http"
	"qooked/internal/auth/jwt"
	"qooked/internal/managers/user"
	"qooked/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userManager user.UserManager
}

func NewAuthController(userManager user.UserManager) *AuthController {
	return &AuthController{
		userManager: userManager,
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {

	var userData models.User
	// var err error

	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user object
	newUser := models.User{
		Username:    userData.Username,
		Email:       userData.Email,
		Password:    string(hashedPassword),
		ProfileName: userData.ProfileName,
	}

	err = ac.userManager.UpsertUser(userData.Username, &newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// auto-login after registration
	token, err := jwt.GenerateJWT(newUser.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := ac.userManager.GetUser(req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare hashed password stored vs password attempt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := jwt.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
