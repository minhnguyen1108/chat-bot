package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenhoang246/go-ai-bot/internal/client"
)

const logo = `
  ____  ____  ______          __     __
 / __ \/ __ \/_  __/__ ____  / /__  / /_
/ /_/ / /_/ / / / / __ \` + "`" + ` _ \/ / _ \/ __/
\____/\____/ /_/  \__/\___/_//_/_/\___/\__/

AI Chat Bot - Powered by Anthropic Claude
Type /help for available commands
`

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func printLogo() {
	clearScreen()
	fmt.Println(logo)
}

func printHelp() {
	fmt.Println(`
Available Commands:
  /help              - Show this help message
  /clear             - Clear chat history
  /models            - List available models
  /model <name>      - Switch to a different model
  /exit              - Exit the chat
  /tokens <number>   - Set max tokens (default: 4096)
  Ctrl+C             - Exit the chat
`)
}

type Chat struct {
	client       *client.Client
	messages     []client.Message
	systemPrompt string
	maxTokens    int
}

func NewChat(apiKey string) *Chat {
	return &Chat{
		client:       client.NewClient(apiKey),
		messages:     []client.Message{},
		systemPrompt: "You are Claude, a helpful AI assistant. Respond clearly and concisely.",
		maxTokens:    4096,
	}
}

func (c *Chat) addMessage(role, content string) {
	c.messages = append(c.messages, client.Message{
		Role:    role,
		Content: content,
	})
}

func (c *Chat) sendMessage(content string) error {
	c.addMessage("user", content)

	err := c.client.SendMessage(c.messages, c.systemPrompt)
	if err != nil {
		c.messages = c.messages[:len(c.messages)-1]
		return err
	}

	c.addMessage("assistant", "")
	return nil
}

func (c *Chat) clearHistory() {
	c.messages = []client.Message{}
}

func (c *Chat) showModels() {
	models := client.GetAvailableModels()
	currentModel := c.client.GetModel()

	fmt.Println("\nAvailable Models:")
	fmt.Println(strings.Repeat("-", 40))
	for _, model := range models {
		marker := "  "
		if model == currentModel {
			marker = "* "
		}
		fmt.Printf("%s%s\n", marker, model)
	}
	fmt.Println()
}

func (c *Chat) switchModel(modelName string) error {
	models := client.GetAvailableModels()
	for _, m := range models {
		if m == modelName {
			c.client.SetModel(modelName)
			fmt.Printf("Switched to model: %s\n\n", modelName)
			return nil
		}
	}
	return fmt.Errorf("unknown model: %s. Use /models to see available options.", modelName)
}

func (c *Chat) setMaxTokens(tokens int) {
	if tokens < 100 {
		tokens = 100
	}
	if tokens > 8192 {
		tokens = 8192
	}
	c.maxTokens = tokens
	fmt.Printf("Max tokens set to: %d\n\n", c.maxTokens)
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "go-ai-bot")
}

func loadAPIKey() string {
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		return apiKey
	}

	configPath := filepath.Join(getConfigDir(), ".env")
	if file, err := os.Open(configPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ANTHROPIC_API_KEY=") {
				apiKey := strings.TrimPrefix(line, "ANTHROPIC_API_KEY=")
				return strings.Trim(apiKey, `"' `)
			}
		}
	}
	return ""
}

func main() {
	apiKey := loadAPIKey()

	if apiKey == "" {
		fmt.Println("Error: ANTHROPIC_API_KEY not found")
		fmt.Println("Please set your API key:")
		fmt.Println("  export ANTHROPIC_API_KEY=your-api-key")
		fmt.Println("  Or create a .env file with: ANTHROPIC_API_KEY=your-api-key")
		os.Exit(1)
	}

	chat := NewChat(apiKey)

	printLogo()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nGoodbye!")
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.HasPrefix(input, "/") {
			cmd := strings.SplitN(input, " ", 2)
			command := strings.ToLower(cmd[0])
			arg := ""
			if len(cmd) > 1 {
				arg = strings.TrimSpace(cmd[1])
			}

			switch command {
			case "/help":
				printHelp()
			case "/clear":
				chat.clearHistory()
				fmt.Println("Chat history cleared!\n")
			case "/models":
				chat.showModels()
			case "/model":
				if arg == "" {
					fmt.Printf("Current model: %s\n", chat.client.GetModel())
					fmt.Println("Usage: /model <model-name> (use /models to see options)\n")
				} else if err := chat.switchModel(arg); err != nil {
					fmt.Printf("Error: %s\n\n", err.Error())
				}
			case "/exit", "/quit":
				fmt.Println("Goodbye!")
				return
			case "/tokens":
				if arg == "" {
					fmt.Printf("Current max tokens: %d\n\n", chat.maxTokens)
				} else {
					var tokens int
					if _, err := fmt.Sscanf(arg, "%d", &tokens); err == nil {
						chat.setMaxTokens(tokens)
					} else {
						fmt.Println("Invalid token count. Usage: /tokens <number>\n")
					}
				}
			default:
				fmt.Printf("Unknown command: %s. Type /help for available commands.\n\n", command)
			}
			continue
		}

		err = chat.sendMessage(input)
		if err != nil {
			fmt.Printf("\nError: %s\n\n", err.Error())
		}
	}
}
