package server

import (
	"qooked/internal/config"
	"qooked/internal/documentdb"
	"qooked/internal/documentdb/azure/cosmos"
	authController "qooked/internal/http/controllers/auth"
	"qooked/internal/http/controllers/health"
	recipeController "qooked/internal/http/controllers/recipe"
	userController "qooked/internal/http/controllers/user"
	"qooked/internal/http/middleware/auth"
	"qooked/internal/http/middleware/unknown"
	"qooked/internal/instrumentation"
	mockInstrumentation "qooked/internal/instrumentation/mock"
	recipeManager "qooked/internal/managers/recipe"
	userManager "qooked/internal/managers/user"

	"github.com/gin-gonic/gin"
)

// Server definition
type Server struct {
	config                 config.Config
	instrumentation        instrumentation.Instrumentation
	documentDatabaseClient documentdb.DocumentDatabaseClient
	router                 *gin.Engine
}

func NewServer(environmentName string) (*Server, error) {
	server := Server{}

	configFilePath := "cmd/api/configs/" + environmentName + ".json"
	err := server.initializeConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	err = server.initializeInstrumentation()
	if err != nil {
		return nil, err
	}

	err = server.initializeDocumentDatabaseClient()
	if err != nil {
		return nil, err
	}

	server.initializeRouter()
	return &server, nil
}

func (server *Server) initializeConfig(fileName string) error {
	config, err := config.NewConfig(fileName)
	if err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	server.config = config
	return nil
}

func (server *Server) initializeInstrumentation() error {
	server.instrumentation = mockInstrumentation.NewMockInstrumentation()

	err := server.instrumentation.InitializeInstrumentation()
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) initializeDocumentDatabaseClient() error {
	server.documentDatabaseClient = cosmos.NewCosmosDocumentDatabaseClient()

	err := server.documentDatabaseClient.InitializeClient(
		server.config.DocumentDatabaseUrl,
		server.config.DatabaseName)

	if err != nil {
		return err
	}

	err = server.documentDatabaseClient.TestConnection()

	if err != nil {
		return err
	}

	return nil
}

func (server *Server) initializeRouter() {
	server.router = gin.Default()

	// health check routes
	server.router.GET("/health", health.HealthCheck)

	// managers and controllers
	userManager := *userManager.NewUserManager(server.documentDatabaseClient, server.instrumentation)
	userController := *userController.NewUserController(userManager)
	recipeManager := *recipeManager.NewRecipeManager(server.documentDatabaseClient, server.instrumentation)

	recipeController := *recipeController.NewRecipeController(recipeManager, userManager)
	authController := *authController.NewAuthController(userManager)

	// public
	server.router.GET("/users", userController.GetUsers)
	server.router.POST("/register", authController.Register)
	server.router.POST("/login", authController.Login)

	// creating user group
	userGroup := server.router.Group("/users/:username").Use(auth.JWTAuthMiddleware())

	// the userGroup above defines a base route that all endpoints in the group will inherit
	userGroup.GET("", userController.GetUser)
	userGroup.DELETE("", userController.DeleteUser)
	userGroup.PUT("", userController.PutUser)
	userGroup.GET("/recipes", recipeController.GetRecipes)
	userGroup.GET("/recipes/:recipe-name", recipeController.GetRecipe)
	userGroup.PUT("/recipes/:recipe-name", recipeController.PutRecipe)
	userGroup.DELETE("/recipes/:recipe-name", recipeController.DeleteRecipe)

	server.router.Use(unknown.UnknownPath)
}

func (server *Server) Run() error {
	return server.router.Run()
}
