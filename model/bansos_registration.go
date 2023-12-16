package model

import (
	"time"
	"github.com/shopspring/decimal"
)

type StatusEnum string

const (
	StatusOnProgress StatusEnum = "ON_PROGRESS"
	StatusAccepted   StatusEnum = "ACCEPTED"
	StatusRejected   StatusEnum = "REJECTED"
	StatusTaken      StatusEnum = "TAKEN"
)

type BansosRegistration struct {
	ID         int        `json:"id"              gorm:"primary_key;auto_increment"`
	BansosID   int        `json:"bansos_id"       gorm:"not null"     sql:"type:int REFERENCES bansos(id)"`
	UserID     int        `json:"user_id"         gorm:"not null"     sql:"type:int REFERENCES users(id)"`
	Name       string     `json:"name"            gorm:"size:255;null"`
	Nik        string     `json:"nik"             gorm:"size:255;null"`
	NoKK       string     `json:"nokk"            gorm:"size:255;not null"`
	Income             string                 `json:"income"                  gorm:"size:255;null"`
	FloorArea          string                 `json:"floor_area"              gorm:"size:255;null"`
	WallQuality        string                 `json:"wall_quality"            gorm:"size:255;null"`
	NumberOfMeals      string                 `json:"number_of_meals"         gorm:"size:255;null"`
	Fuel               string                 `json:"fuel"                    gorm:"size:255;null"`
	Education          string                 `json:"education"               gorm:"size:255;null"`
	TotalAsset         string                 `json:"total_asset"             gorm:"size:255;null"`
	Treatment          string                 `json:"treatment"               gorm:"size:255;null"`
	NumberOfDependents string                 `json:"number_of_dependents"    gorm:"size:255;null"`
	Status     StatusEnum `json:"status"          gorm:"type:enum('ON_PROGRESS','ACCEPTED','REJECTED','TAKEN');not null;default:'ON_PROGRESS'"`
	Point			decimal.Decimal `json:"point"      gorm:"type:decimal(10,2);not null;default:0"`
	ApprovalAt time.Time  `json:"approval_at"     gorm:"type:datetime;autoCreateTime"`
	CreatedAt  time.Time  `json:"created_at"      gorm:"type:datetime;autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at"      gorm:"type:datetime;autoUpdateTime"`
}
