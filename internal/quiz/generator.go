package quiz

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeffory/scorm-generator-poc/internal/openrouter"
	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

// Generator generates quiz questions using an LLM
type Generator struct {
	client *openrouter.Client
}

// NewGenerator creates a new quiz generator
func NewGenerator(client *openrouter.Client) *Generator {
	return &Generator{
		client: client,
	}
}

// Generate generates a quiz on the given subject with the specified number of questions
func (g *Generator) Generate(subject string, numQuestions int) (*models.Quiz, error) {
	prompt := g.buildPrompt(subject, numQuestions)

	response, err := g.client.SendPrompt(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	quiz, err := g.parseQuizResponse(response, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz response: %w", err)
	}

	return quiz, nil
}

func (g *Generator) buildPrompt(subject string, numQuestions int) string {
	return fmt.Sprintf(`Generate %d multiple-choice questions about "%s".

For each question, provide:
1. The question text
2. Four answer options (A, B, C, D)
3. The correct answer (indicate which option is correct)

Return ONLY a valid JSON array with the following structure:
[
  {
    "question": "Question text here?",
    "options": ["Option A", "Option B", "Option C", "Option D"],
    "answer": 0
  }
]

The "answer" field should be the index (0-3) of the correct option.
Make the questions educational and relevant to the subject.
Ensure the JSON is properly formatted and valid.
Do not include any markdown formatting or additional text - only the JSON array.`, numQuestions, subject)
}

func (g *Generator) parseQuizResponse(response, subject string) (*models.Quiz, error) {
	// Clean up the response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Find the JSON array
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]")

	if startIdx == -1 || endIdx == -1 {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	jsonStr := response[startIdx : endIdx+1]

	var questions []models.Question
	if err := json.Unmarshal([]byte(jsonStr), &questions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal questions: %w (response: %s)", err, jsonStr)
	}

	// Assign IDs to questions
	for i := range questions {
		questions[i].ID = i + 1
	}

	quiz := &models.Quiz{
		Subject:   subject,
		Questions: questions,
	}

	return quiz, nil
}
