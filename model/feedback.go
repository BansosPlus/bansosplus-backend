package model

import (
	"time"
)

type Feedback struct {
    ID        	int             `json:"id"          	gorm:"primary_key;auto_increment"`
    UserID      int          	`json:"user_id"        	gorm:"not null" 	sql:"type:int REFERENCES users(id)"`
    Score      	int     	    `json:"score"         	gorm:"not null"`
    Description string     		`json:"description"     gorm:"size:255;null"`
    CreatedAt 	time.Time       `json:"created_at"  	gorm:"type:datetime;autoCreateTime"`
    UpdatedAt 	time.Time       `json:"updated_at"  	gorm:"type:datetime;autoUpdateTime"`
}