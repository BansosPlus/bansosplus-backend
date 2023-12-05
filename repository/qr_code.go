package repository

import (
    "gorm.io/gorm"

    "github.com/BansosPlus/bansosplus-backend.git/model"
)

type QRCodeRepository interface {
	AddQRCode(qrCode *model.QRCode) error
    GetQRCodeByID(id int) (*model.QRCode, error)
    GetQRCodeByUUID(uuid string) (*model.QRCode, error)
}

type QRCodeRepositoryImpl struct {
    db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
    return &QRCodeRepositoryImpl{
        db: db,
    }
}

func (r *QRCodeRepositoryImpl) AddQRCode(qrCode *model.QRCode) error {
    return r.db.Create(qrCode).Error
}

func (r *QRCodeRepositoryImpl) GetQRCodeByID(id int) (*model.QRCode, error) {
	var qrCode model.QRCode
	if err := r.db.Table("qr_codes").Where("id = ?", id).First(&qrCode).Error; err != nil {
		return nil, err
	}
	return &qrCode, nil
}

func (r *QRCodeRepositoryImpl) GetQRCodeByUUID(uuid string) (*model.QRCode, error) {
	var qrCode model.QRCode
	if err := r.db.Table("qr_codes").Where("uuid = ?", uuid).First(&qrCode).Error; err != nil {
		return nil, err
	}
	return &qrCode, nil
}