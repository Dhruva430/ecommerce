package service

import (
	"api/errors"
	"api/internals/data/request"
	"api/models/db"
	"api/util"
	"context"
	"database/sql"
	"net/http"
)

type AuthService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewAuthService(queries *db.Queries, conn *sql.DB) AuthService {
	return AuthService{
		Queries: queries,
		Conn:    conn,
	}
}
func (a *AuthService) Register(ctx context.Context, req request.RegisterRequest) error {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &errors.AppError{Message: "failed to hash password", Code: http.StatusInternalServerError}
	}
	tx, err := a.Conn.BeginTx(ctx, nil)
	if err != nil {
		return &errors.AppError{Message: "failed to start transaction", Code: http.StatusInternalServerError}
	}
	defer tx.Rollback()

	qtx := a.Queries.WithTx(tx)

	data := db.CreateUserParams{
		Email: req.Email,
		Role:  db.Role(req.Role),
	}
	user, err := qtx.CreateUser(
		ctx, data,
	)
	if err != nil {
		return &errors.AppError{Message: "failed to create user", Code: http.StatusInternalServerError}
	}
	account := db.CreateAccountParams{
		AccountID: user.Email,
		Password:  sql.NullString{String: hashedPassword, Valid: true},
		Provider:  db.ProviderCREDENTIALS,
		UserID:    user.ID,
	}
	if err = qtx.CreateAccount(
		ctx, account,
	); err != nil {
		return &errors.AppError{Message: "failed to create account", Code: http.StatusInternalServerError}
	}
	if err := tx.Commit(); err != nil {
		return &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}
	return nil

}

func (a *AuthService) Login() error {
	return &errors.AppError{Message: "login failed", Code: http.StatusUnauthorized}
}
