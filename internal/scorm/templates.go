package scorm

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"

	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

func (g *Generator) generateIndexHTML(dir string, quiz *models.Quiz) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}} Quiz</title>
    <link rel="stylesheet" href="styles.css">
    <script src="scormapi.js"></script>
</head>
<body>
    <div class="container">
        <h1>{{.Subject}} Quiz</h1>
        <div id="quiz-container">
            <div id="question-container"></div>
            <div id="navigation">
                <button id="prev-btn" onclick="previousQuestion()" disabled>Previous</button>
                <span id="question-counter"></span>
                <button id="next-btn" onclick="nextQuestion()">Next</button>
            </div>
            <button id="submit-btn" onclick="submitQuiz()" style="display:none;">Submit Quiz</button>
        </div>
        <div id="results" style="display:none;">
            <h2>Quiz Results</h2>
            <p id="score"></p>
            <div id="review"></div>
        </div>
    </div>

    <script>
        const quizData = {{.QuestionsJSON}};
        let currentQuestion = 0;
        let userAnswers = new Array(quizData.length).fill(null);
        let quizSubmitted = false;

        // Initialize SCORM
        window.onload = function() {
            initSCORM();
            showQuestion(currentQuestion);
        };

        window.onbeforeunload = function() {
            if (!quizSubmitted) {
                setSCORMValue('cmi.core.lesson_status', 'incomplete');
                setSCORMValue('cmi.core.exit', 'suspend');
            }
            finishSCORM();
        };

        function showQuestion(index) {
            const question = quizData[index];
            const container = document.getElementById('question-container');

            let html = '<div class="question">';
            html += '<h3>Question ' + (index + 1) + '</h3>';
            html += '<p class="question-text">' + escapeHtml(question.question) + '</p>';
            html += '<div class="options">';

            question.options.forEach((option, i) => {
                const checked = userAnswers[index] === i ? 'checked' : '';
                const disabled = quizSubmitted ? 'disabled' : '';
                html += '<label class="option ' + (quizSubmitted && i === question.answer ? 'correct' : '') +
                        (quizSubmitted && userAnswers[index] === i && i !== question.answer ? 'incorrect' : '') + '">';
                html += '<input type="radio" name="answer" value="' + i + '" ' + checked + ' ' + disabled +
                        ' onchange="saveAnswer(' + i + ')">';
                html += '<span>' + escapeHtml(option) + '</span>';
                html += '</label>';
            });

            html += '</div></div>';
            container.innerHTML = html;

            updateNavigation();
        }

        function updateNavigation() {
            document.getElementById('prev-btn').disabled = currentQuestion === 0;
            document.getElementById('next-btn').disabled = currentQuestion === quizData.length - 1 || quizSubmitted;
            document.getElementById('question-counter').textContent =
                'Question ' + (currentQuestion + 1) + ' of ' + quizData.length;

            if (currentQuestion === quizData.length - 1 && !quizSubmitted) {
                document.getElementById('submit-btn').style.display = 'block';
            } else {
                document.getElementById('submit-btn').style.display = 'none';
            }
        }

        function saveAnswer(answerIndex) {
            userAnswers[currentQuestion] = answerIndex;
        }

        function nextQuestion() {
            if (currentQuestion < quizData.length - 1) {
                currentQuestion++;
                showQuestion(currentQuestion);
            }
        }

        function previousQuestion() {
            if (currentQuestion > 0) {
                currentQuestion--;
                showQuestion(currentQuestion);
            }
        }

        function submitQuiz() {
            if (!confirm('Are you sure you want to submit your quiz?')) {
                return;
            }

            quizSubmitted = true;
            let correctCount = 0;

            quizData.forEach((question, index) => {
                if (userAnswers[index] === question.answer) {
                    correctCount++;
                }
            });

            const score = Math.round((correctCount / quizData.length) * 100);

            // Update SCORM
            setSCORMValue('cmi.core.score.raw', score);
            setSCORMValue('cmi.core.score.min', 0);
            setSCORMValue('cmi.core.score.max', 100);
            setSCORMValue('cmi.core.lesson_status', score >= 70 ? 'passed' : 'failed');

            // Show results
            document.getElementById('quiz-container').style.display = 'none';
            document.getElementById('results').style.display = 'block';
            document.getElementById('score').textContent =
                'You scored ' + correctCount + ' out of ' + quizData.length + ' (' + score + '%)';

            // Show review
            showReview();
        }

        function showReview() {
            const reviewDiv = document.getElementById('review');
            let html = '<h3>Review Your Answers</h3>';

            quizData.forEach((question, index) => {
                const correct = userAnswers[index] === question.answer;
                html += '<div class="review-item ' + (correct ? 'correct-answer' : 'wrong-answer') + '">';
                html += '<h4>Question ' + (index + 1) + '</h4>';
                html += '<p>' + escapeHtml(question.question) + '</p>';
                html += '<p><strong>Your answer:</strong> ' +
                        (userAnswers[index] !== null ? escapeHtml(question.options[userAnswers[index]]) : 'Not answered') + '</p>';
                if (!correct) {
                    html += '<p><strong>Correct answer:</strong> ' + escapeHtml(question.options[question.answer]) + '</p>';
                }
                html += '</div>';
            });

            reviewDiv.innerHTML = html;
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    </script>
</body>
</html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		return err
	}

	questionsJSON, err := json.Marshal(quiz.Questions)
	if err != nil {
		return err
	}

	data := struct {
		Subject       string
		QuestionsJSON template.JS
	}{
		Subject:       quiz.Subject,
		QuestionsJSON: template.JS(questionsJSON),
	}

	file, err := os.Create(filepath.Join(dir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, data)
}
