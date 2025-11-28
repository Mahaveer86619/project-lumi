package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/utils"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(group *echo.Group, us *services.UserService) *UserHandler {
	handler := &UserHandler{
		userService: us,
	}

	group.GET("/me", handler.GetUserDetails)
	group.PUT("/user", handler.UpdateUser)
	group.DELETE("/user", handler.DeleteUser)

	return handler
}

func (h *UserHandler) GetUserDetails(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage("Unauthorized")
		return failure.JSON(c)
	}

	resp, err := h.userService.GetUserDetails(userID)
	if err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusInternalServerError)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("User details fetched successfully")
	success.SetData(resp)

	return success.JSON(c)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage("Unauthorized")
		return failure.JSON(c)
	}

	var userRequest views.UpdateUserRequest
	if err := c.Bind(&userRequest); err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	if userID != utils.Unmask(userRequest.ID) {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusForbidden)
		failure.SetMessage("You are not authorized to update this profile")
		return failure.JSON(c)
	}

	resp, err := h.userService.UpdateUserDetails(userRequest)
	if err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusInternalServerError)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("User details updated successfully")
	success.SetData(resp)

	return success.JSON(c)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage("Unauthorized")
		return failure.JSON(c)
	}

	type deleteRequest struct {
		ID utils.MaskedId `json:"id"`
	}
	var deleteReq deleteRequest
	if err := c.Bind(&deleteReq); err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	if userID != utils.Unmask(deleteReq.ID) {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusForbidden)
		failure.SetMessage("You are not authorized to update this profile")
		return failure.JSON(c)
	}

	err := h.userService.DeleteUser(deleteReq.ID.Unmask())
	if err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusInternalServerError)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("User details deleted successfully")

	return success.JSON(c)
}
