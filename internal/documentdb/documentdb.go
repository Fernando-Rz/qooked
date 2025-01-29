package documentdb

import "errors"

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrDocumentNotFound   = errors.New("document not found")
)

type Document struct {
	Data []byte
}

type DocumentDatabaseClient interface {
	InitializeClient(endpointUrl string, databaseName string) error
	TestConnection() error
	GetDocuments(collection string) (*[]Document, error)
	GetDocument(collection string, documentId string) (*Document, error)
	UpsertDocument(collection string, documentId string, document *Document) error
	DeleteDocument(collection string, documentId string) error
}
