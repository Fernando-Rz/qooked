package recipe

import (
	"encoding/json"

	"qooked/internal/models"
	"qooked/internal/providers/documentdb"
)

const collectionName = "recipe"

type RecipeManager struct {
    databaseClient documentdb.DocumentDatabaseClient
}

func (recipeManager *RecipeManager) GetRecipes() (*[]models.Recipe, error) {
	documents, err := recipeManager.databaseClient.GetDocuments(collectionName)
	var recipes []models.Recipe

	if err != nil {
		return nil, err
	}

	for _, document := range *documents {
		recipe, err := convertDocToRecipe(&document)

		if err != nil {
			return nil, err
		}

		recipes = append(recipes, *recipe)
	}

	return &recipes, nil
}

func (recipeManager *RecipeManager) GetRecipe(recipeId string) (*models.Recipe, error) {
	document, err := recipeManager.databaseClient.GetDocument(collectionName, recipeId)

	if err != nil {
		return nil, err
	}
    
	recipe, err := convertDocToRecipe(document)

	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (recipeManager *RecipeManager) UpsertRecipe(recipeId string, recipe *models.Recipe) error {
	document, err := convertRecipeToDoc(recipe)

	if err != nil {
		return err
	}

	err = recipeManager.databaseClient.UpsertDocument(collectionName, recipeId, document)
	
	if err != nil {
		return err
	}

	return nil
}

func (recipeManager *RecipeManager) DeleteRecipe(recipeId string) error {
	err := recipeManager.databaseClient.DeleteDocument(collectionName, recipeId)

	if err != nil {
		return err
	}

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

func convertRecipeToDoc(recipe *models.Recipe) (*documentdb.Document, error){
	var document documentdb.Document
	data, err := json.Marshal(*recipe)
   
	if err != nil {
        return nil, err
    }

	document.Data = data
	return &document, nil
}