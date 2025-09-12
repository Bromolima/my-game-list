package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	Email    string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Username string `gorm:"type:varchar(100);not null"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
