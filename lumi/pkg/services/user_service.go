package services

import (
	"errors"

	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/views"
)

type UserService struct {
	ws *WahaService
}

func NewUserService(ws *WahaService) *UserService {
	return &UserService{
		ws: ws,
	}
}

func (us UserService) GetUserDetails(id uint) (*views.UserDetailsResponse, error) {
	var user models.UserProfile
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	response := views.NewUserDetailsResponse(user)
	return response, nil
}

func (us UserService) UpdateUserDetails(req views.UpdateUserRequest) (*views.UserDetailsResponse, error) {
	var user models.UserProfile
	if err := db.DB.Where("id = ?", req.ID).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Username = req.Username
	user.Email = req.Email

	if err := db.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	response := views.NewUserDetailsResponse(user)
	return response, nil
}

func (us UserService) DeleteUser(id uint) error {
	var user models.UserProfile
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.New("invalid credentials")
	}

	if err := db.DB.Model(&user).Unscoped().Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
