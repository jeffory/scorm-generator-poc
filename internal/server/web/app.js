// Application State
const state = {
    sessionId: generateSessionId(),
    currentQuiz: null,
    currentPackagePath: null,
    passingScore: 70
};

// Generate a unique session ID
function generateSessionId() {
    return 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    loadLibrary();
});

// Setup all event listeners
function setupEventListeners() {
    // Tab navigation
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });

    // Quiz form submission
    document.getElementById('quiz-form').addEventListener('submit', handleGenerateQuiz);

    // Editor actions
    document.getElementById('create-scorm-btn').addEventListener('click', handleCreateSCORM);
    document.getElementById('save-to-library-btn').addEventListener('click', () => showSaveModal());
    document.getElementById('export-json-btn').addEventListener('click', handleExportJSON);
    document.getElementById('import-json-btn').addEventListener('click', () => document.getElementById('import-file-input').click());
    document.getElementById('import-file-input').addEventListener('change', handleImportJSON);

    // Success actions
    document.getElementById('download-btn').addEventListener('click', handleDownload);
    document.getElementById('preview-btn').addEventListener('click', handlePreview);
    document.getElementById('new-quiz-btn').addEventListener('click', resetApp);

    // Modal close buttons
    document.querySelectorAll('.close-modal').forEach(btn => {
        btn.addEventListener('click', closeModals);
    });

    // Save modal
    document.getElementById('confirm-save-btn').addEventListener('click', handleSaveToLibrary);

    // Close modals on outside click
    document.querySelectorAll('.modal').forEach(modal => {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) closeModals();
        });
    });
}

// Tab switching
function switchTab(tabName) {
    // Update buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.tab === tabName);
    });

    // Update content
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.toggle('active', content.id === tabName + '-tab');
    });

    // Load library if switching to library tab
    if (tabName === 'library') {
        loadLibrary();
    }
}

// Generate quiz
async function handleGenerateQuiz(e) {
    e.preventDefault();

    const subject = document.getElementById('subject').value;
    const questionCount = parseInt(document.getElementById('question-count').value);
    state.passingScore = parseInt(document.getElementById('passing-score').value);

    const btn = document.getElementById('generate-btn');
    const btnText = btn.querySelector('.btn-text');
    const spinner = btn.querySelector('.loading-spinner');

    try {
        // Show loading state
        btn.disabled = true;
        btnText.textContent = 'Generating...';
        spinner.style.display = 'inline-block';

        const response = await fetch('/api/generate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                subject,
                questions: questionCount,
                sessionId: state.sessionId
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to generate quiz');
        }

        const quiz = await response.json();
        state.currentQuiz = quiz;

        // Show quiz editor
        document.getElementById('generation-form').style.display = 'none';
        document.getElementById('quiz-editor').style.display = 'block';
        renderQuestions();

    } catch (error) {
        alert('Error generating quiz: ' + error.message);
    } finally {
        btn.disabled = false;
        btnText.textContent = 'Generate Questions';
        spinner.style.display = 'none';
    }
}

// Render questions in editor
function renderQuestions() {
    const container = document.getElementById('questions-list');
    container.innerHTML = state.currentQuiz.questions.map((q, index) => `
        <div class="question-item" data-index="${index}">
            <div class="question-header">
                <span class="question-number">Question ${index + 1}</span>
                <div class="question-actions">
                    <button class="btn-icon" onclick="editQuestion(${index})" title="Edit">✏️</button>
                    <button class="btn-icon" onclick="regenerateQuestion(${index})" title="Regenerate">🔄</button>
                </div>
            </div>
            <div class="question-display">
                <div class="question-text">${escapeHtml(q.question)}</div>
                <ul class="options-list">
                    ${q.options.map((opt, i) => `
                        <li class="option-item ${i === q.answer ? 'correct' : ''}">
                            <span class="option-label">${String.fromCharCode(65 + i)}.</span>
                            ${escapeHtml(opt)}
                            ${i === q.answer ? '<span class="correct-badge">✓ Correct</span>' : ''}
                        </li>
                    `).join('')}
                </ul>
            </div>
        </div>
    `).join('');
}

