package shared

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPAddr   string // ":8080"
	DataDir    string // "/data/files"
	DSN        string // "postgres://user:pass@host:5432/db?sslmode=disable"
	ThriftAddr string // "ext-svc:9090"
	ThriftTO   time.Duration
}

func MustLoad() *Config {
	c := &Config{
		HTTPAddr:   getenv("HTTP_ADDR", ":8080"),
		DataDir:    getenv("DATA_DIR", "/data/files"),
		DSN:        getenv("DB_DSN", ""),
		ThriftAddr: getenv("THRIFT_ADDR", "localhost:9090"),
		ThriftTO:   mustDuration(getenv("THRIFT_TIMEOUT", "2s")),
	}
	if c.DSN == "" {
		log.Fatal("DB_DSN required")
	}
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustDuration(s string) time.Duration { d, _ := time.ParseDuration(s); return d }
