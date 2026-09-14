package repository

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	myDB "users-api-v1/src/db"
	"users-api-v1/src/features"
	"users-api-v1/src/models"
)

type UserRepo struct {
	db *sql.DB
}

type scanner interface {
	Scan(dest ...any) error
}

var allowedSortColumns = map[string]string{
	"date_of_birth": "date_of_birth",
	"boss_id":       "boss_id",
	"created_at":    "created_at",
	"id":            "id",
}

var (
	ErrNotFound   = errors.New("Пользователь не найден")
	ErrLoginTaken = errors.New("Пользователь с таким логином уже существует")
)

const userColumns = `id, first_name, last_name, login, password, avatar, boss_id, created_at, updated_at, token, date_of_birth`

func NewUserRepo() (*UserRepo, error) {
	dbPath := features.GetEnv("DB_PATH", "./src/data/app.db")
	conn, err := myDB.Open(dbPath)

	if err != nil {
		return nil, err
	}

	return &UserRepo{db: conn}, nil
}

func scanUser(s scanner) (*models.User, error) {
	var u models.User
	err := s.Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Login, &u.Password,
		&u.Avatar, &u.BossID, &u.CreatedAt, &u.UpdatedAt, &u.Token, &u.DateOfBirth,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
