package repository

import (
	"gorm.io/gorm"
	
	"github.com/BansosPlus/bansosplus-backend.git/model"
)

type UserRepository interface {
	GetUserByID(id int) (*model.User, error)
	UpdateUser(user *model.User) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *model.User) error {
    
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error

}