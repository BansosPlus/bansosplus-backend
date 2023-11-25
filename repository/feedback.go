package repository

import (
    "gorm.io/gorm"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type FeedbackRepository interface {
	AddFeedback(feedback *model.Feedback) error
}

type FeedbackRepositoryImpl struct {
    db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
    return &FeedbackRepositoryImpl{
        db: db,
    }
}

func (r *FeedbackRepositoryImpl) AddFeedback(feedback *model.Feedback) error {
    return r.db.Create(feedback).Error
}