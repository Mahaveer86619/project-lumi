package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/db"
	modelConnections "github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/services/connections"
	"google.golang.org/genai"
)

type BotService struct {
	botClient   *genai.Client
	wahaClient  connections.WahaClient
	chatService *services.ChatService
}

func NewBotService(wahaClient connections.WahaClient, chatService *services.ChatService) *BotService {
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
		wahaClient:  wahaClient,
		botClient:   client,
		chatService: chatService,
	}
}

func (b *BotService) ProcessMessage(msg modelConnections.WAMessage) {
	chatID := msg.From
	if msg.FromMe {
		chatID = msg.To
	}

	text := strings.TrimSpace(msg.Body)
	if text == "" {
		return
	}

	if msg.FromMe && msg.Source == "api" {
		return
	}

	if strings.Contains(chatID, "status") || strings.Contains(chatID, "broadcast") {
		return
	}

	b.chatService.SaveMessage(chatID, "user", text)

	if msg.From == msg.To && !b.chatService.IsChatAllowed(chatID) {
		log.Printf("Auto-registering self-chat: %s", chatID)
		b.chatService.RegisterChat(chatID, "Me (Self)", "self")
	}

	if !b.chatService.IsChatAllowed(chatID) {
		return
	}

	chat, err := b.chatService.GetRegisteredChat(chatID)
	if err != nil {
		log.Printf("Error loading registered chat: %s; error: %s", chatID, err.Error())
		return
	}

	trigger := "@lumi"
	isTrigger := strings.Contains(strings.ToLower(text), trigger)

	if !chat.IsBotActive {
		if isTrigger {
			chat.IsBotActive = true
			b.chatService.UpdateRegisteredChat(chat)
			b.replyAndSave(chatID, "Hello! I'm Lumi. Our session has started. \n\nType *bye*, *exit*, or *stop* to end the session.")

			cleanText := strings.TrimSpace(strings.ReplaceAll(text, trigger, ""))
			if cleanText != "" {
				b.generateAIResponse(chatID, cleanText)
			}
		}
		return
	}

	lowerText := strings.ToLower(text)
	if lowerText == "bye" || lowerText == "exit" || lowerText == "stop" || lowerText == "end session" {
		chat.IsBotActive = false
		db.DB.Save(&chat)
		b.replyAndSave(chatID, "Session ended. 👋 Call me again with @lumi.")
		return
	}

	b.generateAIResponse(chatID, text)
}

func (b *BotService) generateAIResponse(chatID, currentText string) {
	history, err := b.chatService.GetChatHistory(chatID, 10)
	if err != nil {
		log.Printf("Error fetching history: %v", err)
	}

	var parts []*genai.Content
	for _, h := range history {
		prefix := "User: "
		if h.Role == "model" {
			prefix = "Lumi: "
		}
		prompt := prefix + h.Content + fmt.Sprintf("current prompt: %s", currentText)

		parts = append(parts, genai.Text(prompt)...)
	}

	sysPrompt := `You are Lumi, a smart and helpful WhatsApp assistant. 
    - Format your responses using WhatsApp Markdown (e.g., *bold*, _italics_, ~strike~, ` + "`code`" + `).
    - Keep responses concise and easy to read on mobile screens.
    - If the user asks for code, wrap it in code blocks.
    - Be friendly but professional.`

	resp, err := b.botClient.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		parts,
		&genai.GenerateContentConfig{
			SystemInstruction: genai.Text(sysPrompt)[0],
		},
	)

	if err != nil {
		log.Printf("Gemini Error: %v", err)
		b.replyAndSave(chatID, "⚠️ *Error*: I'm having trouble thinking right now. Please try again.")
		return
	}

	responseText := resp.Text()
	b.replyAndSave(chatID, responseText)
}

func (b *BotService) replyAndSave(chatID, text string) {
	_, err := b.wahaClient.SendText(chatID, text)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	b.chatService.SaveMessage(chatID, "model", text)
}
