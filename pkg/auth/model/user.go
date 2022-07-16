package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id" redis:"user_id"`
	FirstName string    `json:"first_name" db:"first_name" redis:"first_name"`
	LastName  string    `json:"last_name" db:"last_name" redis:"last_name"`
	Email     string    `json:"email" db:"email" redis:"email"`
	Password  string    `json:"password" db:"password" redis:"password"`
	Role      string    `json:"role" db:"role" redis:"role"`
	Country   string    `json:"country" db:"country" redis:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" redis:"updated_at"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	u.Role = strings.ToLower(strings.TrimSpace(u.Role))

	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}
