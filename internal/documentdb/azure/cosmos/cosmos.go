package cosmos

import (
	"context"
	"errors"
	"fmt"
	"qooked/internal/documentdb"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type CosmosDocumentDatabaseClient struct {
	client *azcosmos.DatabaseClient
}

func NewCosmosDocumentDatabaseClient() *CosmosDocumentDatabaseClient {
	return &CosmosDocumentDatabaseClient{}
}

func (db *CosmosDocumentDatabaseClient) InitializeClient(endpointUrl string, databaseName string) error {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	clientOptions := azcosmos.ClientOptions{
		EnableContentResponseOnWrite: false,
	}

	client, err := azcosmos.NewClient(endpointUrl, credential, &clientOptions)
	if err != nil {
		return err
	}

	databaseClient, err := client.NewDatabase(databaseName)
	if err != nil {
		return err
	}

	db.client = databaseClient

	return nil
}

func (db *CosmosDocumentDatabaseClient) TestConnection() error {
	_, err := db.client.Read(context.TODO(), nil)

	if err != nil {
		return err
	}

	return nil
}

func (db *CosmosDocumentDatabaseClient) GetDocuments(collection string) (*[]documentdb.Document, error) {
	documents := []documentdb.Document{}

	container, err := db.client.NewContainer(collection)
	if err != nil {
		return nil, err
	}

	// TODO: Update this to userID to partition data by user
	partitionKey := azcosmos.NewPartitionKeyString(collection)
	query := fmt.Sprintf("SELECT * FROM %s c", collection)

	pager := container.NewQueryItemsPager(query, partitionKey, nil)
	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, bytes := range response.Items {
			document := documentdb.Document{
				Data: bytes,
			}

			documents = append(documents, document)
		}
	}

	return &documents, nil
}

func (db *CosmosDocumentDatabaseClient) GetDocument(collection string, documentId string) (*documentdb.Document, error) {
	document := documentdb.Document{}
	
	container, err := db.client.NewContainer(collection)
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(collection)

	response, err := container.ReadItem(context.TODO(), partitionKey, documentId, nil)
	if err != nil {
		var cosmosError *azcore.ResponseError
		if errors.As(err, &cosmosError) && cosmosError.StatusCode == 404 {
			return nil, documentdb.ErrDocumentNotFound
		}

		return nil, err
	}

	document.Data = response.Value
	return &document, nil
}

func (db *CosmosDocumentDatabaseClient) UpsertDocument(collection string, documentId string, document *documentdb.Document) error {
	container, err := db.client.NewContainer(collection)
	if err != nil {
		return err
	}

	partitionKey := azcosmos.NewPartitionKeyString(collection)

	if _, err := container.UpsertItem(context.TODO(), partitionKey, document.Data, nil); err != nil {
		return err
	}

	return nil
}

func (db *CosmosDocumentDatabaseClient) DeleteDocument(collection string, documentId string) error {
	container, err := db.client.NewContainer(collection)
	if err != nil {
		return err
	}

	partitionKey := azcosmos.NewPartitionKeyString(collection)

	if _, err := container.DeleteItem(context.TODO(), partitionKey, documentId, nil); err != nil {
		var cosmosError *azcore.ResponseError
		if errors.As(err, &cosmosError) && cosmosError.StatusCode == 404 {
			return documentdb.ErrDocumentNotFound
		}

		return err
	}

	return nil
}
