package repo

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"monitoring-by-thrift/internal/domain"
)

type DocPG struct{ db *sql.DB }

func NewDocPG(db *sql.DB) *DocPG { return &DocPG{db: db} }

func (r *DocPG) Save(ctx context.Context, d *domain.Transaction) error {
	/*_, err := r.db.ExecContext(ctx,
		`INSERT INTO documents (id, owner_id, name, path, size, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		d.ID, d.OwnerID, d.Name, d.Path, d.Size, d.CreatedAt)
	return err*/
	return nil
}

func (r *DocPG) Get(ctx context.Context, ids *[]string, out chan<- []domain.Transaction) *error {
	const q = `SELECT id, proxy_params, proxy_params2 FROM transactions WHERE id = ANY($1)`
	rows, err := r.db.QueryContext(ctx, q, ids)
	if err != nil {
		return &err
	}
	var chunk []domain.Transaction
	defer rows.Close()
	for rows.Next() {
		var t domain.Transaction
		var proxy_params, proxy_params2 string
		if err := rows.Scan(&t.ID, &proxy_params, &proxy_params2); err != nil {
			rows.Close()
			return &err
		}
		chunk = append(chunk, t)
	}

	select {
	case out <- chunk:
	case <-ctx.Done():
		err = ctx.Err()
		return &err
	}

	return nil

}

/*
func (r *DocPG) List(ctx context.Context, owner string, limit, offset int) ([]*domain.Transaction, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, owner_id, name, path, size, created_at
		   FROM documents
		  WHERE ($1='' OR owner_id=$1)
		  ORDER BY created_at DESC
		  LIMIT $2 OFFSET $3`, owner, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Transaction
	for rows.Next() {
		var d domain.Transaction
		if err := rows.Scan(&d.ID, &d.OwnerID, &d.Name, &d.Path, &d.Size, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}*/
