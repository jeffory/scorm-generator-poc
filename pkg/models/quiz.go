package models

// Question represents a multiple-choice question
type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"` // Index of correct answer in Options
}

// Quiz represents a collection of questions on a subject
type Quiz struct {
	Subject   string     `json:"subject"`
	Questions []Question `json:"questions"`
}
