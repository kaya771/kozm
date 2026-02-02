package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaya771/kozm/internal/api"
)

func main() {
	nodeIDStr := os.Getenv("NODE_ID")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)
	if nodeID == 0 {
		nodeID = 1
	}

	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to create snowflake node: %v\n", err)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	setupDatabase(pool)

	srv := &api.Server{
		DB: pool,
		Node: node,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(5, 1*time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://kozm.link", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Post("/shorten", srv.Shorten)
	r.Get("/{code}", srv.Redirect)

	fmt.Printf("Kozm Node %d listening on :8080...\n", nodeID)
	http.ListenAndServe(":8080", r)
}

func setupDatabase(pool *pgxpool.Pool) {
	schema := `
	CREATE TABLE IF NOT EXISTS links (
		id BIGINT PRIMARY KEY,
		short_code TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := pool.Exec(context.Background(), schema)
	if err != nil {
		fmt.Printf("Failed to initialize schema: %v\n", err)
	}
}