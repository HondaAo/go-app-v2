package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserEntity struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Find user query
type UserEntityWithToken struct {
	User  *UserEntity
	Token string
}

// Hash user password with bcrypt
func (u *UserEntity) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *UserEntity) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}
