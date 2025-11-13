package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jeffory/scorm-generator-poc/internal/openrouter"
	"github.com/jeffory/scorm-generator-poc/internal/quiz"
	"github.com/jeffory/scorm-generator-poc/internal/scorm"
)

func main() {
	// Command-line flags
	subject := flag.String("subject", "", "Subject for the quiz (required)")
	numQuestions := flag.Int("questions", 5, "Number of questions to generate")
	outputDir := flag.String("output", "./output", "Output directory for SCORM package")
	apiKey := flag.String("api-key", "", "OpenRouter API key (or set OPENROUTER_API_KEY env var)")
	model := flag.String("model", openrouter.DefaultModel, "OpenRouter model to use")
	flag.Parse()

	// Validate inputs
	if *subject == "" {
		fmt.Println("Error: -subject is required")
		flag.Usage()
		os.Exit(1)
	}

	// Get API key from flag or environment
	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENROUTER_API_KEY")
	}

	if key == "" {
		fmt.Println("Error: OpenRouter API key is required")
		fmt.Println("Provide it via -api-key flag or OPENROUTER_API_KEY environment variable")
		os.Exit(1)
	}

	// Initialize OpenRouter client
	fmt.Println("Initializing OpenRouter client...")
	client := openrouter.NewClient(key)
	client.Model = *model

	// Generate quiz
	fmt.Printf("Generating %d questions about '%s' using model %s...\n", *numQuestions, *subject, *model)
	quizGen := quiz.NewGenerator(client)
	quizData, err := quizGen.Generate(*subject, *numQuestions)
	if err != nil {
		fmt.Printf("Error generating quiz: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %d questions\n", len(quizData.Questions))

	// Generate SCORM package
	fmt.Println("Creating SCORM package...")
	scormGen := scorm.NewGenerator(*outputDir)
	zipPath, err := scormGen.Generate(quizData)
	if err != nil {
		fmt.Printf("Error creating SCORM package: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ SCORM package created successfully!\n")
	fmt.Printf("  Location: %s\n", zipPath)
	fmt.Printf("  Subject: %s\n", quizData.Subject)
	fmt.Printf("  Questions: %d\n", len(quizData.Questions))
	fmt.Println("\nYou can now upload this SCORM package to your LMS (Learning Management System)")
}
