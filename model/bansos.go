package model

import (
	"time"
)

type Bansos struct {
    ID        	int             `json:"id"          	gorm:"primary_key;auto_increment"`
    Name      	string          `json:"name"        	gorm:"size:255;not null"`
    Type      	string  		`json:"type"        	gorm:"size:255;not null;default:'groceries'"`
    Description string     		`json:"description"     gorm:"size:255;null"`
	ExpiryDate  time.Time  		`json:"expiry_date"     gorm:"type:date;null"`
    CreatedAt 	time.Time       `json:"created_at"  	gorm:"type:datetime;autoCreateTime"`
    UpdatedAt 	time.Time       `json:"updated_at"  	gorm:"type:datetime;autoUpdateTime"`
}