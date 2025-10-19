package domain

import (
	"context"
	"io"
)

type TransactionID string

type Transaction struct {
	ID           string
	Status       string
	ProxyParams  string
	ProxyParams2 string
}

type TransactionRepository interface {
	Save(ctx context.Context, trans *Transaction) error
	Get(ctx context.Context, ids *[]string, out chan<- []Transaction) *error
}

type FileStorage interface {
	WriteRow(transactions *[]Transaction) error
	CreateFile(name string) (string, error)
}
type FileParser interface {
	StreamIDs(ctx context.Context, r io.Reader, out chan<- []string) error
}

type ExternalService interface {
	// пример вызова Thrift-сервиса (например, валидация/обогащение)
	ProcessStatus(ctx context.Context, transactionID string, proxyParams *map[string]string) (transactionStatus string, err error)
}
