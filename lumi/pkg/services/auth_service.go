package services

import (
	"errors"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/utils"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"gorm.io/gorm"
)

const (
	MaxLoginAttempts = 5
	LockDuration     = 15 * time.Minute
)

type AuthService struct {
	AvatarService *AvatarService
}

func NewAuthService(avatarService *AvatarService) *AuthService {
	return &AuthService{
		AvatarService: avatarService,
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*views.AuthResponse, error) {
	var existingUser models.UserProfile
	err := db.DB.Where("username = ?", username).First(&existingUser).Error

	if err == nil {
		return nil, errors.New("username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	avatarHash := s.AvatarService.GenerateHash(email)
	user := models.UserProfile{
		Username:  username,
		Email:     email,
		Password:  hashedPwd,
		AvatarUrl: s.AvatarService.GetAvatarURL(avatarHash),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) LoginWithUsername(username, password string) (*views.AuthResponse, error) {
	var user models.UserProfile
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*views.AuthResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.UserProfile
	if err := db.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) generateAuthResponse(user models.UserProfile) (*views.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return views.NewAuthResponse(user, accessToken, refreshToken), nil
}
