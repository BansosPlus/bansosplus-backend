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

func (r *FeedbackRepositoryImpl) GetFeedbackByUserID(id int) (feedback *model.Feedback) error {
    var feedbacks []*model.Feedback
	if err := r.db.Table("feedbacks").Where("user_id = ?", id).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *FeedbackRepositoryImpl) GetFeedbackByBansosID(id int) (feedback *model.Feedback) error {
    var feedbacks []*model.Feedback
	if err := r.db.Table("feedbacks").Where("bansos_id = ?", id).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *FeedbackRepositoryImpl) UpdateFeedback(feedback *model.Feedback) error {
    
    return r.db.Model(&model.Feedback{}).Where("id = ?", feedback.ID).Updates(feedback).Error

}

func (r *FeedbackRepositoryImpl) DeleteFeedback(feedback *model.Feedback) error {
    
    return r.db.Delete(feedback).Error

}