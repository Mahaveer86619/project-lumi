package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Mahaveer86619/lumi/pkg/config"
	modelConnections "github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services/connections"
	"google.golang.org/genai"
)

type BotService struct {
	botClient  *genai.Client
	wahaClient connections.WahaClient
}

func NewBotService(wahaClient connections.WahaClient) *BotService {
	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey: config.GConfig.GeminiAPIKey,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &BotService{
		wahaClient: wahaClient,
		botClient:  client,
	}
}

func (b *BotService) ProcessMessage(msg modelConnections.WAMessage) {
	chatID := msg.From
	if msg.FromMe {
		chatID = msg.To
	}
	
	text := strings.TrimSpace(msg.Body)

	// --- Loop Protection ---
	// If the message is from me (Note to self) AND starts with our AI prefix "> ",
	// it's likely a bot response. Ignore it to prevent infinite loops.
	if msg.FromMe && strings.HasPrefix(text, "> ") {
		return
	}

	// --- Command Handling ---
	if strings.HasPrefix(text, "/") {
		b.handleCommand(chatID, text)
		return
	}

	// --- AI Response ---
	// Respond to every other message with AI
	b.handleAIResponse(chatID, text)
}

func (b *BotService) handleCommand(chatID, text string) {
	parts := strings.SplitN(text, " ", 2)
	command := parts[0]

	switch command {
	case "/ping":
		b.wahaClient.SendText(chatID, "Pong! 🏓")
	case "/help":
		b.wahaClient.SendText(chatID, "Available commands: /ping, /help, /explain <text>")
	case "/explain":
		if len(parts) < 2 {
			b.wahaClient.SendText(chatID, "Usage: /explain <text>")
			return
		}
		b.handleAIResponse(chatID, fmt.Sprintf("Explain: %s", parts[1]))
	default:
		b.wahaClient.SendText(chatID, "Unknown command.")
	}
}

func (b *BotService) handleAIResponse(chatID, text string) {
	if text == "" {
		return
	}

	responseText, err := b.simpleTextResponse(text)
	if err != nil {
		log.Printf("Gemini Error: %v", err)
		b.wahaClient.SendText(chatID, "Error generating AI response.")
		return
	}

	_, err = b.wahaClient.SendText(chatID, responseText)
	if err != nil {
		log.Printf("Failed to send AI response: %v", err)
	}
}

func (b *BotService) simpleTextResponse(prompt string) (string, error) {

	log.Printf("Prompt: %s", prompt)

	resp, err := b.botClient.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			SystemInstruction: genai.Text("You are a helpful WhatsApp assistant named lumi bot for mahaveer and others will invoke you for responses. Keep answers concise and friendly.")[0],
		},
	)
	if err != nil {
		log.Printf("Gemini Error: %v", err)
		return "", err
	}

	log.Printf("Response: %s", resp.Text())


	return fmt.Sprintf("> %s", resp.Text()), nil
}
