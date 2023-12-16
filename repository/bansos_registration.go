package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/BansosPlus/bansosplus-backend.git/model"
	"github.com/shopspring/decimal"
)

type BansosRegistrationWithBansos struct {
	ID         uint   `json:"id"`
	BansosName string `json:"bansos_name"`
	BansosType string `json:"type"`
	Status     string `json:"status"`
	ImageUrl   string `json:"image_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type DetailBansosRegistrationWithBansos struct {
	ID                uint            `json:"id"`
	BansosName        string          `json:"bansos_name"`
	BansosType        string          `json:"type"`
	BansosDescription string          `json:"description"`
	Status            string          `json:"status"`
	Point             decimal.Decimal `json:"point"`
	ImageUrl          string          `json:"image_url"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
}

type BansosRegistrationRepository interface {
	RegisterBansos(bansosRegistration *model.BansosRegistration) error
	AcceptBansosRegis(bansosRegistration *model.BansosRegistration, point decimal.Decimal) error
	RejectBansosRegis(bansosRegistration *model.BansosRegistration, point decimal.Decimal) error
	ValidateBansosRegis(bansosRegistration *model.BansosRegistration) error
	GetBansosRegisByID(id int) (*model.BansosRegistration, error)
	GetBansosRegisByStatus(status string) ([]*model.BansosRegistration, error)
	GetBansosRegisByUserID(id int, statuses []string) ([]*BansosRegistrationWithBansos, error)
	GetBansosRegisByBansosID(id int) ([]*model.BansosRegistration, error)
	GetDetailBansosRegisByID(id int) (*DetailBansosRegistrationWithBansos, error)
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

func (r *BansosRegistrationRepositoryImpl) GetDetailBansosRegisByID(id int) (*DetailBansosRegistrationWithBansos, error) {
	var result DetailBansosRegistrationWithBansos
	if err := r.db.Table("bansos_registrations").
		Select("bansos_registrations.id, bansos.name as bansos_name, bansos.type, bansos.description, bansos_registrations.status, bansos_registrations.point, bansos.image_url, bansos_registrations.created_at, bansos_registrations.updated_at").
		Joins("JOIN bansos ON bansos_registrations.bansos_id = bansos.id").
		Where("bansos_registrations.id = ?", id).
		First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BansosRegistrationRepositoryImpl) GetBansosRegisByUserID(id int, statuses []string) ([]*BansosRegistrationWithBansos, error) {
	var result []*BansosRegistrationWithBansos
	rows, err := r.db.Table("bansos_registrations").
		Select("bansos_registrations.id, bansos.name as bansos_name, bansos.type, bansos_registrations.status, bansos.image_url, bansos_registrations.created_at, bansos_registrations.updated_at").
		Joins("JOIN bansos ON bansos_registrations.bansos_id = bansos.id").
		Joins("JOIN users ON bansos_registrations.user_id = users.id").
		Where("bansos_registrations.user_id = ?", id).
		Where("bansos_registrations.status IN (?)", statuses).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item BansosRegistrationWithBansos
		if err := rows.Scan(&item.ID, &item.BansosName, &item.BansosType, &item.Status, &item.ImageUrl, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &item)
	}

	return result, nil
}

func (r *BansosRegistrationRepositoryImpl) AcceptBansosRegis(bansosRegistration *model.BansosRegistration, point decimal.Decimal) error {
	return r.db.Model(&model.BansosRegistration{}).
		Where("id = ?", bansosRegistration.ID).
		Updates(map[string]interface{}{"status": "ACCEPTED", "approval_at": time.Now(), "point": point}).
		Error
}

func (r *BansosRegistrationRepositoryImpl) RejectBansosRegis(bansosRegistration *model.BansosRegistration, point decimal.Decimal) error {
	return r.db.Model(&model.BansosRegistration{}).
		Where("id = ?", bansosRegistration.ID).
		Updates(map[string]interface{}{"status": "REJECTED", "approval_at": time.Now(), "point": point}).
		Error
}

func (r *BansosRegistrationRepositoryImpl) ValidateBansosRegis(bansosRegistration *model.BansosRegistration) error {
	return r.db.Model(&model.BansosRegistration{}).
		Where("id = ?", bansosRegistration.ID).
		Updates(map[string]interface{}{"status": "TAKEN", "approval_at": time.Now()}).
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
