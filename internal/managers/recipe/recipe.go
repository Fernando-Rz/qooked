package recipe

import (
	"encoding/json"
	"fmt"

	"qooked/internal/documentdb"
	"qooked/internal/instrumentation"
	"qooked/internal/models"
)

const collectionName = "recipes"

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

func (recipeManager *RecipeManager) GetRecipes() (*[]models.Recipe, error) {
	recipeManager.instrumentation.Log("Getting recipes from database...")
	documents, err := recipeManager.databaseClient.GetDocuments(collectionName)
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

func (recipeManager *RecipeManager) GetRecipe(recipeId string) (*models.Recipe, error) {
	recipeManager.instrumentation.Log(fmt.Sprintf("Getting recipe with recipeID '%s' from database...", recipeId))
	document, err := recipeManager.databaseClient.GetDocument(collectionName, recipeId)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	recipe, err := convertDocToRecipe(document)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Recipe with recipeID '%s' found.", recipeId))
	return recipe, nil
}

func (recipeManager *RecipeManager) UpsertRecipe(recipeId string, recipe *models.Recipe) error {
	document, err := convertRecipeToDoc(recipe)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Attempting to upsert recipe with recipeID '%s' to database...", recipeId))
	err = recipeManager.databaseClient.UpsertDocument(collectionName, recipeId, document)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Recipe with recipeID '%s' successfully upserted to database.", recipeId))
	return nil
}

func (recipeManager *RecipeManager) DeleteRecipe(recipeId string) error {
	recipeManager.instrumentation.Log(fmt.Sprintf("Attempting to delete recipe with recipeID '%s' from database...", recipeId))
	err := recipeManager.databaseClient.DeleteDocument(collectionName, recipeId)

	if err != nil {
		recipeManager.instrumentation.LogError(err.Error())
		return err
	}

	recipeManager.instrumentation.Log(fmt.Sprintf("Recipe with recipeID '%s' successfully deleted from database.", recipeId))
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
