package handler

import (
    "net/http"
    "github.com/labstack/echo"
    "github.com/dgrijalva/jwt-go"
    
    "github.com/BansosPlus/bansosplus-backend.git/model"
    "github.com/BansosPlus/bansosplus-backend.git/repository"
)

type BansosRegistrationHandler struct {
    bansosRegistrationRepository repository.BansosRegistrationRepository
}

func NewBansosRegistrationHandler(bansosRegistrationRepository repository.BansosRegistrationRepository) *BansosRegistrationHandler {
    return &BansosRegistrationHandler{
        bansosRegistrationRepository: bansosRegistrationRepository,
    }
}

func (h *BansosRegistrationHandler) RegisterBansos(c echo.Context) error {
    var bansosRegistration model.BansosRegistration

    token, ok := c.Get("token").(jwt.MapClaims)
    if !ok {
        return c.JSON(http.StatusUnauthorized, echo.Map{
            "code":    http.StatusUnauthorized,
            "status":  "error",
            "message": "Unauthorized",
        })
    }

    userID, ok := token["id"].(float64)
    if !ok {
        return c.JSON(http.StatusUnauthorized, echo.Map{
            "code":    http.StatusUnauthorized,
            "status":  "error",
            "message": "Invalid token format",
        })
    }

    // Set userID in the feedback model
    bansosRegistration.UserID = int(userID)

    // Bind payload
    if err := c.Bind(&bansosRegistration); err != nil || bansosRegistration.BansosID == 0 || bansosRegistration.Name == "" || bansosRegistration.Nik == "" || bansosRegistration.NoKK == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    if err := h.feedbackRepository.RegisterBansos(&bansosRegistration); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code": http.StatusInternalServerError,
            "status": "error",
            "message": "Failed to register user",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Feedback added successfully",
        "data": echo.Map{
            "id": bansosRegistration.ID,
            "user_id": bansosRegistration.UserID,
            "bansos_id": bansosRegistration.BansosID,
        },
    })
}