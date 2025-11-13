package scorm

import (
	"os"
	"path/filepath"
)

func (g *Generator) generateStyles(dir string) error {
	cssContent := `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
    padding: 20px;
}

.container {
    max-width: 800px;
    margin: 0 auto;
    background: white;
    border-radius: 10px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
    padding: 40px;
}

h1 {
    color: #333;
    margin-bottom: 30px;
    text-align: center;
    font-size: 2em;
}

h2 {
    color: #667eea;
    margin-bottom: 20px;
}

h3 {
    color: #555;
    margin-bottom: 15px;
}

h4 {
    color: #666;
    margin-bottom: 10px;
}

.question {
    margin-bottom: 30px;
}

.question-text {
    font-size: 1.2em;
    color: #333;
    margin-bottom: 20px;
    line-height: 1.6;
}

.options {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.option {
    display: flex;
    align-items: center;
    padding: 15px;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.option:hover {
    border-color: #667eea;
    background-color: #f8f9ff;
}

.option input[type="radio"] {
    margin-right: 12px;
    width: 20px;
    height: 20px;
    cursor: pointer;
}

.option span {
    flex: 1;
    color: #333;
    font-size: 1em;
}

.option.correct {
    border-color: #4caf50;
    background-color: #e8f5e9;
}

.option.incorrect {
    border-color: #f44336;
    background-color: #ffebee;
}

#navigation {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 30px;
    padding-top: 20px;
    border-top: 2px solid #e0e0e0;
}

button {
    padding: 12px 24px;
    font-size: 1em;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-weight: 600;
}

#prev-btn, #next-btn {
    background-color: #667eea;
    color: white;
}

#prev-btn:hover:not(:disabled), #next-btn:hover:not(:disabled) {
    background-color: #5568d3;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
    opacity: 0.6;
}

#submit-btn {
    background-color: #4caf50;
    color: white;
    width: 100%;
    margin-top: 20px;
    font-size: 1.1em;
}

#submit-btn:hover {
    background-color: #45a049;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(76, 175, 80, 0.4);
}

#question-counter {
    color: #666;
    font-weight: 500;
}

#results {
    text-align: center;
}

#score {
    font-size: 1.5em;
    color: #667eea;
    margin-bottom: 30px;
    font-weight: 600;
}

.review-item {
    text-align: left;
    margin-bottom: 20px;
    padding: 20px;
    border-radius: 8px;
    border-left: 4px solid #ccc;
}

.review-item.correct-answer {
    background-color: #e8f5e9;
    border-left-color: #4caf50;
}

.review-item.wrong-answer {
    background-color: #ffebee;
    border-left-color: #f44336;
}

.review-item p {
    margin: 10px 0;
    line-height: 1.6;
}

.review-item strong {
    color: #333;
}

@media (max-width: 600px) {
    .container {
        padding: 20px;
    }

    h1 {
        font-size: 1.5em;
    }

    .question-text {
        font-size: 1em;
    }

    #navigation {
        flex-direction: column;
        gap: 10px;
    }

    #prev-btn, #next-btn {
        width: 100%;
    }
}
`

	return os.WriteFile(filepath.Join(dir, "styles.css"), []byte(cssContent), 0644)
}
