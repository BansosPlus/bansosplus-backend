package repository

import (
    "gorm.io/gorm"
    "time"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type BansosRegistrationWithBansos struct {
    model.BansosRegistration
    Bansos model.Bansos `gorm:"column:bansos_id;embedded"`
}

type BansosRegistrationRepository interface {
    RegisterBansos(bansosRegistration *model.BansosRegistration) error
    AcceptBansosRegis(bansosRegistration *model.BansosRegistration) error
    RejectBansosRegis(bansosRegistration *model.BansosRegistration) error
    GetBansosRegisByID(id int) (*model.BansosRegistration, error)
    GetBansosRegisByStatus(status string) ([]*model.BansosRegistration, error)
    GetBansosRegisByUserID(id int) ([]*BansosRegistrationWithBansos, error)
    GetBansosRegisByBansosID(id int) ([]*model.BansosRegistration, error)
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

func (r *BansosRegistrationRepositoryImpl) GetBansosRegisByID(id int) (*model.BansosRegistration, error) {
	var bansosRegistration model.BansosRegistration
	if err := r.db.Table("bansos_registrations").Where("id = ?", id).First(&bansosRegistration).Error; err != nil {
		return nil, err
	}
	return &bansosRegistration, nil
}

func (r *BansosRegistrationRepositoryImpl) GetBansosRegisByUserID(id int) ([]*BansosRegistrationWithBansos, error) {
    var result []*BansosRegistrationWithBansos
    rows, err := r.db.Table("bansos_registrations").
        Select("bansos_registrations.*, bansos.*").
        Joins("JOIN bansos ON bansos_registrations.bansos_id = bansos.id").
        Where("bansos_registrations.user_id = ?", id).
        Rows()

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var item BansosRegistrationWithBansos
        if err := r.db.ScanRows(rows, &item); err != nil {
            return nil, err
        }
        result = append(result, &item)
    }

    return result, nil
}

func (r *BansosRegistrationRepositoryImpl) AcceptBansosRegis(bansosRegistration *model.BansosRegistration) error {    
    return r.db.Model(&model.BansosRegistration{}).
        Where("id = ?", bansosRegistration.ID).
        Updates(map[string]interface{}{"status": "ACCEPTED", "approval_at": time.Now()}).
        Error
}

func (r *BansosRegistrationRepositoryImpl) RejectBansosRegis(bansosRegistration *model.BansosRegistration) error {    
    return r.db.Model(&model.BansosRegistration{}).
        Where("id = ?", bansosRegistration.ID).
        Updates(map[string]interface{}{"status": "REJECTED", "approval_at": time.Now()}).
        Error
}

func (r *BansosRegistrationRepositoryImpl) GetBansosRegisByStatus(status string) ([]*model.BansosRegistration, error) {
    var bansosRegistrations []*model.BansosRegistration
	if err := r.db.Table("bansos_registrations").Where("status = ?", status).Find(&bansosRegistrations).Error; err != nil {
		return nil, err
	}
	return bansosRegistrations, nil
}

func (r *BansosRegistrationRepositoryImpl) GetBansosRegisByBansosID(id int) ([]*model.BansosRegistration, error) {
    var bansosRegistrations []*model.BansosRegistration
	if err := r.db.Table("bansos_registrations").Where("bansos_id = ?", id).Find(&bansosRegistrations).Error; err != nil {
		return nil, err
	}
	return bansosRegistrations, nil
}