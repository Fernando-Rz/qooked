package recipe

import (
	"encoding/json"
	"errors"
	"fmt"

	"qooked/internal/documentdb"
	"qooked/internal/instrumentation"
	"qooked/internal/models"

	"github.com/google/uuid"
)

const collectionName = "recipes"

var (
	ErrConflictingRecipes = errors.New("multiple recipes with the same name")
	ErrRecipeNameExists   = errors.New("recipe name already exists")
)

type RecipeManager struct {
	databaseClient  documentdb.DocumentDatabaseClient
	instrumentation instrumentation.Instrumentation
}

func NewRecipeManager(databaseClient documentdb.DocumentDatabaseClient, instrumentation instrumentation.Instrumentation) *RecipeManager {
	return &RecipeManager{
		databaseClient:  databaseClient,
		instrumentation: instrumentation,
	}
}

func (recipeManager *RecipeManager) GetRecipes(userId string) (*[]models.Recipe, error) {
	query := fmt.Sprintf("SELECT * FROM %s c", collectionName)

	recipeManager.instrumentation.Log(fmt.Sprintf("Getting recipes for user with userId '%s'...", userId))
	documents, err := recipeManager.databaseClient.GetDocuments(collectionName, query, userId)
	recipes := []models.Recipe{}

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	for _, document := range *documents {
		recipe, err := convertDocToRecipe(&document)

		if err != nil {
			recipeManager.instrumentation.LogError(err.Error())
			return nil, err
		}

		recipes = append(recipes, *recipe)
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Number of recipes returned: %d.", len(recipes)))
	return &recipes, nil
}

func (recipeManager *RecipeManager) GetRecipe(recipeName string, userId string) (*models.Recipe, error) {
	// TODO: check if we can pass a blank recipeName for GetRecipes()
	query := fmt.Sprintf("SELECT * FROM %s c WHERE c.recipeName = '%s'", collectionName, recipeName)

	recipeManager.instrumentation.Log(fmt.Sprintf("Getting recipe with recipeName '%s' for user with userId '%s'...", recipeName, userId))
	documents, err := recipeManager.databaseClient.GetDocuments(collectionName, query, userId)
	recipes := []models.Recipe{}

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	for _, document := range *documents {
		recipe, err := convertDocToRecipe(&document)

		if err != nil {
			recipeManager.instrumentation.LogError(err.Error())
			return nil, err
		}

		recipes = append(recipes, *recipe)
	}

	recipeCount := len(recipes)
	recipeManager.instrumentation.Log(fmt.Sprintf("Number of recipes returned: %d.", recipeCount))

	if recipeCount > 1 {
		return nil, ErrConflictingRecipes
	} else if recipeCount == 0 {
		return nil, documentdb.ErrDocumentNotFound
	} else {
		return &recipes[0], nil
	}
}

func (recipeManager *RecipeManager) UpsertRecipe(recipeName string, recipe *models.Recipe, userId string) error {
	creatingNewRecipe := false
	updatingExistingRecipeName := false

	if recipe.RecipeId == "" {
		creatingNewRecipe = true
		recipe.RecipeId = uuid.New().String()
	} else {
		currentRecipeDoc, err := recipeManager.databaseClient.GetDocument(collectionName, recipe.RecipeId, userId)

		if err != nil {
			if err == documentdb.ErrDocumentNotFound {
				creatingNewRecipe = true
				recipe.RecipeId = uuid.New().String()
			} else {
				recipeManager.instrumentation.LogError(err.Error())
				return err
			}
		}

		if !creatingNewRecipe {
			currentRecipe, err := convertDocToRecipe(currentRecipeDoc)

			if err != nil {
				recipeManager.instrumentation.LogError(err.Error())
				return err
			}

			updatingExistingRecipeName = currentRecipe.RecipeName != recipeName
		}
	}

	if creatingNewRecipe || updatingExistingRecipeName {
		existingRecipe, err := recipeManager.GetRecipe(recipeName, userId)

		if err != nil && err != documentdb.ErrDocumentNotFound {
			recipeManager.instrumentation.LogError(err.Error())
			return err
		}

		if err == nil && existingRecipe.RecipeId != recipe.RecipeId {
			// the recipeId passed by the caller does not match what is stored in the db
			recipeManager.instrumentation.LogError(ErrRecipeNameExists.Error())
			return ErrRecipeNameExists
		}
	}

	recipe.UserId = userId
	document, err := convertRecipeToDoc(recipe)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Attempting to upsert recipe with recipeName '%s' for user with userId '%s'...", recipeName, userId))
	err = recipeManager.databaseClient.UpsertDocument(collectionName, recipe.RecipeId, document, userId)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Recipe with recipeName '%s' successfully upserted to database for user with userId '%s'.", recipeName, userId))
	return nil
}

func (recipeManager *RecipeManager) DeleteRecipe(recipeName string, userId string) error {
	currentRecipe, err := recipeManager.GetRecipe(recipeName, userId)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Attempting to delete recipe with recipeName '%s' for user with userId '%s'...", recipeName, userId))
	err = recipeManager.databaseClient.DeleteDocument(collectionName, currentRecipe.RecipeId, userId)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Recipe with recipeName '%s' successfully deleted from database for user with userId '%s'.", recipeName, userId))
	return nil
}

func convertDocToRecipe(document *documentdb.Document) (*models.Recipe, error) {
	var recipe models.Recipe

	err := json.Unmarshal(document.Data, &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func convertRecipeToDoc(recipe *models.Recipe) (*documentdb.Document, error) {
	var document documentdb.Document
	data, err := json.Marshal(*recipe)

	if err != nil {
		return nil, err
	}

	document.Data = data
	return &document, nil
}
