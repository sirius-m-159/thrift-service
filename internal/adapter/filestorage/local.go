package filestorage

import (
	"encoding/csv"
	"fmt"
	"monitoring-by-thrift/internal/domain"
	"os"
	"path/filepath"
	"time"
)

type LocalFS struct{ base string }

func NewLocalFS(base string) *LocalFS { return &LocalFS{base: base} }

func (l *LocalFS) CreateFile(name string) (string, error) {
	if err := os.MkdirAll(l.base, 0o775); err != nil {
		return "", err
	}
	fn := fmt.Sprintf("%d_%s", time.Now().UnixNano(), name)
	full := filepath.Join(l.base, fn)

	f, err := os.Create(full)
	if err != nil {
		panic(err)
	}
	return f.Name(), nil

}
func (l *LocalFS) WriteRow(transactions *[]domain.Transaction) error {
	f, _ := os.OpenFile(l.base, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	w := csv.NewWriter(f)
	defer func() { w.Flush(); _ = w.Error() }()
	for _, t := range *transactions { // где-то генерируешь []string
		if err := w.Write([]string{t.ID, t.Status}); err != nil {
			return err
		}
	}

	return nil
}
