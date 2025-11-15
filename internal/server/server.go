package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"sync"

	"github.com/jeffory/scorm-generator-poc/internal/openrouter"
	"github.com/jeffory/scorm-generator-poc/internal/storage"
	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

//go:embed web/*
var webFiles embed.FS

// Server handles HTTP requests for the SCORM generator
type Server struct {
	mux       *http.ServeMux
	orClient  *openrouter.Client
	storage   *storage.Storage
	sessions  map[string]*models.Quiz
	sessionMu sync.RWMutex
}

// New creates a new server instance
func New(apiKey, model string) *Server {
	orClient := openrouter.NewClient(apiKey)
	orClient.Model = model

	s := &Server{
		mux:      http.NewServeMux(),
		orClient: orClient,
		storage:  storage.New("./data"),
		sessions: make(map[string]*models.Quiz),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Serve static files from embedded filesystem
	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	s.mux.Handle("/", http.FileServer(http.FS(webFS)))

	// API endpoints
	s.mux.HandleFunc("/api/generate", s.handleGenerate)
	s.mux.HandleFunc("/api/regenerate-question", s.handleRegenerateQuestion)
	s.mux.HandleFunc("/api/update-question", s.handleUpdateQuestion)
	s.mux.HandleFunc("/api/generate-scorm", s.handleGenerateSCORM)
	s.mux.HandleFunc("/api/library", s.handleLibrary)
	s.mux.HandleFunc("/api/library/save", s.handleSaveQuiz)
	s.mux.HandleFunc("/api/library/load", s.handleLoadQuiz)
	s.mux.HandleFunc("/api/library/delete", s.handleDeleteQuiz)
	s.mux.HandleFunc("/api/export", s.handleExport)
	s.mux.HandleFunc("/api/import", s.handleImport)

	// SCORM package serving
	s.mux.HandleFunc("/packages/", s.handleServePackage)
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
