package connections

import "encoding/json"

type SessionInfo struct {
	Name   string         `json:"name"`
	Status string         `json:"status"` // Enum: STOPPED, STARTING, SCAN_QR_CODE, WORKING, FAILED
	Me     *MeInfo        `json:"me,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

type MeInfo struct {
	ID       string `json:"id"`
	PushName string `json:"pushName"`
}

type SessionCreateRequest struct {
	Name   string      `json:"name"`
	Start  bool        `json:"start"`
	Config interface{} `json:"config,omitempty"`
}

type MessageTextRequest struct {
	ChatID  string `json:"chatId"`
	Text    string `json:"text"`
	Session string `json:"session"`
	ReplyTo string `json:"reply_to,omitempty"`
}

type FileWrapper struct {
	Mimetype string `json:"mimetype"`
	Data     string `json:"data,omitempty"`
	Filename string `json:"filename,omitempty"`
	Url      string `json:"url,omitempty"`
}

type ImagePayload struct {
	Caption string
	File    FileWrapper
}

type MessageImageRequest struct {
	ChatID  string      `json:"chatId"`
	Session string      `json:"session"`
	File    FileWrapper `json:"file"`
	Caption string      `json:"caption,omitempty"`
	ReplyTo string      `json:"reply_to,omitempty"`
}

type WAMessageID string

func (w *WAMessageID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// Case 1: It's a string (e.g., "true_123@c.us_ABC...")
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*w = WAMessageID(s)
		return nil
	}
	// Case 2: It's an object (e.g., {"_serialized": "...", ...})
	if data[0] == '{' {
		var obj struct {
			Serialized string `json:"_serialized"`
		}
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		*w = WAMessageID(obj.Serialized)
		return nil
	}
	return nil
}

func (w WAMessageID) String() string {
	return string(w)
}

type WAMessage struct {
	ID        WAMessageID            `json:"id"`
	Timestamp int64                  `json:"timestamp"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Body      string                 `json:"body"`
	FromMe    bool                   `json:"fromMe"`
	Source    string                 `json:"source"`
	HasMedia  bool                   `json:"hasMedia"`
	Ack       int                    `json:"ack"`
	AckName   string                 `json:"ackName"`
	Type      string                 `json:"type"` // e.g. "chat", "image", "video"
	Data      map[string]interface{} `json:"_data,omitempty"`
}

type WANumberExistResult struct {
	ChatID       string `json:"chatId,omitempty"`
	NumberExists bool   `json:"numberExists"`
}

type RequestCodeRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Method      string `json:"method,omitempty"`
}

type RequestCodeResponse struct {
	Code string `json:"code,omitempty"`
}

type ChatSummary struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Picture     string     `json:"picture"`
	LastMessage *WAMessage `json:"lastMessage"`
}

type GroupInfo struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
}

type WAHAWebhook struct {
	Event   string          `json:"event"`
	Session string          `json:"session"`
	Payload json.RawMessage `json:"payload"`
}

type SessionStatusPayload struct {
	Status string `json:"status"`
}