// Edit question
function editQuestion(index) {
    const q = state.currentQuiz.questions[index];
    const item = document.querySelector(`[data-index="${index}"]`);

    item.classList.add('editing');
    item.querySelector('.question-display').innerHTML = `
        <div class="form-group">
            <label>Question Text</label>
            <textarea id="edit-q-${index}" class="edit-question-text">${escapeHtml(q.question)}</textarea>
        </div>
        ${q.options.map((opt, i) => `
            <div class="form-group">
                <label>Option ${String.fromCharCode(65 + i)}</label>
                <input type="text" id="edit-opt-${index}-${i}" value="${escapeHtml(opt)}">
            </div>
        `).join('')}
        <div class="form-group">
            <label>Correct Answer</label>
            <select id="edit-answer-${index}">
                ${q.options.map((_, i) => `
                    <option value="${i}" ${i === q.answer ? 'selected' : ''}>
                        Option ${String.fromCharCode(65 + i)}
                    </option>
                `).join('')}
            </select>
        </div>
        <div style="display: flex; gap: 0.5rem; margin-top: 1rem;">
            <button class="btn btn-primary" onclick="saveQuestion(${index})">Save</button>
            <button class="btn btn-secondary" onclick="cancelEdit(${index})">Cancel</button>
        </div>
    `;
}

// Save edited question
async function saveQuestion(index) {
    const question = document.getElementById(`edit-q-${index}`).value;
    const options = [0, 1, 2, 3].map(i =>
        document.getElementById(`edit-opt-${index}-${i}`).value
    );
    const answer = parseInt(document.getElementById(`edit-answer-${index}`).value);

    try {
        const response = await fetch('/api/update-question', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                sessionId: state.sessionId,
                questionIndex: index,
                question: {
                    id: index + 1,
                    question,
                    options,
                    answer
                }
            })
        });

        if (!response.ok) throw new Error('Failed to update question');

        state.currentQuiz.questions[index] = { id: index + 1, question, options, answer };
        renderQuestions();
    } catch (error) {
        alert('Error saving question: ' + error.message);
    }
}

// Cancel edit
function cancelEdit(index) {
    renderQuestions();
}

// Regenerate question
async function regenerateQuestion(index) {
    if (!confirm('Regenerate this question? This will replace it with a new AI-generated question.')) {
        return;
    }

    try {
        const response = await fetch('/api/regenerate-question', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                sessionId: state.sessionId,
                questionIndex: index,
                subject: state.currentQuiz.subject
            })
        });

        if (!response.ok) throw new Error('Failed to regenerate question');

        const newQuestion = await response.json();
        state.currentQuiz.questions[index] = newQuestion;
        renderQuestions();
    } catch (error) {
        alert('Error regenerating question: ' + error.message);
    }
}

// Create SCORM package
async function handleCreateSCORM() {
    const btn = document.getElementById('create-scorm-btn');

    try {
        btn.disabled = true;
        btn.textContent = '📦 Creating Package...';

        const response = await fetch('/api/generate-scorm', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                sessionId: state.sessionId,
                passingScore: state.passingScore
            })
        });

        if (!response.ok) throw new Error('Failed to create SCORM package');

        const data = await response.json();
        state.currentPackagePath = data.downloadUrl;

        // Show success message
        document.getElementById('quiz-editor').style.display = 'none';
        document.getElementById('success-message').style.display = 'block';
        document.getElementById('success-text').textContent =
            `Your SCORM package "${state.currentQuiz.subject}" is ready!`;

    } catch (error) {
        alert('Error creating SCORM package: ' + error.message);
    } finally {
        btn.disabled = false;
        btn.textContent = '📦 Create SCORM Package';
    }
}

// Download package
function handleDownload() {
    if (state.currentPackagePath) {
        window.location.href = state.currentPackagePath;
    }
}

// Preview SCORM
function handlePreview() {
    if (!state.currentPackagePath) return;

    // For preview, we need to extract and serve the SCORM content
    // Since we can't easily unzip client-side, we'll show a message
    alert('Preview functionality: To preview the SCORM package, please:\n\n' +
          '1. Download the package\n' +
          '2. Extract the ZIP file\n' +
          '3. Open index.html in your browser\n\n' +
          'Or upload to your LMS for full SCORM preview.');
}

