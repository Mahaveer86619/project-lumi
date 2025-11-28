package models

import "gorm.io/gorm"

type WhatsAppSession struct {
	gorm.Model

	UserID          uint   `gorm:"uniqueIndex;not null"`
	WahaSessionName string `gorm:"unique;not null"`
	Status          string `gorm:"default:'PENDING'"`
	DeviceID        string
}

type WahaProfile struct {
	ID      string `json:"id"`      // WhatsApp ID (Phone Number + @c.us)
	Name    string `json:"name"`    // User's display name
	Picture string `json:"picture"` // URL to profile picture
}

type WahaSessionInfo struct {
	Name   string      `json:"name"`
	Status string      `json:"status"` // STOPPED, STARTING, SCAN_QR_CODE, WORKING, FAILED
	Me     *WahaMeInfo `json:"me,omitempty"`
}

type WahaMeInfo struct {
	ID       string `json:"id"`       // e.g. "123456789@c.us"
	PushName string `json:"pushName"` // User's display name
}
