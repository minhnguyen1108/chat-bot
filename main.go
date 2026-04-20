package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nguyenhoang246/go-ai-bot/internal/client"
)

var botClient *client.Client
var apiKeyLog string

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	response, err := botClient.SendMessage([]client.Message{{Role: "user", Content: req.Message}}, "You are Claude, an AI assistant made by Anthropic. Be helpful, harmless, and honest. Respond clearly and concisely. Use markdown formatting when appropriate.")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func modelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		Model string `json:"model"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	botClient.SetModel(req.Model)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	testResp, err := botClient.SendMessage([]client.Message{{Role: "user", Content: "Hi"}}, "")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiKeyLoaded": apiKeyLog,
			"model":        botClient.GetModel(),
			"status":       "error",
			"error":        err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"apiKeyLoaded": apiKeyLog,
		"model":        botClient.GetModel(),
		"status":       "ok",
		"testResponse": strings.TrimSpace(testResp),
	})
}

func main() {
	apiKey := os.Getenv("AISHOP24H_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("AISHOP24H_API_KEY or ANTHROPIC_API_KEY is required")
	}

	if len(apiKey) > 10 {
		apiKeyLog = apiKey[:10] + "..."
	} else {
		apiKeyLog = apiKey
	}
	log.Printf("API Key loaded: %s", apiKeyLog)

	botClient = client.NewClient(apiKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/chat", chatHandler)
	http.HandleFunc("/api/model", modelHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/debug", debugHandler)

	fmt.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
