package model

import (
    "time"
)

type QRCode struct {
    ID        	                int             `json:"id"                          gorm:"primary_key;auto_increment"`
    BansosRegistrationID        int             `json:"bansos_registration_id"      gorm:"not null"     sql:"type:int REFERENCES bansos_registrations(id)"`
    Uuid                        string          `json:"uuid"                        gorm:"size:255;not null"`
    Score      	                int     	    `json:"score"                       gorm:"null"`
    Description                 string     		`json:"description"                 gorm:"size:255;null"`
    CreatedAt 	                time.Time       `json:"created_at"                  gorm:"type:datetime;autoCreateTime"`
    UpdatedAt 	                time.Time       `json:"updated_at"                  gorm:"type:datetime;autoUpdateTime"`
}