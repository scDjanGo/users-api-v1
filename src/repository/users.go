package repository

import (
	"database/sql"
	"errors"
	"strings"

	"users-api-v1/src/models"
)

func (r *UserRepo) GetByID(id int64) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = ?`, id)
	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return user, err
}

func (r *UserRepo) GetAll(userQuery models.ListUserQuery) ([]*models.User, error) {
	query := "SELECT " + userColumns + " FROM users"

	var conditions []string
	var args []any
	if userQuery.HasBoss != nil {
		if *userQuery.HasBoss {
			conditions = append(conditions, "boss_id IS NOT NULL")
		} else {
			conditions = append(conditions, "boss_id IS NULL")
		}
	}

	if userQuery.BossID != nil {
		conditions = append(conditions, "boss_id = ?")
		args = append(args, *userQuery.BossID)
	}

	if userQuery.HasAvatar != nil {
		if *userQuery.HasAvatar {
			conditions = append(conditions, "avatar IS NOT NULL")
		} else {
			conditions = append(conditions, "avatar IS NULL")
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sortColumn, ok := allowedSortColumns[userQuery.SortBy]
	if !ok {
		sortColumn = "id"
	}
	order := "ASC"
	if strings.EqualFold(userQuery.Order, "desc") {
		order = "DESC"
	}
	query += " ORDER BY " + sortColumn + " " + order

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)

	for rows.Next() {
		u, err := scanUser(rows)

		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) Create(req models.CreateUserRequest, passwordHash string) (*models.User, error) {
	row := r.db.QueryRow(`
	INSERT INTO users (first_name, last_name, login, password, avatar, date_of_birth) 
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING `+userColumns, req.FirstName, req.LastName, req.Login, passwordHash, req.Avatar, req.DateOfBirth)

	user, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrLoginTaken
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Update(id int64, req models.UpdatedUserRequest) (*models.User, error) {
	row := r.db.QueryRow(`
	UPDATE users SET
		first_name = COALESCE(?, first_name),
		last_name = COALESCE(?, last_name),
		avatar = COALESCE(?, avatar),
		boss_id = COALESCE(?, boss_id),
		date_of_birth = COALESCE(?, date_of_birth),
		updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
	WHERE id = ?
	RETURNING `+userColumns, req.FirstName, req.LastName, req.Avatar, req.BossID, req.DateOfBirth, id)

	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return user, err
}

func (r *UserRepo) Delete(id int64) error {

	res, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil

}

func (r *UserRepo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM users")
	return err
}
