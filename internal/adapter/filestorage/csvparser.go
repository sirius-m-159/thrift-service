package filestorage

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"monitoring-by-thrift/internal/shared"
	"strings"
)

type CsvParser struct{ splitter rune }

func NewCsvParser(splitter rune) *CsvParser { return &CsvParser{splitter: splitter} }

func (p *CsvParser) StreamIDs(ctx context.Context, r io.Reader, out chan<- []string) error {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	cr.FieldsPerRecord = -1
	cr.Comma = p.splitter

	rec, err := cr.Read()
	if err != nil {
		return err
	}
	idIdx, hasHeader := findIDIndex(rec)

	var batch []string
	emit := func() {
		if len(batch) == 0 {
			return
		}
		cp := make([]string, len(batch))
		copy(cp, batch)
		select {
		case out <- cp:
		case <-ctx.Done():
		}
		batch = batch[:0]
	}

	// если заголовка нет — первая строка содержит данные
	if !hasHeader {
		if id, ok := parseID(rec, idIdx); ok {
			batch = append(batch, id)
		}
		if len(batch) >= shared.BatchSize {
			emit()
		}
	}

	for {
		rec, err = cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if id, ok := parseID(rec, idIdx); ok {
			batch = append(batch, id)
			if len(batch) >= shared.BatchSize {
				emit()
			}
		}
	}
	emit()
	return nil
}

func findIDIndex(rec []string) (idx int, hasHeader bool) {
	for i, v := range rec {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "id", "tx_id", "transaction_id":
			return i, true
		}
	}
	return 0, false
}

func parseID(rec []string, idx int) (string, bool) {
	if idx < 0 || idx >= len(rec) {
		return "", false
	}

	return strings.TrimSpace(rec[idx]), true
}
