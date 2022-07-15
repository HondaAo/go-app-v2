package driver

import (
	"context"
	"database/sql"

	"github.com/HondaAo/video-app/pkg/auth/model"
	"github.com/pkg/errors"
)

// Auth Repository
type authRepo struct {
	db *sql.DB
}

// Auth Repository constructor
func NewAuthRepository(db *sql.DB) Repository {
	return &authRepo{db: db}
}

// Create new user
func (r *authRepo) Register(ctx context.Context, user *model.User) (*model.User, error) {
	u := &model.User{}
	stmt, err := r.db.Prepare("INSERT INTO user(first_name,last_name,email,password,role,country) VALUES(?,?,?,?) RETURNING *")
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.User.Insert Error")
	}
	defer stmt.Close()

	err = stmt.QueryRow(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Country).Scan(&u)
	return u, nil
}

// Find user by email
func (r *authRepo) FindByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	foundUser := &model.User{}
	row, err := r.db.Query(`SELECT user_id, first_name, last_name, email, role, country, created_at, updated_at, password FROM users WHERE email = $1`)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail. Error")
	}
	row.Scan(foundUser)
	return foundUser, nil
}
