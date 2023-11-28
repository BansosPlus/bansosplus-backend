package handler

import (
    "net/http"
    "strconv"
    "time"
    "github.com/labstack/echo"
    "github.com/google/uuid"
    
    "github.com/BansosPlus/bansosplus-backend.git/model"
    "github.com/BansosPlus/bansosplus-backend.git/repository"
    "github.com/BansosPlus/bansosplus-backend.git/utility"
)

type BansosHandler struct {
    bansosRepository repository.BansosRepository
    bucketName string
    credentials string
}

func NewBansosHandler(bansosRepository repository.BansosRepository, bucketName string, credentials string) *BansosHandler {
    return &BansosHandler{
        bansosRepository: bansosRepository,
        bucketName: bucketName,
        credentials: credentials,
    }
}

func (h *BansosHandler) AddBansosHandler(c echo.Context) error {
    var bansos model.Bansos

    // Bind payload
    if err := c.Bind(&bansos); err != nil || bansos.Name == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Retrieve the file from the form data
    file, fileHeader, err := c.Request().FormFile("file")
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Failed to get file from request",
        })
    }
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

    // Set ExpiryDate to 3 months from now
    bansos.ExpiryDate = time.Now().AddDate(0, 3, 0)

    // Set the image URL in the Bansos model
    bansos.ImageURL = imageURL
    
    if err := h.bansosRepository.AddBansos(&bansos); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code": http.StatusInternalServerError,
            "status": "error",
            "message": "Failed to create bansos",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Bansos added successfully",
        "data": echo.Map{
            "id": bansos.ID,
            "name": bansos.Name,
            "type": bansos.Type,
            "description": bansos.Description,
            "expiry_date": bansos.ExpiryDate,
            "image_url": bansos.ImageURL,
        },
    })
}

func (h *BansosHandler) GetBansosHandler(c echo.Context) error {
    bansos, err := h.bansosRepository.GetBansos()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to retrieve bansos",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "code": http.StatusOK,
        "status": "success",
        "message": "Bansos retrieved successfully",
        "data": bansos,
    })
}

func (h *BansosHandler) GetBansosByIDHandler(c echo.Context) error {
    // Get ID from the request URL
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid ID parameter",
        })
    }

    bansos, err := h.bansosRepository.GetBansosByID(id)
    if err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{
            "code":    http.StatusNotFound,
            "status":  "error",
            "message": "Bansos not found",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "code": http.StatusOK,
        "status": "success",
        "message": "Bansos retrieved successfully",
        "data": bansos,
    })
}

func (h *BansosHandler) UpdateBansosHandler(c echo.Context) error {
    var bansos model.Bansos
    var exist_bansos *model.Bansos

    // Get ID from the request URL
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid ID parameter",
        })
    }

    // Exist bansos
    exist_bansos, _ = h.bansosRepository.GetBansosByID(id)
    if exist_bansos == nil {
        return c.JSON(http.StatusNotFound, echo.Map{
            "code":    http.StatusNotFound,
            "status":  "error",
            "message": "Bansos not found",
        })
    }

    // Set ID in the bansos model
    bansos.ID = id
    bansos.Name = exist_bansos.Name
    bansos.Type = exist_bansos.Type
    bansos.Description = exist_bansos.Description
    bansos.ExpiryDate = exist_bansos.ExpiryDate
    bansos.ImageURL = exist_bansos.ImageURL

    // Bind payload
    if err := c.Bind(&bansos); err != nil || bansos.Name == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Retrieve the file from the form data
    file, fileHeader, err := c.Request().FormFile("file")
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

        // Set the image URL in the Bansos model
        bansos.ImageURL = imageURL
    }

    // Assuming you have a method in your repository to update bansos
    if err := h.bansosRepository.UpdateBansos(&bansos); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to update bansos",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Bansos updated successfully",
        "data": echo.Map{
            "id": bansos.ID,
            "name": bansos.Name,
            "type": bansos.Type,
            "description": bansos.Description,
            "expiry_date": bansos.ExpiryDate,
            "image_url": bansos.ImageURL,
        },
    })
}

func (h *BansosHandler) DeleteBansosHandler(c echo.Context) error {
    var exist_bansos *model.Bansos

    // Get ID from the request URL
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid ID parameter",
        })
    }

    // Exist bansos
    exist_bansos, _ = h.bansosRepository.GetBansosByID(id)
    if exist_bansos == nil {
        return c.JSON(http.StatusNotFound, echo.Map{
            "code":    http.StatusNotFound,
            "status":  "error",
            "message": "Bansos not found",
        })
    }

    bansos, _ := h.bansosRepository.GetBansosByID(id)
    if err := h.bansosRepository.DeleteBansos(bansos); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to delete bansos",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Bansos deleted successfully",
        "data": echo.Map{
            "id":  bansos.ID,
        },
    })
}