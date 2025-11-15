package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

// Storage handles persistent storage of quizzes
type Storage struct {
	dataDir string
}

// QuizMetadata contains information about a saved quiz
type QuizMetadata struct {
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Questions int       `json:"questions"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// New creates a new storage instance
func New(dataDir string) *Storage {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create data directory: %v", err))
	}

	return &Storage{
		dataDir: dataDir,
	}
}

// Save saves a quiz to disk
func (s *Storage) Save(name string, quiz *models.Quiz) error {
	// Sanitize filename
	filename := sanitizeFilename(name) + ".json"
	path := filepath.Join(s.dataDir, filename)

	// Create file
	data, err := json.MarshalIndent(quiz, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal quiz: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load loads a quiz from disk
func (s *Storage) Load(name string) (*models.Quiz, error) {
	filename := sanitizeFilename(name) + ".json"
	path := filepath.Join(s.dataDir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var quiz models.Quiz
	if err := json.Unmarshal(data, &quiz); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quiz: %w", err)
	}

	return &quiz, nil
}

// List returns metadata for all saved quizzes
func (s *Storage) List() ([]QuizMetadata, error) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var quizzes []QuizMetadata
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Load quiz to get metadata
		name := strings.TrimSuffix(entry.Name(), ".json")
		quiz, err := s.Load(name)
		if err != nil {
			continue
		}

		quizzes = append(quizzes, QuizMetadata{
			Name:      name,
			Subject:   quiz.Subject,
			Questions: len(quiz.Questions),
			CreatedAt: info.ModTime(),
			UpdatedAt: info.ModTime(),
		})
	}

	// Sort by updated time (newest first)
	sort.Slice(quizzes, func(i, j int) bool {
		return quizzes[i].UpdatedAt.After(quizzes[j].UpdatedAt)
	})

	return quizzes, nil
}

// Delete removes a quiz from disk
func (s *Storage) Delete(name string) error {
	filename := sanitizeFilename(name) + ".json"
	path := filepath.Join(s.dataDir, filename)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// sanitizeFilename converts a string into a safe filename
func sanitizeFilename(s string) string {
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")

	// Keep only alphanumeric, dash, and underscore
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
