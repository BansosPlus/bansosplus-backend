package repository

import (
    "gorm.io/gorm"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type BansosRegistrationRepository interface {

}

type BansosRegistrationRepositoryImpl struct {
    db *gorm.DB
}

func NewBansosRegistrationRepository(db *gorm.DB) BansosRegistrationRepository {
    return &BansosRegistrationRepositoryImpl{
        db: db,
    }
}

func (r *BansosRegistrationRepositoryImpl) RegisterBansos(bansosRegistration *model.BansosRegistration) error {
    return r.db.Create(bansosRegistration).Error
}