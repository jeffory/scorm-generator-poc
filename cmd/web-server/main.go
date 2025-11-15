package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jeffory/scorm-generator-poc/internal/server"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to run the web server on")
	apiKey := flag.String("api-key", os.Getenv("OPENROUTER_API_KEY"), "OpenRouter API key")
	model := flag.String("model", "anthropic/claude-3.5-sonnet", "OpenRouter model to use")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("OpenRouter API key is required. Set OPENROUTER_API_KEY environment variable or use -api-key flag")
	}

	// Create and start server
	srv := server.New(*apiKey, *model)

	log.Printf("Starting SCORM Generator Web Server on port %s", *port)
	log.Printf("Open http://localhost:%s in your browser", *port)

	if err := http.ListenAndServe(":"+*port, srv); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
