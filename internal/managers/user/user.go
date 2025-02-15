package user

import (
	"encoding/json"
	"fmt"

	"qooked/internal/documentdb"
	"qooked/internal/instrumentation"
	"qooked/internal/models"
)

const collectionName = "users"

type UserManager struct {
	databaseClient  documentdb.DocumentDatabaseClient
	instrumentation instrumentation.Instrumentation
}

func NewUserManager(databaseClient documentdb.DocumentDatabaseClient, instrumentation instrumentation.Instrumentation) *UserManager {
	return &UserManager{
		databaseClient:  databaseClient,
		instrumentation: instrumentation,
	}
}

func (userManager *UserManager) GetUsers() (*[]models.User, error) {
	userManager.instrumentation.Log("Getting users from database...")
	documents, err := userManager.databaseClient.GetDocuments(collectionName)
	users := []models.User{}

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	for _, document := range *documents {
		user, err := convertDocToUser(&document)

		if err != nil {
			userManager.instrumentation.LogError(err.Error())
			return nil, err
		}

		users = append(users, *user)
	}

	userManager.instrumentation.Log(fmt.Sprintf("Number of users returned: %d.", len(users)))
	return &users, nil
}

func (userManager *UserManager) GetUser(userId string) (*models.User, error) {
	userManager.instrumentation.Log(fmt.Sprintf("Getting user with userID '%s' from database...", userId))
	document, err := userManager.databaseClient.GetDocument(collectionName, userId, userId)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	user, err := convertDocToUser(document)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return nil, err
	}

	userManager.instrumentation.Log(fmt.Sprintf("User with userID '%s' found.", userId))
	return user, nil
}

func (userManager *UserManager) UpsertUser(userId string, user *models.User) error {
	document, err := convertUserToDoc(user)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("Attempting to upsert user with userID '%s' to database...", userId))
	err = userManager.databaseClient.UpsertDocument(collectionName, userId, document, userId)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("User with userID '%s' successfully upserted to database.", userId))
	return nil
}

func (userManager *UserManager) DeleteUser(userId string) error {
	userManager.instrumentation.Log(fmt.Sprintf("Attempting to delete user with userID '%s' from database...", userId))
	err := userManager.databaseClient.DeleteDocument(collectionName, userId, userId)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("User with userID '%s' successfully deleted from database.", userId))
	return nil
}

func convertDocToUser(document *documentdb.Document) (*models.User, error) {
	var user models.User

	err := json.Unmarshal(document.Data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func convertUserToDoc(user *models.User) (*documentdb.Document, error) {
	var document documentdb.Document
	data, err := json.Marshal(*user)

	if err != nil {
		return nil, err
	}

	document.Data = data
	return &document, nil
}
