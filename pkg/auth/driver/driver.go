package driver

import (
	"context"
	"database/sql"

	"github.com/HondaAo/video-app/pkg/auth/model"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "Register")
	defer span.Finish()

	u := &model.User{}
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users(user_id,first_name,last_name,email,password,role,country) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.User.Insert Error")
	}
	defer stmt.Close()

	id := uuid.New()
	user.UserID = id.String()

	if err = stmt.QueryRowContext(ctx, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Country).Scan(u); err != nil {
		return nil, err
	}
	return u, nil
}

// Find user by email
func (r *authRepo) FindByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FindUserByEmail")
	defer span.Finish()

	foundUser := &model.User{}
	row, err := r.db.QueryContext(ctx, `SELECT user_id, password FROM users WHERE email = ?`, user.Email)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail. Error")
	}

	for row.Next() {
		if err = row.Scan(&foundUser.UserID, &foundUser.Password); err != nil {
			return nil, errors.Wrap(err, "authRepo.FindByEmail. Error")
		}
	}

	return foundUser, nil
}

// Get user by id
func (r *authRepo) GetByID(ctx context.Context, userID string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.GetByID")
	defer span.Finish()

	user := &model.User{}
	row, err := r.db.QueryContext(ctx, `SELECT * FROM users WHERE user_id = ?`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail. Error")
	}

	for row.Next() {
		if err = row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Country, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "authRepo.FindByEmail. Error")
		}
	}
	return user, nil
}
