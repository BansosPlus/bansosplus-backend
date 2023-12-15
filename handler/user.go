package handler

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	
	"github.com/BansosPlus/bansosplus-backend.git/model"
	"github.com/BansosPlus/bansosplus-backend.git/repository"
	"github.com/BansosPlus/bansosplus-backend.git/utility"
)

type UserHandler struct {
	userRepository repository.UserRepository
	bucketName string
	credentials string
}

func NewUserHandler(userRepository repository.UserRepository, bucketName string, credentials string) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
		bucketName: bucketName,
		credentials: credentials,
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
			"role": user.Role,
			"nik": user.Nik,
			"no_kk": user.NoKK,
			"image_url": user.ImageURL,
			"income": user.Income,
		},
	})
}

func (h *UserHandler) UpdateUserHandler(c echo.Context) error {
	var user model.User
	var exist_user *model.User

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
	exist_user, _ = h.userRepository.GetUserByID(int(userID))
	if exist_user == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to fetch user details",
		})
	}

	// Set ID in the user model
	user.ID = exist_user.ID
	user.Name = exist_user.Name
	user.Nik = exist_user.Nik
	user.NoKK = exist_user.NoKK
	user.Income = exist_user.Income
	user.FloorArea = exist_user.FloorArea
	user.WallQuality = exist_user.WallQuality
	user.NumberOfMeals = exist_user.NumberOfMeals
	user.Fuel = exist_user.Fuel
	user.Education = exist_user.Education
	user.TotalAsset = exist_user.TotalAsset
	user.Treatment = exist_user.Treatment
	user.NumberOfDependents = exist_user.NumberOfDependents
	user.Email = exist_user.Email
	user.Phone = exist_user.Phone
	user.Password = exist_user.Password
	user.Role = exist_user.Role
	user.ImageURL = exist_user.ImageURL

	// Bind payload
	if err := c.Bind(&user); err != nil || user.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
				"code":    http.StatusBadRequest,
				"status":  "error",
				"message": "Invalid request payload",
		})
	} 

	// Retrieve the file from the form data
	file, fileHeader, _ := c.Request().FormFile("file")
	if file != nil {
		defer file.Close()

		// Generate a unique filename for the uploaded file
		filename := uuid.New().String() + "_" + fileHeader.Filename

		// Upload the file to Google Cloud Storage
		imageURL, err := utility.UploadFileToGCS(file, filename, h.bucketName, h.credentials)
		if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
						"code":    http.StatusInternalServerError,
						"status":  "error",
						"message": "Failed to upload file",
				})
		}

		// Set the image URL in the User model
		user.ImageURL = imageURL
	}

	// Update user details in the database
	if err := h.userRepository.UpdateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to update user",
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
			"role": user.Role,
			"nik": user.Nik,
			"no_kk": user.NoKK,
			"image_url": user.ImageURL,
			"income": user.Income,
			"floor_area": user.FloorArea,
			"wall_quality": user.WallQuality,
			"number_of_meals": user.NumberOfMeals,
			"fuel": user.Fuel,
			"education": user.Education,
			"total_asset": user.TotalAsset,
			"treatment": user.Treatment,
			"number_of_dependents": user.NumberOfDependents,
		},
	})
}