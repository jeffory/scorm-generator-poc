# SCORM Generator POC

A Go application that generates SCORM 1.2 compliant quiz packages using OpenRouter's LLM API.

## Features

- 🤖 **AI-Powered**: Uses OpenRouter to generate educational multiple-choice questions on any subject
- 📦 **SCORM 1.2 Compliant**: Creates standard SCORM packages compatible with most Learning Management Systems (LMS)
- 🎨 **Modern UI**: Includes a responsive, attractive quiz interface
- 📊 **Progress Tracking**: Tracks quiz completion and scores via SCORM API
- ✅ **Instant Feedback**: Shows correct/incorrect answers after submission

## Prerequisites

- Go 1.19 or later
- OpenRouter API key ([Get one here](https://openrouter.ai/))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/jeffory/scorm-generator-poc.git
cd scorm-generator-poc
```

2. Build the application:
```bash
go build -o scorm-generator ./cmd/scorm-generator
```

## Usage

### Basic Usage

```bash
export OPENROUTER_API_KEY="your-api-key-here"
./scorm-generator -subject "World History"
```

### Advanced Options

```bash
./scorm-generator \
  -subject "Quantum Physics" \
  -questions 10 \
  -output ./my-courses \
  -model "anthropic/claude-3.5-sonnet" \
  -api-key "your-api-key"
```

### Command-Line Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-subject` | Subject for the quiz | - | Yes |
| `-questions` | Number of questions to generate | 5 | No |
| `-output` | Output directory for SCORM package | `./output` | No |
| `-model` | OpenRouter model to use | `anthropic/claude-3.5-sonnet` | No |
| `-api-key` | OpenRouter API key | From env: `OPENROUTER_API_KEY` | Yes |

### Environment Variables

- `OPENROUTER_API_KEY`: Your OpenRouter API key (alternative to `-api-key` flag)

## Output

The application generates a SCORM 1.2 package as a ZIP file containing:

- `imsmanifest.xml` - SCORM manifest file
- `index.html` - Interactive quiz interface
- `scormapi.js` - SCORM API wrapper for LMS communication
- `styles.css` - Modern, responsive styling

The ZIP file can be directly uploaded to any SCORM 1.2 compatible LMS (Moodle, Canvas, Blackboard, etc.)

## Example Output

```
Initializing OpenRouter client...
Generating 5 questions about 'World History' using model anthropic/claude-3.5-sonnet...
Successfully generated 5 questions
Creating SCORM package...

✓ SCORM package created successfully!
  Location: ./output/World_History.zip
  Subject: World History
  Questions: 5

You can now upload this SCORM package to your LMS (Learning Management System)
```

## SCORM Package Features

### For Learners
- Clean, modern interface
- Navigate between questions
- Review answers before submission
- See detailed results with correct answers
- Progress is tracked in the LMS

### For LMS Administrators
- SCORM 1.2 compliant
- Reports completion status
- Tracks scores (0-100%)
- Pass/fail threshold at 70%
- Supports suspend/resume

## Available OpenRouter Models

Some popular models you can use:

- `anthropic/claude-3.5-sonnet` (default) - Best quality
- `anthropic/claude-3-haiku` - Faster, more economical
- `openai/gpt-4-turbo` - OpenAI's latest
- `openai/gpt-3.5-turbo` - Fast and economical
- `meta-llama/llama-3-70b-instruct` - Open source option

See [OpenRouter's model list](https://openrouter.ai/models) for all available models.

## Project Structure

```
.
├── cmd/
│   └── scorm-generator/    # Main application entry point
│       └── main.go
├── internal/
│   ├── openrouter/         # OpenRouter API client
│   │   └── client.go
│   ├── quiz/               # Quiz generation logic
│   │   └── generator.go
│   └── scorm/              # SCORM package generator
│       ├── generator.go    # Main SCORM generator
│       ├── manifest.go     # imsmanifest.xml generator
│       ├── templates.go    # HTML template generator
│       ├── scormapi.go     # SCORM API JavaScript
│       └── styles.go       # CSS styles
└── pkg/
    └── models/             # Data models
        └── quiz.go         # Quiz and Question models
```

## Development

### Building

```bash
go build -o scorm-generator ./cmd/scorm-generator
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

## How It Works

1. **User Input**: You provide a subject and number of questions
2. **LLM Query**: The app sends a prompt to OpenRouter requesting quiz questions
3. **JSON Parsing**: The LLM response is parsed into structured question data
4. **SCORM Generation**:
   - Creates `imsmanifest.xml` with SCORM metadata
   - Generates interactive HTML quiz interface
   - Includes SCORM API wrapper for LMS communication
   - Adds responsive CSS styling
5. **Package Creation**: All files are zipped into a SCORM package

## Troubleshooting

### "Error: OpenRouter API key is required"
Set your API key via environment variable:
```bash
export OPENROUTER_API_KEY="your-key-here"
```

Or pass it directly:
```bash
./scorm-generator -api-key "your-key-here" -subject "Math"
```

### "Failed to parse quiz response"
The LLM might have returned unexpected output. Try:
- Using a different model with `-model`
- Reducing the number of questions
- Simplifying the subject name

### SCORM Package Not Working in LMS
- Ensure your LMS supports SCORM 1.2
- Check that the ZIP file structure is intact
- Review LMS logs for specific errors

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - feel free to use this project for any purpose.

## Acknowledgments

- Built with [OpenRouter](https://openrouter.ai/) for LLM access
- SCORM 1.2 specification by ADL (Advanced Distributed Learning)
- Go standard library for robust functionality

## Future Enhancements

- [ ] SCORM 2004 support
- [ ] Support for different question types (true/false, fill-in-the-blank)
- [ ] Customizable quiz themes
- [ ] Include explanations for correct answers
- [ ] Support for images in questions
- [ ] Batch generation from CSV
- [ ] Web interface

## Contact

For questions or feedback, please open an issue on GitHub.
