package handler

import (
    "net/http"
    "strconv"
    "fmt"
    "github.com/labstack/echo"
    "github.com/google/uuid"
    "github.com/skip2/go-qrcode"
    
    "github.com/BansosPlus/bansosplus-backend.git/model"
    "github.com/BansosPlus/bansosplus-backend.git/repository"
)

type QRCodeHandler struct {
    qrCodeRepository repository.QRCodeRepository
    protocol string
    host string
    port string
}

func NewQRCodeHandler(qrCodeRepository repository.QRCodeRepository, protocol string, host string, port string) *QRCodeHandler {
    return &QRCodeHandler{
        qrCodeRepository: qrCodeRepository,
        protocol: protocol,
        host: host,
        port: port,
    }
}

func (h *QRCodeHandler) CreateQRCodeHandler(c echo.Context) error {
    var qrCode model.QRCode

    if err := c.Bind(&qrCode); err != nil || qrCode.BansosRegistrationID == 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Generate a random UUID
	uuid := uuid.New()

	// Set the UUID to your model
	qrCode.Uuid = uuid.String()

    if err := h.qrCodeRepository.AddQRCode(&qrCode); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code": http.StatusInternalServerError,
            "status": "error",
            "message": "Failed to create QR code",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "QR code added successfully",
        "data": echo.Map{
            "id": qrCode.ID,
            "uuid": qrCode.Uuid,
        },
    })
}

func (h *QRCodeHandler) ShowQRCodeByIDHandler(c echo.Context) error {
    // Get ID from the request URL
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid ID parameter",
        })
    }

    qrCode, err := h.qrCodeRepository.GetQRCodeByID(id)
    if err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{
            "code":    http.StatusNotFound,
            "status":  "error",
            "message": "QR code not found",
        })
    }

	// Generate a link based on your QRCode struct fields
	link := fmt.Sprintf("%s://%s:%s/api/qr-codes/%s", h.protocol, h.host, h.port, qrCode.Uuid)

	// Generate QR code
	qrCodeImage, err := qrcode.New(link, qrcode.Medium)
	if err != nil {
		return err
	}

    // Generate QR code image bytes with the specified size
    qrCodeBytes, err := qrCodeImage.PNG(256)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "Failed to generate QR code image",
		})
	}

	// Serve the QR code as an image
	c.Response().Header().Set(echo.HeaderContentType, "image/png")
	return c.Blob(http.StatusOK, "image/png", qrCodeBytes)
}

func (h *QRCodeHandler) GetQRCodeByUUIDHandler(c echo.Context) error {
    // Get UUID from the request URL
    uuid := c.Param("uuid")

    qrCode, err := h.qrCodeRepository.GetQRCodeByUUID(uuid)
    if err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{
            "code":    http.StatusNotFound,
            "status":  "error",
            "message": "QR code not found",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "QR code retrieved successfully",
        "data":    qrCode,
    })
}