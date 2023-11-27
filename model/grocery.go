package model

import (
    "time"
    "github.com/shopspring/decimal"
)

type Grocery struct {
    ID        	int                 `json:"id"              gorm:"not null"     sql:"type:int REFERENCES bansos(id)"`
    Income      decimal.Decimal     `json:"income"          gorm:"type:decimal(20,8);not null"`
    CreatedAt 	time.Time           `json:"created_at"      gorm:"type:datetime;autoCreateTime"`
    UpdatedAt 	time.Time           `json:"updated_at"      gorm:"type:datetime;autoUpdateTime"`
}