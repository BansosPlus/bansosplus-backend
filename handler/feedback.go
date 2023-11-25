package handler

import (
    "net/http"
    "github.com/labstack/echo"
    
    "github.com/BansosPlus/bansosplus-backend.git/middleware"
    "github.com/BansosPlus/bansosplus-backend.git/model"
    "github.com/BansosPlus/bansosplus-backend.git/repository"
)

type FeedbackHandler struct {
    feedbackRepository repository.FeedbackRepository
}

func NewFeedbackHandler(feedbackRepository repository.FeedbackRepository) *FeedbackHandler {
    return &FeedbackHandler{
        feedbackRepository: feedbackRepository,
    }
}

func (h *FeedbackHandler) AddFeedbackHandler(c echo.Context) error {
    var feedback model.Feedback

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
    feedback.UserID = int(userID)

    // Bind payload
    if err := c.Bind(&feedback); err != nil || feedback.BansosID == "" || feedback.Score == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    if err := h.feedbackRepository.AddFeedback(&feedback); err != nil {
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
            "feedback_id":  feedback.ID,
            "user_id":      feedback.UserID,
            "score":        feedback.Score,
            "description":  feedback.Description,
        },
    })
}


func (h *FeedbackHandler) GetFeedbackByUserIDHandler(c echo.Context) error {

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

    feedbacks, err := h.feedbackRepository.GetFeedbacksByUserID(int(userID))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to retrieve feedbacks",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "code": http.StatusOK,
        "status": "success",
        "message": "Feedbacks retrieved successfully",
        "data": feedbacks,
    })
}

func (h *FeedbackHandler) GetFeedbackByBansosIDHandler(c echo.Context) error {
    // Bind payload
    var request struct {
        BansosID int `json:"bansos_id" form:"bansos_id"`
    }

    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    feedbacks, err := h.feedbackRepository.GetFeedbacksByBansosID(int(userID))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to retrieve feedbacks",
        })
    }

    return c.JSON(http.StatusOK, echo.Map{
        "code": http.StatusOK,
        "status": "success",
        "message": "Feedbacks retrieved successfully",
        "data": feedbacks,
    })
}

func (h *FeedbackHandler) UpdateFeedbackHandler(c echo.Context) error {
    var feedback model.Feedback

    // Bind payload
    if err := c.Bind(&feedback); err != nil || feedback.ID == 0 || feedback.Score == ""  {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Assuming you have a method in your repository to update feedback
    if err := h.feedbackRepository.UpdateFeedback(&feedback); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to update feedback",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Feedback updated successfully",
        "data": echo.Map{
            "feedback_id":  feedback.ID,
            "user_id":      feedback.UserID,
            "bansos_id":    feedback.BansosID,
            "score":        feedback.Score,
            "description":  feedback.Description,
        },
    })
}

func (h *FeedbackHandler) DeleteFeedbackHandler(c echo.Context) error {
    var feedback model.Feedback

    // Bind payload
    if err := c.Bind(&feedback); err != nil || feedback.ID == 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code":    http.StatusBadRequest,
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Assuming you have a method in your repository to delete feedback
    if err := h.feedbackRepository.DeleteFeedback(&feedback); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "code":    http.StatusInternalServerError,
            "status":  "error",
            "message": "Failed to delete feedback",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code":    http.StatusOK,
        "status":  "success",
        "message": "Feedback deleted successfully",
        "data": echo.Map{
            "feedback_id":  feedback.ID,
            "user_id":      feedback.UserID,
            "score":        feedback.Score,
            "description":  feedback.Description,
        },
    })
}