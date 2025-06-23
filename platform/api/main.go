package main

import (
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/rasha-hantash/boilerplate-repo/platform/api/config"
	"github.com/rasha-hantash/boilerplate-repo/platform/api/handler"
	todov1connect "github.com/rasha-hantash/boilerplate-repo/platform/api/gen/proto/todo/v1/todov1connect"
	services "github.com/rasha-hantash/boilerplate-repo/platform/api/services/todo"
)

func main() {
	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := services.NewRepository(db)

	// Initialize database schema
	if err := repo.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize handler
	todoHandler := handler.NewTodoHandler(repo)

	// Set up HTTP server with CORS
	mux := http.NewServeMux()
	// Create ConnectRPC service
	mux.Handle(todov1connect.NewTodoServiceHandler(todoHandler))

	// Add CORS middleware
	corsHandler := corsMiddleware(mux)

	// Create server with HTTP/2 support
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),
	}

	log.Printf("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers to allow frontend requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
