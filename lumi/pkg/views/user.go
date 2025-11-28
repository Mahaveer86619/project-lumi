package views

import (
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/utils"
)

type UserDetailsResponse struct {
	ID           utils.MaskedId     `json:"id"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	AvatarUrl    string             `json:"avatar_url"`
	IsLocked     bool               `json:"is_locked"`
	IsSuspended  bool               `json:"is_suspended"`
	TokenVersion int                `json:"token_version"`
}

func NewUserDetailsResponse(user models.UserProfile) *UserDetailsResponse {
	return &UserDetailsResponse{
		ID:           utils.Mask(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		AvatarUrl:    user.AvatarUrl,
	}
}

type UpdateUserRequest struct {
	ID       utils.MaskedId `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
}
