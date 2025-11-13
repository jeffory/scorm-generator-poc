package scorm

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

// Generator generates SCORM packages
type Generator struct {
	outputDir string
}

// NewGenerator creates a new SCORM generator
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// Generate creates a SCORM package from a quiz
func (g *Generator) Generate(quiz *models.Quiz) (string, error) {
	// Create a temporary directory for the SCORM package
	tempDir, err := os.MkdirTemp("", "scorm-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate all SCORM files
	if err := g.generateManifest(tempDir, quiz); err != nil {
		return "", fmt.Errorf("failed to generate manifest: %w", err)
	}

	if err := g.generateIndexHTML(tempDir, quiz); err != nil {
		return "", fmt.Errorf("failed to generate index.html: %w", err)
	}

	if err := g.generateSCORMAPI(tempDir); err != nil {
		return "", fmt.Errorf("failed to generate SCORM API: %w", err)
	}

	if err := g.generateStyles(tempDir); err != nil {
		return "", fmt.Errorf("failed to generate styles: %w", err)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the zip file
	zipPath := filepath.Join(g.outputDir, sanitizeFilename(quiz.Subject)+".zip")
	if err := g.createZip(tempDir, zipPath); err != nil {
		return "", fmt.Errorf("failed to create zip: %w", err)
	}

	return zipPath, nil
}

func sanitizeFilename(s string) string {
	// Simple sanitization - replace spaces and special characters
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	return result
}

func (g *Generator) createZip(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})
}
