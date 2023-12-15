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
	Income              IncomeEnum              `json:"income"                  gorm:"type:enum('<500 ribu','500 ribu-1 juta','1 juta-1.5 juta','>1.5 juta');null"`
	FloorArea           FloorAreaEnum           `json:"floor_area"              gorm:"type:enum('Diatas 8m²','Dibawah 8m²');null"`
	WallQuality         WallQualityEnum         `json:"wall_quality"            gorm:"type:enum('Buruk','Normal','Bagus');null"`
	NumberOfMeals       NumberOfMealsEnum       `json:"number_of_meals"         gorm:"type:enum('0','1','2','3');null"`
	Fuel                FuelEnum                `json:"fuel"                    gorm:"type:enum('Kayu/ Arang','Gas/ LPG');null"`
	Education           EducationEnum           `json:"education"               gorm:"type:enum('SD','SMP','SMA','Sarjana');null"`
	TotalAsset          TotalAssetEnum          `json:"total_asset"             gorm:"type:enum('<500 ribu','500 ribu-1 juta','1 juta-1.5 juta','>1.5 juta');null"`
	Treatment           TreatmentEnum           `json:"treatment"               gorm:"type:enum('Mampu','Tidak Mampu');null"`
	NumberOfDependents  NumberOfDependentsEnum  `json:"number_of_dependents"    gorm:"type:enum('0','1','2','>2');null"`
	Status     StatusEnum `json:"status"          gorm:"type:enum('ON_PROGRESS','ACCEPTED','REJECTED','TAKEN');not null;default:'ON_PROGRESS'"`
	Point			decimal.Decimal `json:"point"      gorm:"type:decimal(10,2);not null;default:0"`
	ApprovalAt time.Time  `json:"approval_at"     gorm:"type:datetime;autoCreateTime"`
	CreatedAt  time.Time  `json:"created_at"      gorm:"type:datetime;autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at"      gorm:"type:datetime;autoUpdateTime"`
}
