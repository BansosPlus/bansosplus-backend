package model

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Income             string                 `json:"income"                  gorm:"size:255;null"`
	FloorArea          string                 `json:"floor_area"              gorm:"size:255;null"`
	WallQuality        string                 `json:"wall_quality"            gorm:"size:255;null"`
	NumberOfMeals      string                 `json:"number_of_meals"         gorm:"size:255;null"`
	Fuel               string                 `json:"fuel"                    gorm:"size:255;null"`
	Education          string                 `json:"education"               gorm:"size:255;null"`
	TotalAsset         string                 `json:"total_asset"             gorm:"size:255;null"`
	Treatment          string                 `json:"treatment"               gorm:"size:255;null"`
	NumberOfDependents string                 `json:"number_of_dependents"    gorm:"size:255;null"`
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
