package repository

import (
    "gorm.io/gorm"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type BansosRepository interface {
	AddBansos(bansos *model.Bansos) error
    GetBansos() ([]*model.Bansos, error)
    GetBansosByID(id int) (*model.Bansos, error)
    UpdateBansos(bansos *model.Bansos) error
    DeleteBansos(bansos *model.Bansos) error
}

type BansosRepositoryImpl struct {
    db *gorm.DB
}

func NewBansosRepository(db *gorm.DB) BansosRepository {
    return &BansosRepositoryImpl{
        db: db,
    }
}

func (r *BansosRepositoryImpl) AddBansos(bansos *model.Bansos) error {
    return r.db.Create(bansos).Error
}

func (r *BansosRepositoryImpl) GetBansos() ([]*model.Bansos, error) {
    var bansos []*model.Bansos
	if err := r.db.Table("bansos").Find(&bansos).Error; err != nil {
		return nil, err
	}
	return bansos, nil
}

func (r *BansosRepositoryImpl) GetBansosByID(id int) (*model.Bansos, error) {
	var bansos model.Bansos
	if err := r.db.Table("bansos").Where("id = ?", id).First(&bansos).Error; err != nil {
		return nil, err
	}
	return &bansos, nil
}

func (r *BansosRepositoryImpl) UpdateBansos(bansos *model.Bansos) error {
    
    return r.db.Model(&model.Bansos{}).Where("id = ?", bansos.ID).Updates(bansos).Error

}

func (r *BansosRepositoryImpl) DeleteBansos(bansos *model.Bansos) error {
    
    return r.db.Delete(bansos).Error

}