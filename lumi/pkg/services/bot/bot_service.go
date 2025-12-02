package bot

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	modelConnections "github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/services/connections"
	"google.golang.org/genai"
)

const SessionTimeout = 5 * time.Minute

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
	if msg.FromMe && msg.Source == "api" {
		return
	}

	chatID := msg.From
	if msg.FromMe {
		chatID = msg.To
	}

	if strings.Contains(chatID, "status") || strings.Contains(chatID, "broadcast") {
		return
	}

	text := strings.TrimSpace(msg.Body)
	if text == "" {
		return
	}

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

	if chat.IsBotActive && time.Since(chat.UpdatedAt) > SessionTimeout {
		log.Printf("Session timed out for %s", chatID)
		chat.IsBotActive = false
		b.chatService.UpdateRegisteredChat(chat)
		b.chatService.ClearHistory(chatID)
		b.wahaClient.SendText(chatID, "💤 LumiThread timed out due to inactivity.")
	}

	triggerKeyword := "@lumi"
	lowerText := strings.ToLower(text)
	isTrigger := strings.Contains(lowerText, triggerKeyword)
	isExit := strings.Contains(lowerText, "bye") || strings.Contains(lowerText, "exit") || strings.Contains(lowerText, "stop")

	if chat.IsBotActive && isExit {
		chat.IsBotActive = false
		b.chatService.UpdateRegisteredChat(chat)
		b.chatService.ClearHistory(chatID)
		b.wahaClient.SendText(chatID, "LumiThread ended. Data cleared. 👋")
		return
	}

	if !chat.IsBotActive {
		if isTrigger {
			chat.IsBotActive = true
			b.chatService.UpdateRegisteredChat(chat)
			b.chatService.ClearHistory(chatID)

			cleanText := strings.TrimSpace(strings.ReplaceAll(text, triggerKeyword, ""))

			if cleanText != "" {
				b.chatService.SaveMessage(chatID, "user", cleanText)
				b.generateAIResponse(chatID, cleanText)
			} else {
				b.replyAndSave(chatID, "Hey! LumiThread started. 🧠\nI'm listening. Type *bye* to exit.")
			}
		}
		return
	}

	b.chatService.UpdateRegisteredChat(chat)

	b.chatService.SaveMessage(chatID, "user", text)

	cleanPrompt := text
	if isTrigger {
		cleanPrompt = strings.TrimSpace(strings.ReplaceAll(text, triggerKeyword, ""))
	}

	b.generateAIResponse(chatID, cleanPrompt)
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

		prompt := prefix + h.Content
		parts = append(parts, genai.Text(prompt)...)
	}

	sysPrompt := config.GConfig.WahaBotSystemPrompt

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
		b.replyAndSave(chatID, "⚠️ *Error*: My brain connection timed out.")
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
