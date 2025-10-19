// internal/usecase/document_service.go
package usecase

import (
	"context"
	"fmt"
	"io"
	"monitoring-by-thrift/internal/adapter/repo"
	"monitoring-by-thrift/internal/adapter/thriftclient"
	"monitoring-by-thrift/internal/shared"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"monitoring-by-thrift/internal/domain"
)

type TransactionService struct {
	repo      domain.TransactionRepository
	files     domain.FileStorage
	ext       domain.ExternalService
	csvParser domain.FileParser
	nowFn     func() time.Time
}

func NewDocumentService(r *repo.DocPG, f domain.FileStorage, e *thriftclient.Client, parser domain.FileParser) *TransactionService {
	return &TransactionService{repo: r, files: f, ext: e, csvParser: parser, nowFn: time.Now}
}

func (s *TransactionService) Parser(ctx context.Context, fileReader io.Reader, closeFn func(), providerID string, flusher http.Flusher) (*int64, error) {
	var counter *int64
	defer closeFn()

	// 2) конвейер: CSV → батчи id → воркеры DB → вывод NDJSON
	ids := make(chan []string, shared.InflightBatches)
	results := make(chan []domain.Transaction, shared.InflightBatches)
	errs := make(chan error, 1)
	s.files.CreateFile(fmt.Sprintf("%s.csv", time.Now().Format("02.01.2006_15:04:05")))

	var wgProd sync.WaitGroup
	wgProd.Add(1)
	go func() {
		defer wgProd.Done()
		defer close(ids)
		if err := s.csvParser.StreamIDs(ctx, fileReader, ids); err != nil {
			select {
			case errs <- err:
			default:
			}
		}
	}()

	var wgWorkers sync.WaitGroup
	for i := 0; i < shared.Workers; i++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			if err := s.fetchBatches(ctx, ids, results); err != nil {
				select {

				case errs <- err:
				default:
				}
			}
		}()
	}
	go func() { wgWorkers.Wait(); close(results) }()

loop:
	for {
		select {
		case err := <-errs:
			if err != nil {
				return nil, err
			}
		case chunk, ok := <-results:
			if !ok {
				break loop
			}
			for i := range chunk {
				proxyParams := make(map[string]string)

				id := chunk[i].ID
				providerStatus, err := s.ext.ProcessStatus(ctx, id, &proxyParams) // Thrift вызов
				if err != nil {
					return nil, err
				}
				chunk[i].Status = providerStatus

			}
			s.files.WriteRow(&chunk)

			atomic.AddInt64(counter, int64(len(chunk)))

			if flusher != nil {
				flusher.Flush()
			}
		case <-ctx.Done():
			return counter, nil
		}
	}

	return counter, nil
}

func (s *TransactionService) fetchBatches(ctx context.Context, in <-chan []string, out chan<- []domain.Transaction) error {
	// pgx stdlib понимает массивы int8 через pq-style ($1)
	// Запрос под индекс PK(id) — обязателен для скорости.

	for ids := range in {
		s.repo.Get(ctx, &ids, out)

		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
