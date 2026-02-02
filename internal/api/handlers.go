package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaya771/kozm/internal/encoder"
)

type Server struct {
	DB *pgxpool.Pool
	Node *snowflake.Node
}

func (s *Server) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id := s.Node.Generate()
	code := encoder.Encode(uint64(id.Int64()))

	_, err := s.DB.Exec(r.Context(),
		"INSERT INTO links (id, short_code, original_url) VALUES ($1, $2, $3)",
		id.Int64(), code, req.URL)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"short_code": code,
		"original":   req.URL,
		"link":       fmt.Sprintf("http://localhost:8080/%s", code),
	})
}


func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	var originalURL string

	err := s.DB.QueryRow(r.Context(),
		"SELECT original_url FROM links WHERE short_code = $1",
		code).Scan(&originalURL)

	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}