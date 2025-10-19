run:
	DB_DSN="postgres://user:pass@host:5432/db?sslmode=disable" \
	THRIFT_ADDR="ext-svc:9090" \
	DATA_DIR="./data" \
	go run ./cmd/server