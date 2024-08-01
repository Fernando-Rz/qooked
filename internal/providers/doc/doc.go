package doc

type Document struct{
	Data []byte
}

type DocumentDatabaseClient interface {
	InitializeClient(endpoint string) error
	GetDocuments(collection string) (*[]Document, error)
	GetDocument(collection string, documentId string) (*Document, error)
	UpsertDocument(collection string, documentId string, document *Document) error
	DeleteDocument(collection string, documentId string) error
}