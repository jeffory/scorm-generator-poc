# SCORM Generator - Web Application

A modern web application that generates SCORM 1.2 compliant quiz packages using AI. Features an intuitive web interface with quiz editing, preview, and library management.

## ✨ Features

### Core Functionality
- 🌐 **Modern Web Interface**: Beautiful, responsive single-page application
- 🤖 **AI-Powered Generation**: Uses OpenRouter to create educational multiple-choice questions on any subject
- ✏️ **Inline Editing**: Edit questions and answers directly in the browser
- 🔄 **Smart Regeneration**: Regenerate individual questions with one click
- 📦 **SCORM 1.2 Compliant**: Creates standard packages compatible with most LMS platforms
- 👁️ **SCORM Preview**: Preview your quiz before deployment
- 🎯 **Customizable Settings**: Configure passing scores and quiz parameters

### Library & Management
- 📚 **Quiz Library**: Save and manage multiple quiz configurations
- 💾 **Import/Export**: Import and export quizzes as JSON files
- 📊 **Quiz History**: Track all your created quizzes with metadata
- 🔍 **Quick Search**: Find and load saved quizzes instantly

### User Experience
- 🎨 **Beautiful UI**: Modern gradient design with smooth animations
- 📱 **Fully Responsive**: Works perfectly on desktop, tablet, and mobile
- ⚡ **Real-time Updates**: Instant feedback on all operations
- 🎓 **Educational Focus**: Clean, distraction-free interface

## Prerequisites

