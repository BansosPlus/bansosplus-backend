package handler

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	
	"github.com/BansosPlus/bansosplus-backend.git/model"
	"github.com/BansosPlus/bansosplus-backend.git/repository"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

func (h *UserHandler) GetUserHandler(c echo.Context) error {
	var user *model.User

	// Extract user token from the request attributes
	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to extract user claims",
		})
	}

	// Extract user ID from token
	userID, ok := token["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to extract user ID from claims",
		})
	}

	// Fetch user details from the database using the user ID
	user, _ = h.userRepository.GetUserByID(int(userID))
	if user == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to fetch user details",
		})
	}

	// Return user details in the response
	return c.JSON(http.StatusOK, echo.Map{
		"code": http.StatusOK,
		"status": "success",
		"data": echo.Map{
			"name": user.Name,
			"email": user.Email,
			"phone": user.Phone,
		},
	})
}