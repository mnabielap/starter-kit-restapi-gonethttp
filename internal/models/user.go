package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	Password        string    `gorm:"not null" json:"-"` // json:"-" prevents password from being returned in API
	Role            string    `gorm:"default:'user'" json:"role"`
	IsEmailVerified bool      `gorm:"default:false" json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// BeforeCreate is a GORM hook that generates a UUID before saving
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// BeforeSave is a GORM hook that hashes the password if it has been modified
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" && len(u.Password) < 60 { 
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return
}

// ComparePassword checks if the provided password matches the hashed password
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}