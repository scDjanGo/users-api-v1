package repository

import (
	"database/sql"
	"errors"

	"users-api-v1/src/models"
)

func (r *UserRepo) GetUserByToken(token string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE token = ?`, token)

	user, err := scanUser(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return user, err
}

func (r *UserRepo) UserRegistration(req models.RegistrationRequest, passwordHash string) (*models.User, error) {
	row := r.db.QueryRow(`
	INSERT INTO users (first_name, last_name, login, password, avatar, boss_id, date_of_birth)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING `+userColumns, req.FirstName, req.LastName, req.Login, passwordHash, req.Avatar, req.BossID, req.DateOfBirth)

	user, err := scanUser(row)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrLoginTaken
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByLogin(login string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE login = ?`, login)

	user, err := scanUser(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return user, err
}
