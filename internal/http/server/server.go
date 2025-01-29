package server

import (
	"qooked/internal/config"
	"qooked/internal/documentdb"
	"qooked/internal/documentdb/azure/cosmos"
	"qooked/internal/http/controllers/health"
	recipeController "qooked/internal/http/controllers/recipe"
	"qooked/internal/http/middleware/unknown"
	"qooked/internal/instrumentation"
	mockInstrumentation "qooked/internal/instrumentation/mock"
	recipeManager "qooked/internal/managers/recipe"

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

	// recipe scope routes
	recipeManager := *recipeManager.NewRecipeManager(server.documentDatabaseClient, server.instrumentation)
	recipeController := *recipeController.NewRecipeController(recipeManager)

	server.router.GET("/recipes", recipeController.GetRecipes)
	server.router.GET("/recipes/:recipe-name", recipeController.GetRecipe)
	server.router.PUT("/recipes/:recipe-name", recipeController.PutRecipe)
	server.router.DELETE("/recipes/:recipe-name", recipeController.DeleteRecipe)

	// middlewares
	server.router.Use(unknown.UnknownPath)
}

func (server *Server) Run() error {
	return server.router.Run()
}
