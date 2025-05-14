package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	h "scripts/components/helper"
	"strconv"
	"time"
)

const (
	telegramApiUrl   = "https://api.telegram.org/bot"
	apiTimeout       = 10 * time.Second
	maxMessageLength = 4096
)

type Bot struct {
	client *http.Client
	chatId uint64
	uri    string
}

func NewBot(token string, chatId string) *Bot {
	tgChatId := validate(token, chatId)
	return &Bot{
		client: &http.Client{Timeout: apiTimeout},
		chatId: tgChatId,
		uri:    fmt.Sprintf("%s%s", telegramApiUrl, token),
	}
}

func (b *Bot) Send(message string) {
	wd, _ := os.Getwd()
	message = wd + "\n" + message
	if len(message) > maxMessageLength {
		h.Fatal("message too long")
	}

	requestBody := map[string]any{
		"chat_id": b.chatId,
		"text":    message,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		h.Fatal("Failed to marshal payload:", err)
	}

	resp, err := b.client.Post(fmt.Sprintf("%s/sendMessage", b.uri), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		h.Fatal("HTTP request failed:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		h.Fatal("API error. status:", resp.Status, "response:", string(body), "request:", string(jsonBody))
	}
}

func validate(token string, chatId string) uint64 {
	if token == "" || chatId == "" {
		h.Fatal("TELEGRAM_TOKEN OR TELEGRAM_CHAT_ID not found in OS ENV")
	}

	if !regexp.MustCompile(`^\d+:[A-Za-z0-9_-]+$`).MatchString(token) {
		h.Fatal("Incorrect TELEGRAM_TOKEN:", token)
	}

	tgChatId, err := strconv.ParseUint(chatId, 10, 64)
	if err != nil {
		h.Fatal("incorrect TELEGRAM_CHAT_ID:", chatId)
	}
	return tgChatId
}
