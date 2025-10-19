FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/server

FROM gcr.io/distroless/base-debian12
ENV HTTP_ADDR=:8080 DATA_DIR=/data/files
WORKDIR /app
COPY --from=build /bin/app /app/app
VOLUME ["/data/files"]
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/app"]