// Export JSON
function handleExportJSON() {
    if (!state.currentQuiz) return;

    const dataStr = JSON.stringify(state.currentQuiz, null, 2);
    const blob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${state.currentQuiz.subject}.json`;
    a.click();
    URL.revokeObjectURL(url);
}

// Import JSON
async function handleImportJSON(e) {
    const file = e.target.files[0];
    if (!file) return;

    try {
        const text = await file.text();
        const quiz = JSON.parse(text);

        const response = await fetch(`/api/import?sessionId=${state.sessionId}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: text
        });

        if (!response.ok) throw new Error('Failed to import quiz');

        state.currentQuiz = quiz;
        document.getElementById('generation-form').style.display = 'none';
        document.getElementById('quiz-editor').style.display = 'block';
        renderQuestions();

    } catch (error) {
        alert('Error importing quiz: ' + error.message);
    }

    // Reset file input
    e.target.value = '';
}

// Show save modal
function showSaveModal() {
    document.getElementById('quiz-name').value = state.currentQuiz.subject;
    document.getElementById('save-modal').classList.add('active');
}

// Save to library
async function handleSaveToLibrary() {
    const name = document.getElementById('quiz-name').value.trim();
    if (!name) {
        alert('Please enter a quiz name');
        return;
    }

    try {
        const response = await fetch('/api/library/save', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name,
                quiz: state.currentQuiz,
                sessionId: state.sessionId
            })
        });

        if (!response.ok) throw new Error('Failed to save quiz');

        closeModals();
        alert('Quiz saved to library!');
    } catch (error) {
        alert('Error saving quiz: ' + error.message);
    }
}

// Load library
async function loadLibrary() {
    try {
        const response = await fetch('/api/library');
        if (!response.ok) throw new Error('Failed to load library');

        const quizzes = await response.json();
        renderLibrary(quizzes);
    } catch (error) {
        console.error('Error loading library:', error);
    }
}

// Render library
function renderLibrary(quizzes) {
    const container = document.getElementById('library-list');

    if (quizzes.length === 0) {
        container.innerHTML = '<div class="library-empty">No saved quizzes yet. Create and save your first quiz!</div>';
        return;
    }

    container.innerHTML = quizzes.map(quiz => `
        <div class="library-item">
            <div class="library-info">
                <h3>${escapeHtml(quiz.name)}</h3>
                <div class="library-meta">
                    Subject: ${escapeHtml(quiz.subject)} •
                    ${quiz.questions} questions •
                    Updated: ${new Date(quiz.updatedAt).toLocaleDateString()}
                </div>
            </div>
            <div class="library-actions">
                <button class="btn btn-secondary" onclick="loadQuizFromLibrary('${escapeHtml(quiz.name)}')">
                    Load
                </button>
                <button class="btn btn-danger" onclick="deleteQuizFromLibrary('${escapeHtml(quiz.name)}')">
                    Delete
                </button>
            </div>
        </div>
    `).join('');
}

// Load quiz from library
async function loadQuizFromLibrary(name) {
    try {
        const response = await fetch(`/api/library/load?name=${encodeURIComponent(name)}`);
        if (!response.ok) throw new Error('Failed to load quiz');

        const quiz = await response.json();
        state.currentQuiz = quiz;

        // Switch to generate tab and show editor
        switchTab('generate');
        document.getElementById('generation-form').style.display = 'none';
        document.getElementById('quiz-editor').style.display = 'block';
        renderQuestions();

    } catch (error) {
        alert('Error loading quiz: ' + error.message);
    }
}

// Delete quiz from library
async function deleteQuizFromLibrary(name) {
    if (!confirm(`Delete quiz "${name}"?`)) return;

    try {
        const response = await fetch(`/api/library/delete?name=${encodeURIComponent(name)}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete quiz');

        loadLibrary();
    } catch (error) {
        alert('Error deleting quiz: ' + error.message);
    }
}

// Close all modals
function closeModals() {
    document.querySelectorAll('.modal').forEach(modal => {
        modal.classList.remove('active');
    });
}

// Reset app
function resetApp() {
    state.sessionId = generateSessionId();
    state.currentQuiz = null;
    state.currentPackagePath = null;

    document.getElementById('quiz-form').reset();
    document.getElementById('generation-form').style.display = 'block';
    document.getElementById('quiz-editor').style.display = 'none';
    document.getElementById('success-message').style.display = 'none';
}

// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
