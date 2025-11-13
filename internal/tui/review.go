package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	questionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	optionStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			MarginTop(0)

	correctStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(2)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000")).
			MarginTop(1)
)

type Model struct {
	quiz          *models.Quiz
	currentIndex  int
	approved      bool
	rejected      bool
	width         int
	height        int
}

func NewModel(quiz *models.Quiz) Model {
	return Model{
		quiz:         quiz,
		currentIndex: 0,
		approved:     false,
		rejected:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.rejected = true
			return m, tea.Quit

		case "left", "h":
			if m.currentIndex > 0 {
				m.currentIndex--
			}
			return m, nil

		case "right", "l":
			if m.currentIndex < len(m.quiz.Questions)-1 {
				m.currentIndex++
			}
			return m, nil

		case "enter", "y":
			m.approved = true
			return m, tea.Quit

		case "n":
			m.rejected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("Quiz: %s", m.quiz.Subject)))
	b.WriteString("\n\n")

	// Current question
	question := m.quiz.Questions[m.currentIndex]

	// Progress
	progress := fmt.Sprintf("Question %d of %d", m.currentIndex+1, len(m.quiz.Questions))
	b.WriteString(progressStyle.Render(progress))
	b.WriteString("\n\n")

	// Question text
	b.WriteString(questionStyle.Render(fmt.Sprintf("Q%d: %s", question.ID, question.Question)))
	b.WriteString("\n\n")

	// Options
	b.WriteString("Options:\n\n")
	for i, option := range question.Options {
		prefix := fmt.Sprintf("%c) ", 'A'+i)
		if i == question.Answer {
			b.WriteString(correctStyle.Render(prefix + option + " ✓ (correct)"))
		} else {
			b.WriteString(optionStyle.Render(prefix + option))
		}
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Navigation: ← → or h/l  •  Approve: Enter or y  •  Reject: n or q"))

	return b.String()
}

func (m Model) Approved() bool {
	return m.approved
}

func (m Model) Rejected() bool {
	return m.rejected
}

// ReviewQuiz displays the quiz in a TUI and returns true if approved
func ReviewQuiz(quiz *models.Quiz) (bool, error) {
	m := NewModel(quiz)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running TUI: %w", err)
	}

	result := finalModel.(Model)
	return result.Approved(), nil
}
