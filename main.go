package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nguyenhoang246/go-ai-bot/internal/client"
)

var botClient *client.Client

const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Claude AI Bot</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #1a1a2e; color: #eee; min-height: 100vh; display: flex; flex-direction: column; }
        .header { background: #16213e; padding: 15px 20px; border-bottom: 1px solid #0f3460; display: flex; justify-content: space-between; align-items: center; }
        .header h1 { color: #00d9ff; font-size: 1.4rem; }
        .model-select { background: #0f3460; color: #fff; border: 1px solid #00d9ff; padding: 8px 12px; border-radius: 8px; cursor: pointer; }
        .chat-container { flex: 1; max-width: 900px; margin: 0 auto; width: 100%; padding: 20px; overflow-y: auto; }
        .message { margin: 15px 0; padding: 15px 20px; border-radius: 15px; max-width: 80%; line-height: 1.6; white-space: pre-wrap; word-wrap: break-word; }
        .user { background: #0f3460; margin-left: auto; border-bottom-right-radius: 5px; }
        .bot { background: #16213e; border: 1px solid #0f3460; margin-right: auto; border-bottom-left-radius: 5px; }
        .input-area { background: #16213e; padding: 20px; border-top: 1px solid #0f3460; }
        .input-container { max-width: 900px; margin: 0 auto; display: flex; gap: 10px; }
        textarea { flex: 1; background: #0f3460; border: 1px solid #333; color: #fff; padding: 15px; border-radius: 10px; resize: none; font-size: 1rem; min-height: 50px; max-height: 200px; font-family: inherit; }
        textarea:focus { outline: none; border-color: #00d9ff; }
        button { background: #00d9ff; color: #1a1a2e; border: none; padding: 15px 30px; border-radius: 10px; cursor: pointer; font-weight: bold; font-size: 1rem; }
        button:hover { background: #00b8d9; }
        button:disabled { background: #555; cursor: not-allowed; }
        .typing { color: #00d9ff; font-style: italic; }
        .clear-btn { background: #ff6b6b; padding: 8px 15px; border-radius: 8px; font-size: 0.9rem; }
        .clear-btn:hover { background: #ee5a5a; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Claude AI Bot</h1>
        <div>
            <button class="clear-btn" onclick="clearChat()">Clear</button>
            <select class="model-select" onchange="changeModel(this.value)">
                <option value="claude-3-5-haiku-20241007">Haiku (Fast)</option>
                <option value="claude-3-5-sonnet-20241022">Sonnet (Balanced)</option>
                <option value="claude-3-opus-20240229">Opus (Powerful)</option>
            </select>
        </div>
    </div>
    <div class="chat-container" id="chat"></div>
    <div class="input-area">
        <div class="input-container">
            <textarea id="input" placeholder="Type your message..." onkeydown="handleKey(event)" rows="1"></textarea>
            <button onclick="sendMessage()" id="sendBtn">Send</button>
        </div>
    </div>
    <script>
        const chat = document.getElementById('chat');
        const input = document.getElementById('input');
        const sendBtn = document.getElementById('sendBtn');

        async function sendMessage() {
            const text = input.value.trim();
            if (!text) return;

            addMessage('user', text);
            input.value = '';
            sendBtn.disabled = true;

            const botMsg = addMessage('bot', 'Claude is thinking...');

            try {
                const response = await fetch('/api/chat', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ message: text })
                });
                const data = await response.json();

                if (data.error) {
                    botMsg.innerHTML = '<strong>Claude:</strong>\n<span style="color:#ff6b6b">Error: ' + data.error + '</span>';
                } else {
                    botMsg.innerHTML = '<strong>Claude:</strong>\n' + escapeHtml(data.response);
                }
            } catch (e) {
                botMsg.innerHTML = '<strong>Claude:</strong>\n<span style="color:#ff6b6b">Error: ' + e.message + '</span>';
            }

            sendBtn.disabled = false;
            input.focus();
        }

        function escapeHtml(text) {
            return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\n/g, '<br>');
        }

        function addMessage(role, text) {
            const div = document.createElement('div');
            div.className = 'message ' + role;
            div.innerHTML = text;
            chat.appendChild(div);
            chat.scrollTop = chat.scrollHeight;
            return div;
        }

        function handleKey(e) {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        }

        async function changeModel(model) {
            await fetch('/api/model', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ model: model })
            });
            addMessage('bot', '<strong>Claude:</strong>\nModel switched!');
        }

        async function clearChat() {
            chat.innerHTML = '';
        }
    </script>
</body>
</html>`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlPage))
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	response, err := botClient.SendMessage([]client.Message{{Role: "user", Content: req.Message}}, "You are Claude, a helpful AI assistant. Be concise and clear.")
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
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is required")
	}

	botClient = client.NewClient(apiKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/chat", chatHandler)
	http.HandleFunc("/api/model", modelHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
