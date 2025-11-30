package views

import (
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/utils"
)

type UserDetailsResponse struct {
	ID             utils.MaskedId `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	AvatarUrl      string         `json:"avatar_url"`
	WhatsAppStatus string         `json:"whatsapp_status"`
}

func NewUserDetailsResponse(user models.UserProfile, waStatus string) *UserDetailsResponse {
	return &UserDetailsResponse{
		ID:        utils.Mask(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
		WhatsAppStatus: waStatus,
	}
}

type UpdateUserRequest struct {
	ID       utils.MaskedId `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
}
