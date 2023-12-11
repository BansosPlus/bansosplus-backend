package repository

import (
    "gorm.io/gorm"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type FeedbackWithUsername struct {
    ID          uint   `json:"id"`
    UserID      uint   `json:"user_id"`
    UserName    string `json:"user_name"`
    BansosID    uint   `json:"bansos_id"`
    Score       int    `json:"score"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type FeedbackRepository interface {
	AddFeedback(feedback *model.Feedback) error
    GetFeedbackByBansosID(id int) ([]*FeedbackWithUsername, error)
    GetFeedbackByUserID(id int) ([]*model.Feedback, error)
    GetFeedbackByID(id int) (*model.Feedback, error)
    UpdateFeedback(feedback *model.Feedback) error
    DeleteFeedback(feedback *model.Feedback) error
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

func (r *FeedbackRepositoryImpl) GetFeedbackByID(id int) (*model.Feedback, error) {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *FeedbackRepositoryImpl) GetFeedbackByUserID(id int) ([]*model.Feedback, error) {
    var feedbacks []*model.Feedback
	if err := r.db.Table("feedbacks").Where("user_id = ?", id).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *FeedbackRepositoryImpl) GetFeedbackByBansosID(id int) ([]*FeedbackWithUsername, error) {
    var feedbacks []*FeedbackWithUsername
	if err := r.db.Table("feedbacks").
        Select("feedbacks.id, feedbacks.user_id, feedbacks.bansos_id, feedbacks.score, feedbacks.description, feedbacks.created_at, feedbacks.updated_at, users.name as user_name").
        Joins("INNER JOIN users ON feedbacks.user_id = users.id").
        Where("feedbacks.bansos_id = ?", id).
        Find(&feedbacks).Error; err != nil {
        return nil, err
    }
	
    return feedbacks, nil
}

func (r *FeedbackRepositoryImpl) UpdateFeedback(feedback *model.Feedback) error {
    
    return r.db.Model(&model.Feedback{}).Where("id = ?", feedback.ID).Updates(feedback).Error

}

func (r *FeedbackRepositoryImpl) DeleteFeedback(feedback *model.Feedback) error {
    
    return r.db.Delete(feedback).Error

}