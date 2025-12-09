package service

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/data/response"
	"api/models/db"
	"api/util"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
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
func (a *AuthService) BuyerRegister(ctx context.Context, req request.RegisterRequest) error {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &errors.AppError{Message: "failed to hash password", Code: http.StatusInternalServerError}
	}
	tx, err := a.Conn.BeginTx(ctx, nil)
	if err != nil {
		return &errors.AppError{Message: "failed to start transaction", Code: http.StatusInternalServerError}
	}
	savedUser, err := a.Queries.GetUserByEmail(ctx, req.Email)
	if err == nil && savedUser.ID != 0 {
		return &errors.AppError{Message: "Email already in used", Code: http.StatusConflict}
	}

	qtx := a.Queries.WithTx(tx)

	data := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
	}
	user, err := qtx.CreateUser(
		ctx, data,
	)
	if err != nil {
		fmt.Print(err)
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
		fmt.Print(err)
		return &errors.AppError{Message: "failed to create account", Code: http.StatusInternalServerError}
	}
	_, err = qtx.CreateBuyer(ctx, user.ID)
	if err != nil {
		tx.Rollback()
		return &errors.AppError{Message: "failed to create buyer", Code: http.StatusInternalServerError}
	}

	if err := tx.Commit(); err != nil {
		return &errors.AppError{Message: "failed to commit transaction", Code: http.StatusInternalServerError}
	}
	return nil
}

func (a *AuthService) BuyerLogin(ctx context.Context, req request.LoginRequest, ip string) (response.UserResponse, error) {
	var res response.UserResponse
	ip = "" + ip
	{
		user, err := a.Queries.GetUserByAccountID(ctx, req.Email)
		if err != nil {
			return res, &errors.AppError{Message: "user not found", Code: http.StatusNotFound}
		}
		if user.Provider != db.ProviderCREDENTIALS {
			return res, &errors.AppError{Message: "Email already in used", Code: http.StatusUnauthorized}
		}
		isMatch := util.ComparePassword(user.Password.String, req.Password)
		if !isMatch {
			return res, &errors.AppError{Message: "invalid password", Code: http.StatusUnauthorized}
		}

		refreshToken, claims, err := util.GenerateRefreshToken(user.ID, ip)
		if err != nil {
			return res, &errors.AppError{Message: "failed to generate refresh token", Code: http.StatusInternalServerError}
		}
		accessToken, _, err := util.GenerateAccessToken(user.ID, ip, claims.Permissions)
		if err != nil {
			return res, &errors.AppError{Message: "failed to generate access token", Code: http.StatusInternalServerError}
		}

		a.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			ID:        claims.ID,
			Token:     refreshToken,
			UserID:    user.ID,
			IpAddress: sql.NullString{String: ip, Valid: true},
			ExpiresAt: claims.ExpiresAt.Time,
		})
		res = response.UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			Username:       user.Username,
			RefreshToken:   refreshToken,
			AccessToken:    accessToken,
			TokenExpiresAt: claims.ExpiresAt.Time,
		}
	}
	return res, nil
}
func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string, ip string) (response.TokenResponse, error) {
	var res response.TokenResponse
	{
		claims, err := util.ParseJWT(refreshToken)
		if err != nil {
			return res, &errors.AppError{Message: "invalid refresh token", Code: http.StatusUnauthorized}
		}
		savedToken, err := a.Queries.GetRefreshToken(ctx, refreshToken)
		if err != nil {
			return res, &errors.AppError{Message: "refresh token not found", Code: http.StatusUnauthorized}
		}

		if savedToken.Revoked {
			return res, &errors.AppError{Message: "refresh token revoked", Code: http.StatusUnauthorized}
		}
		a.Queries.UpdateRefreshTokenRevoked(
			ctx, db.UpdateRefreshTokenRevokedParams{
				ID:       savedToken.ID,
				Revoked:  true,
				LastUsed: sql.NullTime{Time: time.Now(), Valid: true},
			},
		)
		newRefreshToken, newClaims, err := util.GenerateRefreshToken(claims.UserID, ip)
		if err != nil {
			return res, &errors.AppError{Message: "failed to generate refresh token", Code: http.StatusInternalServerError}
		}
		if err := a.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			ID:        newClaims.ID,
			Token:     newRefreshToken,
			UserID:    int64(claims.UserID),
			IpAddress: sql.NullString{String: ip, Valid: true},
			ExpiresAt: newClaims.ExpiresAt.Time,
		}); err != nil {
			return res, &errors.AppError{Message: "failed to save refresh token", Code: http.StatusInternalServerError}
		}
		if err := a.Queries.DeleteRefreshTokensByID(ctx, claims.ID); err != nil {
			return res, &errors.AppError{Message: "failed to delete old refresh token", Code: http.StatusInternalServerError}
		}

		newAccessToken, exp, err := util.GenerateAccessToken(claims.UserID, ip, claims.Permissions)
		if err != nil {
			return res, &errors.AppError{Message: "failed to generate access token", Code: http.StatusInternalServerError}
		}
		res = response.TokenResponse{
			RefreshToken:   newRefreshToken,
			AccessToken:    newAccessToken,
			TokenExpiresAt: exp,
		}
	}
	return res, nil
}

func (a *AuthService) GetUserByID(ctx context.Context, userID int64) (response.UserDetails, error) {
	user, err := a.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return response.UserDetails{}, &errors.AppError{Message: "user not found", Code: http.StatusNotFound}
	}
	return response.UserDetails{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := util.ParseJWT(refreshToken)
	if err != nil {
		return &errors.AppError{Message: "invalid refresh token", Code: http.StatusUnauthorized}
	}
	if err := a.Queries.DeleteRefreshTokensByID(ctx, claims.ID); err != nil {
		return &errors.AppError{Message: "failed to delete refresh token", Code: http.StatusInternalServerError}
	}
	return nil
}
