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

func (h *FeedbackHandler) AddFeedback(c echo.Context) error {
    var feedback model.Feedback

    // Bind payload
    if err := c.Bind(&feedback); err != nil || feedback.UserID == "" || feedback.Score == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "code": http.StatusBadRequest,
            "status": "error",
            "message": "Invalid request payload",
        })
    }

    // Success
    return c.JSON(http.StatusOK, echo.Map{
        "code": http.StatusOK,
        "status": "success",
        "message": "Feedback added successfully",
        "data": echo.Map{
            "feedback_id": feedback.ID,
            "user_id": feedback.UserID,
            "score": feedback.Score,
            "description": feedback.Description,
        },
    })
}