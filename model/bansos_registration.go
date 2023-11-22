package model

import (
    "time"
)

type StatusEnum string
const (
	StatusOnProgress    StatusEnum = "ON_PROGRESS"
    StatusAccepted      StatusEnum = "ACCEPTED"
    StatusRejected      StatusEnum = "REJECTED"
)

type BansosRegistration struct {
    ID        	int             `json:"id"          	gorm:"primary_key;auto_increment"`
    BansosID    int          	`json:"bansos_id"       gorm:"not null"     sql:"type:int REFERENCES bansos(id)"`
    UserID      int          	`json:"user_id"        	gorm:"not null"     sql:"type:int REFERENCES users(id)"`
    Name        string          `json:"name"            gorm:"size:255;null"`
    Nik         string          `json:"nik"             gorm:"size:255;null;unique"`
    NoKK        string          `json:"nokk"            gorm:"size:255;null"`
    Status      StatusEnum      `json:"status"          gorm:"type:enum('ON_PROGRESS','ACCEPTED','REJECTED');not null;default:'ON_PROGRESS'"`
    ApprovalAt 	time.Time       `json:"approval_at"  	gorm:"type:datetime;autoCreateTime"`
    CreatedAt 	time.Time       `json:"created_at"  	gorm:"type:datetime;autoCreateTime"`
    UpdatedAt 	time.Time       `json:"updated_at"  	gorm:"type:datetime;autoUpdateTime"`
}