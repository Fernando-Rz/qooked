package recipe

import (
	"encoding/json"

	"qooked/internal/models"
	"qooked/internal/providers/doc"
)

const collectionName = "recipe"

type RecipeManager struct {
    databaseClient doc.DocumentDatabaseClient
}

func (recipeManager *RecipeManager) GetRecipe(recipeId string) (models.Recipe, error) {
	document, err := recipeManager.databaseClient.GetDocument(collectionName, recipeId)

	if err != nil {
		return models.Recipe{}, err
	}
    
	recipe, err := convertDocToRecipe(document)

	if err != nil {
		return models.Recipe{}, err
	}

	return recipe, nil
}

func (recipeManager *RecipeManager) GetRecipes() ([]models.Recipe, error) {
	documents, err := recipeManager.databaseClient.GetDocuments(collectionName)
	var recipes []models.Recipe

	if err != nil {
		return []models.Recipe{}, err
	}

	for _, document := range documents {
		recipe, err := convertDocToRecipe(document)

		if err != nil {
			return []models.Recipe{}, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (recipeManager *RecipeManager) UpsertRecipe(recipeId string, recipe models.Recipe) error {
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

func convertDocToRecipe(document doc.Document) (models.Recipe, error) {
	var recipe models.Recipe

    err := json.Unmarshal(document.Data, &recipe)
    if err != nil {
        return models.Recipe{}, err
    }

	return recipe, nil
}

func convertRecipeToDoc(recipe models.Recipe) (doc.Document, error){
	var document doc.Document
	data, err := json.Marshal(recipe)
   
	if err != nil {
        return doc.Document{}, err
    }

	document.Data = data
	return document, nil
}