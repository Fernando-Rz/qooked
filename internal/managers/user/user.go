package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"qooked/internal/documentdb"
	"qooked/internal/instrumentation"
	"qooked/internal/models"

	"github.com/google/uuid"
)

const collectionName = "users"
const universalGroupId = "1"

var (
	ErrConflictingUsers = errors.New("multiple users with the same username")
	ErrUsernameExists   = errors.New("username already exists")
)

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
	query := fmt.Sprintf("SELECT * FROM %s c", collectionName)

	userManager.instrumentation.Log("Getting users from database...")
	documents, err := userManager.databaseClient.GetDocuments(collectionName, query, universalGroupId)
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

func (userManager *UserManager) GetUser(username string) (*models.User, error) {
	// TODO: check if we can pass a blank username for GetUsers()
	query := fmt.Sprintf("SELECT * FROM %s c WHERE c.username = '%s'", collectionName, username)

	userManager.instrumentation.Log(fmt.Sprintf("Getting user with username '%s' from database...", username))
	documents, err := userManager.databaseClient.GetDocuments(collectionName, query, universalGroupId)
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

	userCount := len(users)
	userManager.instrumentation.Log(fmt.Sprintf("Number of users returned: %d.", userCount))

	if userCount > 1 {
		return nil, ErrConflictingUsers
	} else if userCount == 0 {
		return nil, documentdb.ErrDocumentNotFound
	} else {
		return &users[0], nil
	}
}

func (userManager *UserManager) UpsertUser(username string, user *models.User) error {
	creatingNewUser := false
	updatingExistingUsername := false
	currentUserDoc, err := userManager.databaseClient.GetDocument(collectionName, user.UserId, universalGroupId)

	if err != nil {
		if err == documentdb.ErrDocumentNotFound {
			creatingNewUser = true
			user.UserId = uuid.New().String()
		} else {
			userManager.instrumentation.LogError(err.Error())
			return err
		}
	}

	if !creatingNewUser {
		currentUser, err := convertDocToUser(currentUserDoc)

		if err != nil {
			userManager.instrumentation.LogError(err.Error())
			return err
		}

		updatingExistingUsername = currentUser.Username != username
	}

	if creatingNewUser || updatingExistingUsername {
		existingUser, err := userManager.GetUser(username)

		if err != nil && err != documentdb.ErrDocumentNotFound {
			userManager.instrumentation.LogError(err.Error())
			return err
		}

		if err == nil && existingUser.UserId != user.UserId {
			// the userId passed by the caller does not match what is stored in the db
			userManager.instrumentation.LogError(ErrUsernameExists.Error())
			return ErrUsernameExists
		}
	}

	document, err := convertUserToDoc(user)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("Attempting to upsert user with username '%s' to database...", username))
	err = userManager.databaseClient.UpsertDocument(collectionName, user.UserId, document, universalGroupId)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("User with username '%s' successfully upserted to database.", username))
	return nil
}

func (userManager *UserManager) DeleteUser(username string) error {
	currentUser, err := userManager.GetUser(username)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("Attempting to delete user with username '%s' from database...", username))
	err = userManager.databaseClient.DeleteDocument(collectionName, currentUser.UserId, universalGroupId)

	if err != nil {
		userManager.instrumentation.LogError(err.Error())
		return err
	}

	userManager.instrumentation.Log(fmt.Sprintf("User with username '%s' successfully deleted from database.", username))
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
