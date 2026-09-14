package models

type User struct {
	ID          int64   `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Login       string  `json:"login"`
	Password    string  `json:"-"`
	Avatar      *string `json:"avatar"`
	DateOfBirth string  `json:"date_of_birth"`
	BossID      *int64  `json:"boss_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Token       string  `json:"-"`
}

type CreateUserRequest struct {
	FirstName   string  `json:"first_name" form:"first_name" binding:"required,min=1"`
	LastName    string  `json:"last_name" form:"last_name" binding:"required,min=1"`
	Login       string  `json:"login" form:"login" binding:"required,min=3"`
	Password    string  `json:"password" form:"password" binding:"required,min=6"`
	Avatar      *string `json:"avatar" form:"-"`
	DateOfBirth string  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
}

type UpdatedUserRequest struct {
	FirstName   *string `json:"first_name" form:"first_name" binding:"omitempty,min=1"`
	LastName    *string `json:"last_name" form:"last_name" binding:"omitempty,min=1"`
	Avatar      *string `json:"avatar" form:"-"`
	BossID      *int64  `json:"boss_id" form:"boss_id" binding:"omitempty,gt=0"`
	DateOfBirth *string `json:"date_of_birth" form:"date_of_birth" binding:"omitempty"`
}

type LoginRequest struct {
	Login    string `json:"login" form:"login" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegistrationRequest struct {
	Login          string  `json:"login" form:"login" binding:"required,min=3"`
	Password       string  `json:"password" form:"password" binding:"required,min=6"`
	RepeatPassword string  `json:"repeat_password" form:"repeat_password" binding:"required,eqfield=Password"`
	FirstName      string  `json:"first_name" form:"first_name" binding:"required,min=1"`
	LastName       string  `json:"last_name" form:"last_name" binding:"required,min=1"`
	Avatar         *string `json:"avatar" form:"-"`
	DateOfBirth    string  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	BossID         *int64  `json:"boss_id" form:"boss_id" binding:"omitempty,gt=0"`
}

type ListUserQuery struct {
	SortBy    string `form:"sort_by"`
	Order     string `form:"order"`
	HasBoss   *bool  `form:"has_boss"`
	HasAvatar *bool  `form:"has_avatar"`
	BossID    *int64 `form:"boss_id"`
}

type UserWithToken struct {
	ID          int64   `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Login       string  `json:"login"`
	Avatar      *string `json:"avatar"`
	DateOfBirth string  `json:"date_of_birth"`
	BossID      *int64  `json:"boss_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Token       string  `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Success string `json:"success"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields"`
}