- Go 1.19 or later
- OpenRouter API key ([Get one here](https://openrouter.ai/))
- Modern web browser (Chrome, Firefox, Safari, Edge)

## Quick Start

### 1. Clone and Build

```bash
git clone https://github.com/jeffory/scorm-generator-poc.git
cd scorm-generator-poc
go build -o bin/web-server ./cmd/web-server
```

### 2. Set Your API Key

```bash
export OPENROUTER_API_KEY="your-api-key-here"
```

### 3. Start the Server

```bash
./bin/web-server
```

The server will start on `http://localhost:8080` by default.

### 4. Open Your Browser

Navigate to `http://localhost:8080` and start creating SCORM quizzes!

## Configuration Options

### Command-Line Flags

```bash
./bin/web-server \
  -port 8080 \
  -api-key "your-api-key" \
  -model "anthropic/claude-3.5-sonnet"
```

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-port` | Port to run the web server on | `8080` | No |
| `-api-key` | OpenRouter API key | From env: `OPENROUTER_API_KEY` | Yes |
| `-model` | OpenRouter model to use | `anthropic/claude-3.5-sonnet` | No |

### Environment Variables

- `OPENROUTER_API_KEY`: Your OpenRouter API key (recommended method)

## Using the Web Interface

### Generating a Quiz

1. **Enter Quiz Details**
   - Subject/Topic (e.g., "World History", "Python Programming")
   - Number of questions (1-50)
   - Passing score percentage (1-100%)

2. **Click "Generate Questions"**
   - AI will create questions based on your subject
   - Takes 10-30 seconds depending on question count

3. **Review & Edit**
   - See all generated questions
   - Click ✏️ to edit any question
   - Click 🔄 to regenerate individual questions
   - All edits are saved automatically

### Managing Questions

#### Edit a Question
- Click the ✏️ edit icon
- Modify question text, options, or correct answer
- Click "Save" to apply changes

#### Regenerate a Question
- Click the 🔄 regenerate icon
- AI will create a new question to replace it
- Review and edit if needed

### Creating SCORM Package

1. Review all questions
2. Click "📦 Create SCORM Package"
3. Package is generated in seconds
4. Options:
   - **⬇️ Download**: Download the ZIP file
   - **👁️ Preview**: Preview the quiz (extract and open index.html)
   - **➕ New Quiz**: Start a new quiz

### Quiz Library Features

#### Save a Quiz
- Click "💾 Save to Library" in the editor
- Enter a memorable name
- Quiz is saved with all metadata

#### Load a Quiz
- Switch to "Quiz Library" tab
- Browse saved quizzes
- Click "Load" to continue editing
- Click "Delete" to remove from library

#### Export/Import JSON
- **Export**: Download quiz as JSON file for backup or sharing
- **Import**: Load quiz from JSON file (drag & drop or file picker)

## Web UI Features in Detail

### Generate Tab
- **Generation Form**: Simple 3-field form to create new quizzes
- **Question Editor**: Rich editor with inline editing capabilities
- **Action Buttons**: Save, export, import, and generate SCORM
- **Success Screen**: Download and preview options

### Library Tab
- **Quiz Cards**: Visual cards showing quiz metadata
- **Quick Actions**: Load or delete with one click
- **Metadata Display**: Subject, question count, and last updated date

### Modals
- **Preview Modal**: Full-screen SCORM preview
- **Save Modal**: Name your quiz before saving

## Project Structure

```
.
├── cmd/
│   ├── scorm-generator/    # Original CLI application (legacy)
│   └── web-server/         # Web server application (main)
├── internal/
│   ├── openrouter/         # OpenRouter API client
│   ├── quiz/               # Quiz generation logic
│   ├── scorm/              # SCORM package generator
│   ├── server/             # Web server and HTTP handlers
│   │   ├── server.go       # Server setup and routing
│   │   ├── handlers.go     # API endpoint handlers
│   │   └── web/            # Static web files
│   │       ├── index.html  # Main web UI
│   │       ├── styles.css  # UI styling
│   │       └── app.js      # Frontend logic
│   ├── storage/            # Quiz library storage
│   └── tui/                # Terminal UI (for CLI version)
├── pkg/
│   └── models/             # Data models (Quiz, Question)
├── data/                   # Saved quizzes (created on first run)
└── output/                 # Generated SCORM packages
```

## API Endpoints

The web server exposes the following REST API endpoints:

### Quiz Generation
- `POST /api/generate` - Generate quiz questions
- `POST /api/regenerate-question` - Regenerate a single question
- `POST /api/update-question` - Update a question with edits
- `POST /api/generate-scorm` - Create SCORM package

### Library Management
- `GET /api/library` - List all saved quizzes
- `POST /api/library/save` - Save quiz to library
- `GET /api/library/load?name=xyz` - Load quiz from library
- `DELETE /api/library/delete?name=xyz` - Delete quiz from library

### Import/Export
- `GET /api/export?sessionId=xyz` - Export quiz as JSON
- `POST /api/import?sessionId=xyz` - Import quiz from JSON

### Package Serving
- `GET /packages/<filename>.zip` - Download SCORM package

## SCORM Package Features

### For Learners
- Clean, modern interface with gradient design
- Navigate between questions freely
- Review answers before submission
- Detailed results with correct/incorrect feedback
- Progress tracked in the LMS

### For LMS Administrators
- SCORM 1.2 compliant (works with Moodle, Canvas, Blackboard, etc.)
- Reports completion status
- Tracks scores (0-100%)
- Configurable pass/fail threshold
- Supports suspend/resume functionality

## Available OpenRouter Models

Popular models you can use with the `-model` flag:

- `anthropic/claude-3.5-sonnet` (default) - Best quality, most accurate
- `anthropic/claude-3-haiku` - Faster, more economical
- `openai/gpt-4-turbo` - OpenAI's latest high-quality model
- `openai/gpt-3.5-turbo` - Fast and economical
- `meta-llama/llama-3-70b-instruct` - Open source option

See [OpenRouter's model list](https://openrouter.ai/models) for all available models.

## Development

### Building from Source

```bash
# Build web server
go build -o bin/web-server ./cmd/web-server

# Build CLI (legacy)
go build -o bin/scorm-generator ./cmd/scorm-generator
```

### Running Tests

```bash
go test ./...
```

### Adding Dependencies

```bash
go get <package>
go mod tidy
```

### Development Mode

For development with auto-reload, use `air` or similar:

```bash
go install github.com/cosmtrek/air@latest
air
```

## Deployment

### Internal Hosting

Since no authentication is required, this is designed for internal hosting:

```bash
# Run on custom port
./bin/web-server -port 3000

# Run with production settings
OPENROUTER_API_KEY="your-key" ./bin/web-server -port 80
```

### Docker Deployment (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o web-server ./cmd/web-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web-server .
RUN mkdir -p data output
EXPOSE 8080
CMD ["./web-server"]
```

### Systemd Service

Create `/etc/systemd/system/scorm-generator.service`:

```ini
[Unit]
Description=SCORM Generator Web Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/scorm-generator
Environment="OPENROUTER_API_KEY=your-key-here"
ExecStart=/opt/scorm-generator/bin/web-server -port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl enable scorm-generator
sudo systemctl start scorm-generator
```

## Troubleshooting

### Server Won't Start

**Error: "OpenRouter API key is required"**
```bash
export OPENROUTER_API_KEY="your-key-here"
./bin/web-server
```

**Error: "Port already in use"**
```bash
./bin/web-server -port 8081
```

### Quiz Generation Issues

**"Failed to generate quiz"**
- Check your API key is valid
- Ensure you have OpenRouter credits
- Try a different model
- Reduce the number of questions

**Questions are low quality**
- Try `anthropic/claude-3.5-sonnet` for best results
- Make your subject more specific
- Generate fewer questions at once

### SCORM Package Issues

**Package won't upload to LMS**
- Ensure your LMS supports SCORM 1.2
- Check the ZIP file isn't corrupted
- Review LMS logs for specific errors

**Quiz doesn't track progress**
- Verify LMS SCORM support is enabled
- Check SCORM API is accessible
- Test with a different LMS

## Security Notes

⚠️ **This application has NO authentication** - it's designed for internal hosting only.

For production deployment:
- Run behind a VPN or firewall
- Use reverse proxy with authentication (nginx + basic auth, OAuth, etc.)
- Limit network access to trusted users
- Monitor API usage and costs

## CLI Version (Legacy)

The original CLI version is still available:

```bash
go build -o bin/scorm-generator ./cmd/scorm-generator
export OPENROUTER_API_KEY="your-key"
./bin/scorm-generator -subject "Math" -questions 5
```

See git history for CLI documentation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Ideas for Contributions
- Add authentication system
- Implement user accounts and permissions
- Add more question types (true/false, short answer)
- Improve SCORM preview (in-browser unzip and display)
- Add quiz templates and categories
- Implement analytics dashboard
- Support for images in questions
- Multi-language support

## License

MIT License - feel free to use this project for any purpose.

## Acknowledgments

- Built with [OpenRouter](https://openrouter.ai/) for LLM access
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the CLI TUI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for TUI styling
- SCORM 1.2 specification by ADL (Advanced Distributed Learning)
- Go standard library for robust functionality

## Roadmap

- [x] Web interface
- [x] Quiz library/history
- [x] Import/export JSON
- [x] Inline editing
- [x] Question regeneration
- [x] SCORM preview
- [ ] User authentication
- [ ] SCORM 2004 support
- [ ] Additional question types
- [ ] Image support in questions
- [ ] Quiz templates
- [ ] Analytics dashboard
- [ ] Batch import from CSV
- [ ] LMS integration testing

## Contact

For questions or feedback, please open an issue on GitHub.

---

**Built with ❤️ for educators and instructional designers**
