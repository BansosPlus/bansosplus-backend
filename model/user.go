package model

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type IncomeEnum string

const (
	Income_500rb    IncomeEnum = "<500 ribu"
	Income500rb_1jt IncomeEnum = "500 ribu-1 juta"
	Income1jt_1_5jt IncomeEnum = "1 juta-1.5 juta"
	Income1_5jt_2   IncomeEnum = ">1.5 juta"
)

type FloorAreaEnum string

const (
	FloorArea_8 FloorAreaEnum = "Diatas 8m²"
	FloorArea8_ FloorAreaEnum = "Dibawah 8m²"
)

type WallQualityEnum string

const (
	WallQualityBad    WallQualityEnum = "Buruk"
	WallQualityNormal WallQualityEnum = "Normal"
	WallQualityGood   WallQualityEnum = "Bagus"
)

type NumberOfMealsEnum string

const (
	NumberOfMeals_0 NumberOfMealsEnum = "0"
	NumberOfMeals_1 NumberOfMealsEnum = "1"
	NumberOfMeals_2 NumberOfMealsEnum = "2"
	NumberOfMeals_3 NumberOfMealsEnum = "3"
)

type FuelEnum string

const (
	FuelKayuArang FuelEnum = "Kayu/Arang"
	FuelGasLPG    FuelEnum = "Gas/LPG"
)

type EducationEnum string

const (
	EducationSD      EducationEnum = "SD"
	EducationSMP     EducationEnum = "SMP"
	EducationSMA     EducationEnum = "SMA"
	EducationSarjana EducationEnum = "Sarjana"
)

type TotalAssetEnum string

const (
	TotalAsset_500rb    TotalAssetEnum = "<500 ribu"
	TotalAsset500rb_1jt TotalAssetEnum = "500 ribu-1 juta"
	TotalAsset1jt_1_5jt TotalAssetEnum = "1 juta-1.5 juta"
	TotalAsset1_5jt_2   TotalAssetEnum = ">1.5 juta"
)

type TreatmentEnum string

const (
	TreatmentMampu      TreatmentEnum = "Mampu"
	TreatmentTidakMampu TreatmentEnum = "Tidak Mampu"
)

type NumberOfDependentsEnum string

const (
	NumberOfDependents0  NumberOfDependentsEnum = "0"
	NumberOfDependents1  NumberOfDependentsEnum = "1"
	NumberOfDependents2  NumberOfDependentsEnum = "2"
	NumberOfDependents2_ NumberOfDependentsEnum = ">2"
)

type RoleEnum string

const (
	RoleAdmin RoleEnum = "admin"
	RoleUser  RoleEnum = "user"
)

type User struct {
	ID                 int                    `json:"id"                      gorm:"primary_key;auto_increment"`
	Name               string                 `json:"name"                    gorm:"size:255;not null"`
	Nik                string                 `json:"nik"                     gorm:"size:255;null"`
	NoKK               string                 `json:"nokk"                    gorm:"size:255;null"`
	Income             IncomeEnum             `json:"income"                  gorm:"type:enum('<500 ribu','500 ribu-1 juta','1 juta-1.5 juta','>1.5 juta');null"`
	FloorArea          FloorAreaEnum          `json:"floor_area"              gorm:"type:enum('Diatas 8m²','Dibawah 8m²');null"`
	WallQuality        WallQualityEnum        `json:"wall_quality"            gorm:"type:enum('Buruk','Normal','Bagus');null"`
	NumberOfMeals      NumberOfMealsEnum      `json:"number_of_meals"         gorm:"type:enum('0','1','2','3');null"`
	Fuel               FuelEnum               `json:"fuel"                    gorm:"type:enum('Kayu/Arang','Gas/LPG');null"`
	Education          EducationEnum          `json:"education"               gorm:"type:enum('SD','SMP','SMA','Sarjana');null"`
	TotalAsset         TotalAssetEnum         `json:"total_asset"             gorm:"type:enum('<500 ribu','500 ribu-1 juta','1 juta-1.5 juta','>1.5 juta');null"`
	Treatment          TreatmentEnum          `json:"treatment"               gorm:"type:enum('Mampu','Tidak Mampu');null"`
	NumberOfDependents NumberOfDependentsEnum `json:"number_of_dependents"    gorm:"type:enum('0','1','2','>2');null"`
	Email              string                 `json:"email"                   gorm:"size:255;not null;unique"`
	Phone              string                 `json:"phone"                   gorm:"size:255;not null;unique"`
	Password           string                 `json:"password"                gorm:"size:255;not null"`
	Role               RoleEnum               `json:"role"                    gorm:"type:enum('admin','user');not null;default:'user'"`
	ImageURL           string                 `json:"image_url"               gorm:"size:255;null"`
	CreatedAt          time.Time              `json:"created_at"              gorm:"type:datetime;autoCreateTime"`
	UpdatedAt          time.Time              `json:"updated_at"              gorm:"type:datetime;autoUpdateTime"`
}

func ValidateEmail(email string) bool {
	// Regular expression for basic email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func ValidatePhone(phone string) bool {
	// Regular expression for phone number validation with only numeric digits and a minimum of 10 digits
	pattern := `^[0-9]{10,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phone)
}

func (u *User) HashPassword(plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) VerifyPassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}
