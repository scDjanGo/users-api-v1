package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"users-api-v1/src/models"
	"users-api-v1/src/repository"
)

// RegisterUser godoc
// @Summary      Регистрация
// @Description  Публичная самостоятельная регистрация. Токен возвращается только один раз — сохраните его.
// @Tags         auth
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        body  body      models.RegistrationRequest  true  "Данные регистрации"
// @Success      201   {object}  models.UserWithToken
// @Failure      400   {object}  models.ValidationErrorResponse
// @Failure      409   {object}  models.ErrorResponse  "Логин уже занят"
// @Failure      413   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /auth/registration [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req models.RegistrationRequest

	if err := c.ShouldBind(&req); err != nil {
		bindError(c, err)
		return
	}

	if err := validateDateOfBirth(req.DateOfBirth); err != nil {
		badRequest(c, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		bindPassword(c, err)
		return
	}

	avatar, savedFile, err := h.readAvatar(c, req.Avatar)

	if err != nil {
		avatarError(c, err)
		return
	}

	req.Avatar = avatar

	user, err := h.repo.UserRegistration(req, string(hash))

	if err != nil {
		h.storage.Delete(savedFile)
		if errors.Is(err, repository.ErrLoginTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, withToken(h.present(user)))
}

// GetMe godoc
// @Summary      Текущий пользователь
// @Description  Возвращает профиль владельца переданного токена.
// @Tags         auth
// @Produce      json
// @Security     TokenAuth
// @Success      200  {object}  models.User
// @Failure      401  {object}  models.ErrorResponse
// @Router       /auth/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуеться авторизация"})
		return
	}

	token, ok := strings.CutPrefix(header, "Token ")

	if !ok || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка Authorization"})
		return
	}

	user, err := h.repo.GetUserByToken(token)

	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}

	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, h.present(user))
}

// Login godoc
// @Summary      Вход
// @Description  Проверяет логин/пароль, возвращает токен для дальнейшей авторизации.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Логин и пароль"
// @Success      200   {object}  models.UserWithToken
// @Failure      400   {object}  models.ValidationErrorResponse
// @Failure      401   {object}  models.ErrorResponse  "Неверный логин или пароль"
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		bindError(c, err)
		return
	}

	user, err := h.repo.GetByLogin(req.Login)

	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if err != nil {
		internalError(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	c.JSON(http.StatusOK, withToken(h.present(user)))
}
