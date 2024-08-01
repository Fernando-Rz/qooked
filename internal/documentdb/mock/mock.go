package mock

import (
	"qooked/internal/documentdb"
	"sync"
)

type MockDocumentDatabaseClient struct {
	mutex       sync.RWMutex
    collections map[string]map[string]documentdb.Document
}

func NewMockDocumentDatabaseClient() *MockDocumentDatabaseClient {
	return &MockDocumentDatabaseClient{
		collections: make(map[string]map[string]documentdb.Document),
	}
}

func (db *MockDocumentDatabaseClient) InitializeClient(endpoint string) error {
	// No initialization required for in-memory implementation
	return nil
}

func (db *MockDocumentDatabaseClient) GetDocuments(collection string) (*[]documentdb.Document, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	documents, ok := db.collections[collection]
	if !ok {
		return nil, documentdb.ErrCollectionNotFound
	}

	var result []documentdb.Document
	for _, doc := range documents {
		result = append(result, doc)
	}

	return &result, nil
}

func (db *MockDocumentDatabaseClient) GetDocument(collection string, documentId string) (*documentdb.Document, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	documents, ok := db.collections[collection]
	if !ok {
		return nil, documentdb.ErrCollectionNotFound
	}

	document, ok := documents[documentId]
	if !ok {
		return nil, documentdb.ErrDocumentNotFound
	}

	return &document, nil
}

func (db *MockDocumentDatabaseClient) UpsertDocument(collection string, documentId string, document *documentdb.Document) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.collections[collection]; !ok {
		db.collections[collection] = make(map[string]documentdb.Document)
	}

	db.collections[collection][documentId] = *document
	return nil
}

func (db *MockDocumentDatabaseClient) DeleteDocument(collection string, documentId string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	documents, ok := db.collections[collection]
	if !ok {
		return documentdb.ErrCollectionNotFound
	}

	if _, ok := documents[documentId]; !ok {
		return documentdb.ErrDocumentNotFound
	}

	delete(documents, documentId)
	return nil
}
