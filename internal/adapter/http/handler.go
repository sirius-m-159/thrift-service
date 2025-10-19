package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"monitoring-by-thrift/internal/usecase"
)

type Handler struct {
	s *usecase.TransactionService
}

func NewHandler(s *usecase.TransactionService) *Handler { return &Handler{s: s} }

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(Recoverer, RequestID) // см. middleware.go
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })

	r.Route("/v1/transactions", func(r chi.Router) {
		r.Post("/", h.createStatuses)
		//r.Get("/{id}", h.get)
		//r.Get("/", h.list)
	})
	return r
}

func (h *Handler) createStatuses(w http.ResponseWriter, r *http.Request) {

	fileReader, fileCloser, err := h.csvFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	providerID := r.FormValue("provider_id")
	w.Header().Set("Content-Type", "application/x-ndjson")

	flusher, _ := w.(http.Flusher)

	doc, err := h.s.Parser(r.Context(), fileReader, fileCloser, providerID, flusher)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(doc)
}

func (h *Handler) csvFromRequest(r *http.Request) (io.Reader, func(), error) {
	ct := r.Header.Get("Content-Type")
	// multipart/form-data: file=<csv>
	if strings.HasPrefix(ct, "multipart/form-data") {
		// ограничим память на парсинг формы
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
			return nil, func() {}, err
		}
		f, _, err := r.FormFile("file")
		if err != nil {
			return nil, func() {}, err
		}
		return f, func() { _ = f.Close() }, nil
	}
	// сырое тело text/csv
	if strings.Contains(ct, "text/csv") || ct == "" {
		return r.Body, func() { _ = r.Body.Close() }, nil
	}
	return nil, func() {}, errors.New("unsupported Content-Type: " + ct)
}
