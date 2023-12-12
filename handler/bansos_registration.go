package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

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

func (h *BansosRegistrationHandler) RegisterBansosHandler(c echo.Context) error {
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

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	bansosRegistration.UserID = int(userID)

	if err := c.Bind(&bansosRegistration); err != nil || bansosRegistration.BansosID == 0 || bansosRegistration.Name == "" || bansosRegistration.Nik == "" || bansosRegistration.NoKK == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid request payload",
		})
	}

	if err := h.bansosRegistrationRepository.RegisterBansos(&bansosRegistration); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to do bansos registeration",
		})
	}

	if userRole == "admin" {
		if err := h.bansosRegistrationRepository.AcceptBansosRegis(&bansosRegistration); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to accept registration",
			})
		}
	}

	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registration added successfully",
		"data": echo.Map{
			"id":        bansosRegistration.ID,
			"user_id":   bansosRegistration.UserID,
			"bansos_id": bansosRegistration.BansosID,
			"status":    bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) AcceptBansosRegisHandler(c echo.Context) error {
	_, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	bansosRegistration, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Registration Not Found",
		})
	}

	if err := h.bansosRegistrationRepository.AcceptBansosRegis(bansosRegistration); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to accept registration",
		})
	}
	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Registration accepted successfully",
		"data": echo.Map{
			"bansos_registration_id": bansosRegistration.ID,
			"user_id":                bansosRegistration.UserID,
			"bansos_id":              bansosRegistration.BansosID,
			"status":                 bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) RejectBansosRegisHandler(c echo.Context) error {
	_, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	bansosRegistration, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Registration Not Found",
		})
	}

	if err := h.bansosRegistrationRepository.RejectBansosRegis(bansosRegistration); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to reject registration",
		})
	}
	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Registration rejected successfully",
		"data": echo.Map{
			"bansos_registration_id": bansosRegistration.ID,
			"user_id":                bansosRegistration.UserID,
			"bansos_id":              bansosRegistration.BansosID,
			"status":                 bansosRegistration.Status,
		},
	})
}

func (h *BansosRegistrationHandler) GetOnProgressBansosRegisHandler(c echo.Context) error {
	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	if userRole == "admin" {
		registrations, err := h.bansosRegistrationRepository.GetBansosRegisByStatus("ON_PROGRESS")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to retrieve bansos registrations",
			})
		}

		// Success
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Bansos registrations retrieved successfully",
			"data":    registrations,
		})
	}

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": "Unauthorized",
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByUserIDHandler(c echo.Context) error {

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

	statusStr := c.QueryParam("status")
	statusValues := strings.Split(statusStr, ",")

	bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByUserID(int(userID), statusValues)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to retrieve bansos registration",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registration retrieved successfully",
		"data":    bansosRegistrations,
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByBansosIDHandler(c echo.Context) error {

	bansosIDStr := c.QueryParam("bansos_id")

	bansosID, err := strconv.Atoi(bansosIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	token, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	userRole, ok := token["role"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Invalid token format",
		})
	}

	if userRole == "admin" {
		bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByBansosID(bansosID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    http.StatusInternalServerError,
				"status":  "error",
				"message": "Failed to retrieve bansos registrations",
			})
		}

		// Success
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "Bansos registrations retrieved successfully",
			"data":    bansosRegistrations,
		})
	}

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": "Unauthorized",
	})
}

func (h *BansosRegistrationHandler) GetBansosRegisByIDHandler(c echo.Context) error {

	bansosRegisIDStr := c.QueryParam("bansos_registration_id")

	bansosRegisID, err := strconv.Atoi(bansosRegisIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "Invalid bansos_id parameter",
		})
	}

	_, ok := c.Get("token").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"code":    http.StatusUnauthorized,
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	bansosRegistrations, err := h.bansosRegistrationRepository.GetBansosRegisByID(bansosRegisID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to retrieve bansos registrations",
		})
	}

	// Success
	return c.JSON(http.StatusOK, echo.Map{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Bansos registrations retrieved successfully",
		"data":    bansosRegistrations,
	})
}
