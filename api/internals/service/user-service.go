package service

import (
	"api/models/db"
	"database/sql"
)

type UserService struct {
	Queries *db.Queries
	Conn    *sql.DB
}

func NewUserService(queries *db.Queries, conn *sql.DB) UserService {
	return UserService{
		Queries: queries,
		Conn:    conn,
	}
}
