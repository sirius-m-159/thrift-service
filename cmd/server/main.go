// cmd/server/main.go
package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"monitoring-by-thrift/internal/adapter/filestorage"
	httpad "monitoring-by-thrift/internal/adapter/http"
	"monitoring-by-thrift/internal/adapter/repo"
	"monitoring-by-thrift/internal/adapter/thriftclient"
	"monitoring-by-thrift/internal/shared"
	"monitoring-by-thrift/internal/usecase"
)

func main() {
	cfg := shared.MustLoad()

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	fs := filestorage.NewLocalFS(cfg.DataDir)
	ext, err := thriftclient.New(cfg.ThriftAddr, cfg.ThriftTO)
	if err != nil {
		log.Fatal(err)
	}

	rp := repo.NewDocPG(db)
	parser := filestorage.NewCsvParser(',')
	svc := usecase.NewDocumentService(rp, fs, ext, parser)

	h := httpad.NewHandler(svc)
	log.Printf("HTTP listen on %s", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, h.Routes()))
}
