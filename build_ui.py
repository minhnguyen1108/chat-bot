#!/usr/bin/env python3
import os

html = '''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Claude</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        :root {
            --sidebar-bg: #1a1a1a;
            --sidebar-border: #2e2e2e;
            --sidebar-hover: #2a2a2a;
            --sidebar-active: #333;
            --main-bg: #0a0a0a;
            --header-bg: #0f0f0f;
            --border: #333;
            --text: #ececf1;
            --text-secondary: #8b949e;
            --text-muted: #6e7681;
            --accent: #d4a574;
            --user-bg: #1a1a2e;
            --bot-bg: #0f0f0f;
            --input-bg: #262626;
            --scrollbar: #3a3a3c;
        }
        body { font-family: Inter, -apple-system, sans-serif; background: var(--main-bg); color: var(--text); height: 100vh; overflow: hidden; }
        .app { display: flex; height: 100vh; }
        .sidebar { width: 280px; background: var(--sidebar-bg); border-right: 1px solid var(--sidebar-border); display: flex; flex-direction: column; }
        .sidebar-header { padding: 16px; border-bottom: 1px solid var(--sidebar-border); }
        .new-chat-btn { width: 100%; padding: 12px 16px; background: var(--accent); color: #0a0a0a; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 8px; }
        .new-chat-btn:hover { background: #e5b985; }
        .chat-history { flex: 1; overflow-y: auto; padding: 8px; }
        .chat-history::-webkit-scrollbar { width: 6px; }
        .chat-history::-webkit-scrollbar-thumb { background: var(--scrollbar); border-radius: 3px; }
        .history-title { font-size: 11px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; padding: 16px 8px 8px; letter-spacing: 0.5px; }
        .chat-item { padding: 12px; border-radius: 8px; cursor: pointer; margin-bottom: 4px; display: flex; align-items: center; gap: 10px; }
        .chat-item:hover { background: var(--sidebar-hover); }
        .chat-item-icon { width: 20px; height: 20px; color: var(--text-secondary); }
        .chat-item-text { flex: 1; font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .sidebar-footer { padding: 12px; border-top: 1px solid var(--sidebar-border); }
        .user-menu { display: flex; align-items: center; gap: 12px; padding: 10px; border-radius: 8px; cursor: pointer; }
        .user-menu:hover { background: var(--sidebar-hover); }
        .user-avatar { width: 32px; height: 32px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 600; }
        .user-info { flex: 1; }
        .user-name { font-size: 14px; font-weight: 500; }
        .user-status { font-size: 12px; color: var(--text-muted); }
        .main { flex: 1; display: flex; flex-direction: column; min-width: 0; }
        .header { height: 56px; background: var(--header-bg); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 20px; }
        .header-left, .header-right { display: flex; align-items: center; gap: 12px; }
        .menu-btn, .header-btn { background: none; border: none; color: var(--text); cursor: pointer; padding: 8px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
        .menu-btn:hover, .header-btn:hover { background: var(--sidebar-hover); }
        .model-selector { display: flex; align-items: center; gap: 8px; padding: 8px 16px; background: var(--input-bg); border: 1px solid var(--border); border-radius: 20px; cursor: pointer; }
        .model-selector select { background: transparent; border: none; color: var(--text); font-size: 14px; font-weight: 500; cursor: pointer; outline: none; }
        .chat-container { flex: 1; overflow-y: auto; padding: 20px; }
        .chat-container::-webkit-scrollbar { width: 8px; }
        .chat-container::-webkit-scrollbar-thumb { background: var(--scrollbar); border-radius: 4px; }
        .welcome-screen { max-width: 800px; margin: 0 auto; padding: 60px 20px; text-align: center; }
        .welcome-logo { font-size: 64px; margin-bottom: 24px; }
        .welcome-title { font-size: 32px; font-weight: 700; margin-bottom: 12px; background: linear-gradient(135deg, #d4a574 0%, #f5d09a 50%, #d4a574 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        .welcome-subtitle { font-size: 16px; color: var(--text-secondary); margin-bottom: 40px; }
        .capabilities { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; text-align: left; }
        .capability-card { background: var(--bot-bg); border: 1px solid var(--border); border-radius: 12px; padding: 20px; cursor: pointer; transition: all 0.2s; }
        .capability-card:hover { border-color: var(--accent); transform: translateY(-2px); }
        .capability-icon { font-size: 28px; margin-bottom: 12px; }
        .capability-title { font-size: 14px; font-weight: 600; margin-bottom: 8px; }
        .capability-desc { font-size: 12px; color: var(--text-secondary); line-height: 1.5; }
        .message { max-width: 800px; margin: 0 auto 24px; display: flex; gap: 16px; animation: fadeIn 0.3s ease; }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        .message-avatar { width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0; display: flex; align-items: center; justify-content: center; font-size: 16px; }
        .message-avatar.user { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
        .message-avatar.assistant { background: linear-gradient(135deg, #d4a574 0%, #f5d09a 100%); }
        .message-content { flex: 1; min-width: 0; }
        .message-role { font-size: 14px; font-weight: 600; margin-bottom: 8px; }
        .message-text { font-size: 15px; line-height: 1.7; white-space: pre-wrap; word-wrap: break-word; }
        .message-text code { background: #262626; padding: 2px 6px; border-radius: 4px; font-family: monospace; font-size: 13px; }
        .message-text pre { background: #262626; padding: 16px; border-radius: 8px; overflow-x: auto; margin: 12px 0; }
        .message-text pre code { background: none; padding: 0; }
        .message-text ul, .message-text ol { margin: 12px 0; padding-left: 24px; }
        .message-text li { margin-bottom: 6px; }
        .message-text blockquote { border-left: 3px solid var(--accent); padding-left: 16px; margin: 12px 0; color: var(--text-secondary); }
        .typing-indicator { display: flex; gap: 4px; padding: 8px 0; }
        .typing-dot { width: 8px; height: 8px; background: var(--accent); border-radius: 50%; animation: typing 1.4s infinite; }
        .typing-dot:nth-child(2) { animation-delay: 0.2s; }
        .typing-dot:nth-child(3) { animation-delay: 0.4s; }
        @keyframes typing { 0%, 60%, 100% { transform: translateY(0); opacity: 0.4; } 30% { transform: translateY(-8px); opacity: 1; } }
        .input-area { padding: 0 20px 20px; }
        .input-container { max-width: 800px; margin: 0 auto; background: var(--input-bg); border: 1px solid var(--border); border-radius: 16px; transition: all 0.2s; }
        .input-container:focus-within { border-color: var(--accent); }
        .input-actions { display: flex; padding: 8px 12px; border-bottom: 1px solid var(--border); gap: 8px; }
        .input-action-btn { background: none; border: none; color: var(--text-muted); cursor: pointer; padding: 6px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
        .input-action-btn:hover { background: var(--sidebar-hover); color: var(--text); }
        .input-textarea { width: 100%; background: transparent; border: none; color: var(--text); font-size: 15px; padding: 12px 16px; resize: none; outline: none; font-family: inherit; min-height: 52px; max-height: 200px; line-height: 1.5; }
        .input-textarea::placeholder { color: var(--text-muted); }
        .input-footer { display: flex; align-items: center; justify-content: space-between; padding: 8px 12px; }
        .input-info { font-size: 12px; color: var(--text-muted); }
        .send-btn { background: var(--accent); color: #0a0a0a; border: none; padding: 8px 16px; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 6px; }
        .send-btn:hover:not(:disabled) { background: #e5b985; }
        .send-btn:disabled { opacity: 0.5; cursor: not-allowed; }
        .tooltip { position: relative; }
        .tooltip-text { position: absolute; bottom: 100%; left: 50%; transform: translateX(-50%); background: var(--sidebar-active); color: var(--text); padding: 6px 10px; border-radius: 6px; font-size: 12px; white-space: nowrap; opacity: 0; visibility: hidden; transition: all 0.2s; margin-bottom: 6px; }
        .tooltip:hover .tooltip-text { opacity: 1; visibility: visible; }
        @media (max-width: 768px) { .sidebar { display: none; } .capabilities { grid-template-columns: 1fr; } }
    </style>
</head>
<body>
<div class="app">
    <div class="sidebar">
        <div class="sidebar-header">
            <button class="new-chat-btn" onclick="newChat()">+ New Chat</button>
        </div>
        <div class="chat-history">
            <div class="history-title">Recent Chats</div>
            <div class="chat-item active" onclick="newChat()">
                <svg class="chat-item-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
                <span class="chat-item-text">New Conversation</span>
            </div>
        </div>
        <div class="sidebar-footer">
            <div class="user-menu">
                <div class="user-avatar">U</div>
                <div class="user-info">
                    <div class="user-name">User</div>
                    <div class="user-status">Online</div>
                </div>
            </div>
        </div>
    </div>
    <div class="main">
        <div class="header">
            <div class="header-left">
                <button class="menu-btn tooltip">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>
                    <span class="tooltip-text">Menu</span>
                </button>
                <div class="model-selector">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2L2 7l10 5 10-5-10-5z"></path><path d="M2 17l10 5 10-5"></path><path d="M2 12l10 5 10-5"></path></svg>
                    <select id="modelSelect" onchange="changeModel(this.value)">
                        <option value="anthropic/claude-sonnet-4.5">Claude Sonnet 4.5</option>
                        <option value="anthropic/claude-haiku-4">Claude Haiku 4</option>
                        <option value="anthropic/claude-opus-4">Claude Opus 4</option>
                    </select>
                </div>
            </div>
            <div class="header-right">
                <button class="header-btn tooltip">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
                    <span class="tooltip-text">Settings</span>
                </button>
            </div>
        </div>
        <div class="chat-container" id="chatContainer">
            <div class="welcome-screen" id="welcomeScreen">
                <div class="welcome-logo">&#x1F916;</div>
                <h1 class="welcome-title">How can I help you today?</h1>
                <p class="welcome-subtitle">I am Claude, an AI assistant. Ask me anything.</p>
                <div class="capabilities">
                    <div class="capability-card" onclick="askCapability('Explain a complex topic')">
                        <div class="capability-icon">&#x1F4DA;</div>
                        <div class="capability-title">Explain</div>
                        <div class="capability-desc">I can explain complex topics in simple terms</div>
                    </div>
                    <div class="capability-card" onclick="askCapability('Write code for')">
                        <div class="capability-icon">&#x1F4BB;</div>
                        <div class="capability-title">Write Code</div>
                        <div class="capability-desc">I can help you write, debug, and explain code</div>
                    </div>
                    <div class="capability-card" onclick="askCapability('Brainstorm ideas for')">
                        <div class="capability-icon">&#x1F4A1;</div>
                        <div class="capability-title">Brainstorm</div>
                        <div class="capability-desc">I can help you brainstorm and plan projects</div>
                    </div>
                    <div class="capability-card" onclick="askCapability('Summarize this text')">
                        <div class="capability-icon">&#x1F4DD;</div>
                        <div class="capability-title">Summarize</div>
                        <div class="capability-desc">I can summarize long documents or articles</div>
                    </div>
                    <div class="capability-card" onclick="askCapability('Translate this to Vietnamese')">
                        <div class="capability-icon">&#x1F310;</div>
                        <div class="capability-title">Translate</div>
                        <div class="capability-desc">I can translate text between languages</div>
                    </div>
                    <div class="capability-card" onclick="askCapability('Help me analyze')">
                        <div class="capability-icon">&#x1F50D;</div>
                        <div class="capability-title">Analyze</div>
                        <div class="capability-desc">I can analyze data, documents, or problems</div>
                    </div>
                </div>
            </div>
        </div>
        <div class="input-area">
            <div class="input-container">
                <div class="input-actions">
                    <button class="input-action-btn tooltip" title="Attach file">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"></path></svg>
                    </button>
                    <button class="input-action-btn tooltip" title="Add images">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                    </button>
                </div>
                <textarea class="input-textarea" id="messageInput" placeholder="Message Claude..." rows="1" onkeydown="handleKeyDown(event)" oninput="autoResize(this)"></textarea>
                <div class="input-footer">
                    <span class="input-info">Powered by Claude AI</span>
                    <button class="send-btn" id="sendBtn" onclick="sendMessage()">Send</button>
                </div>
            </div>
        </div>
    </div>
</div>
<script>
const chatContainer = document.getElementById("chatContainer");
const messageInput = document.getElementById("messageInput");
const sendBtn = document.getElementById("sendBtn");
const welcomeScreen = document.getElementById("welcomeScreen");
let messages = [];
let isTyping = false;

function autoResize(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = Math.min(textarea.scrollHeight, 200) + "px";
}

function handleKeyDown(e) {
    if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        sendMessage();
    }
}

function newChat() {
    messages = [];
    const msgs = chatContainer.querySelectorAll(".message");
    msgs.forEach(m => m.remove());
    welcomeScreen.style.display = "block";
}

function askCapability(prompt) {
    messageInput.value = prompt;
    messageInput.focus();
}

function changeModel(model) {
    fetch("/api/model", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ model: model })
    });
}

function escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
}

function renderMarkdown(text) {
    let html = escapeHtml(text);
    html = html.replace(/```(\w+)?\\n([\\s\\S]*?)```/g, "<pre><code>$2</code></pre>");
    html = html.replace(/`([^`]+)`/g, "<code>$1</code>");
    html = html.replace(/\\*\\*([^*]+)\\*\\*/g, "<strong>$1</strong>");
    html = html.replace(/\\*([^*]+)\\*/g, "<em>$1</em>");
    html = html.replace(/^### (.+)$/gm, "<h3>$1</h3>");
    html = html.replace(/^## (.+)$/gm, "<h2>$1</h2>");
    html = html.replace(/^# (.+)$/gm, "<h1>$1</h1>");
    html = html.replace(/^\\- (.+)$/gm, "<li>$1</li>");
    html = html.replace(/\\n\\n/g, "</p><p>");
    html = "<p>" + html + "</p>";
    html = html.replace(/<p><\\/p>/g, "");
    html = html.replace(/<p>(<h[123]>)/g, "$1");
    html = html.replace(/(<\\/h[123]>)<\\/p>/g, "$1");
    return html;
}

function addMessage(role, content) {
    if (welcomeScreen.style.display !== "none") {
        welcomeScreen.style.display = "none";
    }
    const avatar = role === "user" ? "U" : "&#x1F916;";
    const avatarClass = role === "user" ? "user" : "assistant";
    const roleName = role === "user" ? "You" : "Claude";
    const html = '<div class="message"><div class="message-avatar ' + avatarClass + '">' + avatar + '</div><div class="message-content"><div class="message-role">' + roleName + '</div><div class="message-text">' + renderMarkdown(content) + '</div></div></div>';
    chatContainer.insertAdjacentHTML("beforeend", html);
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

function showTyping() {
    const html = '<div class="message" id="typingMsg"><div class="message-avatar assistant">&#x1F916;</div><div class="message-content"><div class="message-role">Claude</div><div class="message-text"><div class="typing-indicator"><div class="typing-dot"></div><div class="typing-dot"></div><div class="typing-dot"></div></div></div></div></div>';
    chatContainer.insertAdjacentHTML("beforeend", html);
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

function removeTyping() {
    const typing = document.getElementById("typingMsg");
    if (typing) typing.remove();
}

async function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || isTyping) return;
    isTyping = true;
    sendBtn.disabled = true;
    messageInput.value = "";
    messageInput.style.height = "auto";
    messages.push({ role: "user", content: text });
    addMessage("user", text);
    showTyping();
    try {
        const response = await fetch("/api/chat", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ message: text })
        });
        const data = await response.json();
        removeTyping();
        if (data.error) {
            addMessage("assistant", "Error: " + data.error);
        } else {
            messages.push({ role: "assistant", content: data.response });
            addMessage("assistant", data.response);
        }
    } catch (e) {
        removeTyping();
        addMessage("assistant", "Error: " + e.message);
    }
    isTyping = false;
    sendBtn.disabled = false;
    messageInput.focus();
}
</script>
</body>
</html>'''

go_code = '''package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nguyenhoang246/go-ai-bot/internal/client"
)

var botClient *client.Client
var apiKeyLog string

const htmlPage = `''' + html + '''

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

	response, err := botClient.SendMessage([]client.Message{{Role: "user", Content: req.Message}}, "You are Claude, an AI assistant. Be helpful, harmless, and honest. Respond clearly.")
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

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
'''

with open('D:/go/ai/main.go', 'w', encoding='utf-8') as f:
    f.write(go_code)

print("Done!")
