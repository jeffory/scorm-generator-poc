package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jeffory/scorm-generator-poc/internal/quiz"
	"github.com/jeffory/scorm-generator-poc/internal/scorm"
	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

// GenerateRequest represents a request to generate quiz questions
type GenerateRequest struct {
	Subject   string `json:"subject"`
	Questions int    `json:"questions"`
	SessionID string `json:"sessionId"`
}

// RegenerateQuestionRequest represents a request to regenerate a single question
type RegenerateQuestionRequest struct {
	SessionID     string `json:"sessionId"`
	QuestionIndex int    `json:"questionIndex"`
	Subject       string `json:"subject"`
}

// UpdateQuestionRequest represents a request to update a question
type UpdateQuestionRequest struct {
	SessionID     string         `json:"sessionId"`
	QuestionIndex int            `json:"questionIndex"`
	Question      models.Question `json:"question"`
}

// GenerateSCORMRequest represents a request to generate a SCORM package
type GenerateSCORMRequest struct {
	SessionID    string `json:"sessionId"`
	PassingScore int    `json:"passingScore,omitempty"`
}

// SaveQuizRequest represents a request to save a quiz to the library
type SaveQuizRequest struct {
	Name      string       `json:"name"`
	Quiz      models.Quiz  `json:"quiz"`
	SessionID string       `json:"sessionId"`
}

// handleGenerate generates quiz questions based on subject and count
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Subject == "" {
		respondError(w, http.StatusBadRequest, "Subject is required")
		return
	}
	if req.Questions < 1 || req.Questions > 50 {
		respondError(w, http.StatusBadRequest, "Questions must be between 1 and 50")
		return
	}

	// Generate questions
	generator := quiz.NewGenerator(s.orClient)
	quizData, err := generator.Generate(req.Subject, req.Questions)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate questions: %v", err))
		return
	}

	s.sessionMu.Lock()
	s.sessions[req.SessionID] = quizData
	s.sessionMu.Unlock()

	respondJSON(w, http.StatusOK, quizData)
}

// handleRegenerateQuestion regenerates a single question
func (s *Server) handleRegenerateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegenerateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get session
	s.sessionMu.RLock()
	quizData, exists := s.sessions[req.SessionID]
	s.sessionMu.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Session not found")
		return
	}

	if req.QuestionIndex < 0 || req.QuestionIndex >= len(quizData.Questions) {
		respondError(w, http.StatusBadRequest, "Invalid question index")
		return
	}

	// Generate a single new question
	generator := quiz.NewGenerator(s.orClient)
	tempQuiz, err := generator.Generate(req.Subject, 1)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to regenerate question: %v", err))
		return
	}

	if len(tempQuiz.Questions) == 0 {
		respondError(w, http.StatusInternalServerError, "No question generated")
		return
	}

	// Update the question in the session
	newQuestion := tempQuiz.Questions[0]
	newQuestion.ID = req.QuestionIndex + 1

	s.sessionMu.Lock()
	quizData.Questions[req.QuestionIndex] = newQuestion
	s.sessionMu.Unlock()

	respondJSON(w, http.StatusOK, newQuestion)
}

// handleUpdateQuestion updates a question with user edits
func (s *Server) handleUpdateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req UpdateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get session
	s.sessionMu.RLock()
	quizData, exists := s.sessions[req.SessionID]
	s.sessionMu.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Session not found")
		return
	}

	if req.QuestionIndex < 0 || req.QuestionIndex >= len(quizData.Questions) {
		respondError(w, http.StatusBadRequest, "Invalid question index")
		return
	}

	// Update the question
	s.sessionMu.Lock()
	quizData.Questions[req.QuestionIndex] = req.Question
	s.sessionMu.Unlock()

	respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// handleGenerateSCORM generates a SCORM package from the session quiz
func (s *Server) handleGenerateSCORM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req GenerateSCORMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get session
	s.sessionMu.RLock()
	quizData, exists := s.sessions[req.SessionID]
	s.sessionMu.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Session not found")
		return
	}

	// Set passing score if provided
	passingScore := 70
	if req.PassingScore > 0 && req.PassingScore <= 100 {
		passingScore = req.PassingScore
	}

	// Generate SCORM package
	generator := scorm.NewGenerator(passingScore)
	packagePath, err := generator.Generate(quizData)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate SCORM package: %v", err))
		return
	}

	// Return the package path for download
	respondJSON(w, http.StatusOK, map[string]string{
		"packagePath": filepath.Base(packagePath),
		"downloadUrl": "/packages/" + filepath.Base(packagePath),
	})
}

// handleLibrary returns all saved quizzes
func (s *Server) handleLibrary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	quizzes, err := s.storage.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list quizzes: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, quizzes)
}

// handleSaveQuiz saves a quiz to the library
func (s *Server) handleSaveQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SaveQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Quiz name is required")
		return
	}

	if err := s.storage.Save(req.Name, &req.Quiz); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save quiz: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

// handleLoadQuiz loads a quiz from the library
func (s *Server) handleLoadQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "Quiz name is required")
		return
	}

	quizData, err := s.storage.Load(name)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Quiz not found: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, quizData)
}

// handleDeleteQuiz deletes a quiz from the library
func (s *Server) handleDeleteQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "Quiz name is required")
		return
	}

	if err := s.storage.Delete(name); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete quiz: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// handleExport exports the current session quiz as JSON
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	s.sessionMu.RLock()
	quizData, exists := s.sessions[sessionID]
	s.sessionMu.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Session not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", quizData.Subject))
	json.NewEncoder(w).Encode(quizData)
}

// handleImport imports a quiz from JSON
func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	var quizData models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quizData); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid quiz JSON")
		return
	}

	// Validate quiz
	if quizData.Subject == "" || len(quizData.Questions) == 0 {
		respondError(w, http.StatusBadRequest, "Invalid quiz data")
		return
	}

	// Store in session
	s.sessionMu.Lock()
	s.sessions[sessionID] = &quizData
	s.sessionMu.Unlock()

	respondJSON(w, http.StatusOK, quizData)
}

// handleServePackage serves generated SCORM packages for download or preview
func (s *Server) handleServePackage(w http.ResponseWriter, r *http.Request) {
	// Get package name from URL path
	packageName := filepath.Base(r.URL.Path)
	packagePath := filepath.Join("./output", packageName)

	// Check if file exists
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		respondError(w, http.StatusNotFound, "Package not found")
		return
	}

	// Serve the file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", packageName))
	http.ServeFile(w, r, packagePath)
}
