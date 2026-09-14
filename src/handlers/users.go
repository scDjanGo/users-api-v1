package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"users-api-v1/src/models"
	"users-api-v1/src/repository"
)

// Create godoc
// @Summary      Создать пользователя
// @Description  Аватар — http(s)-ссылка в JSON, либо файл в multipart/form-data.
// @Tags         users
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        body  body      models.CreateUserRequest  true  "Данные пользователя"
// @Success      201   {object}  models.UserWithToken
// @Failure      400   {object}  models.ValidationErrorResponse
// @Failure      409   {object}  models.ErrorResponse  "Логин уже занят"
// @Failure      413   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
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

	user, err := h.repo.Create(req, string(hash))

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

// Update godoc
// @Summary      Обновить пользователя
// @Description  Частичное обновление — присылайте только те поля, которые нужно изменить.
// @Tags         users
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      int                        true  "ID пользователя"
// @Param        body  body      models.UpdatedUserRequest  true  "Поля для обновления"
// @Success      200   {object}  models.User
// @Failure      400   {object}  models.ValidationErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Failure      413   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		badRequest(c, "Невалидный id пользователя")
		return
	}

	var req models.UpdatedUserRequest
	if err := c.ShouldBind(&req); err != nil {
		bindError(c, err)
		return
	}

	if req.DateOfBirth != nil {
		if err := validateDateOfBirth(*req.DateOfBirth); err != nil {
			badRequest(c, err.Error())
			return
		}
	}

	existing, err := h.repo.GetByID(id)

	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		internalError(c, err)
		return
	}

	avatar, savedFile, err := h.readAvatar(c, req.Avatar)
	if err != nil {
		avatarError(c, err)
		return
	}
	req.Avatar = avatar

	if req.FirstName == nil && req.LastName == nil && req.Avatar == nil && req.BossID == nil && req.DateOfBirth == nil {
		badRequest(c, "Ничего нету для обновления")
		return
	}

	user, err := h.repo.Update(id, req)

	if err != nil {
		h.storage.Delete(savedFile)
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		internalError(c, err)
		return
	}

	if req.Avatar != nil && existing.Avatar != nil && *existing.Avatar != *req.Avatar {
		h.storage.Delete(*existing.Avatar)
	}

	c.JSON(http.StatusOK, h.present(user))
}

// ReadById godoc
// @Summary      Получить пользователя по id
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) ReadById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не валидный параметр"})
		return
	}

	user, err := h.repo.GetByID(id)

	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь с таким id не найден"})
		return
	}

	if err != nil {
		internalError(c, err)
	}

	c.JSON(http.StatusOK, user)
}

// ReadAll godoc
// @Summary      Список пользователей
// @Description  Поддерживает фильтрацию и сортировку через query-параметры.
// @Tags         users
// @Produce      json
// @Param        sort_by     query     string  false  "Поле сортировки"          Enums(id, date_of_birth, boss_id, created_at)
// @Param        order       query     string  false  "Направление сортировки"   Enums(asc, desc)
// @Param        has_boss    query     bool    false  "Только с начальником (true) / без (false)"
// @Param        has_avatar  query     bool    false  "Только с аватаром (true) / без (false)"
// @Param        boss_id     query     int     false  "Точный id начальника"
// @Success      200  {array}   models.User
// @Router       /users [get]
func (h *UserHandler) ReadAll(c *gin.Context) {
	var q models.ListUserQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		bindError(c, err)
		return
	}

	users, err := h.repo.GetAll(q)

	if err != nil {
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)

}

// Delete godoc
// @Summary      Удалить пользователя
// @Description  Подчинённые не удаляются — у них обнуляется boss_id.
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  models.MessageResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не правильный формат id"})
		return
	}

	err = h.repo.Delete(id)

	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Пользователь успешно удалён"})

}
