package auth

import (
	"log"
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
	// required fields according to the model
	var req struct {
		Username    string `json:"username" binding:"required,alphanum"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"` // we can remove these string requirements
		ProfileName string `json:"profileName" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// check if username already exists
	_, err := ac.userManager.GetUser(req.Username)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user object
	newUser := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		ProfileName: req.ProfileName,
	}

	err = ac.userManager.UpsertUser(req.Username, &newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// auto-login after registration, tested it and it returns the same token if you hit the login endpoint afterwards
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
		// remove log debugging
		log.Printf("Login failed: GetUser(%q) error: %v", req.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare hashed password stored vs password attempt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// remove log debugging
		log.Printf("Login failed: password mismatch for %q: %v", req.Username, err)
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
