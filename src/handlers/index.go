package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"users-api-v1/src/media"
	"users-api-v1/src/models"
	"users-api-v1/src/repository"
)

var errInvalidAvatarURL = errors.New("Поля avatar принимает ссылку с http(s) началом или файл с изображением")
var minDateOfBirth = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

type UserHandler struct {
	repo    *repository.UserRepo
	storage *media.Storage
	baseURL string
}

func NewUserHandler(repo *repository.UserRepo, storage *media.Storage, baseURL string) *UserHandler {
	return &UserHandler{repo: repo, storage: storage, baseURL: baseURL}
}

func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func bindError(c *gin.Context, err error) {

	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Тело запроса слишком большое"})
		return
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fields := make([]models.FieldError, 0, len(verrs))

		for _, fe := range verrs {
			fields = append(fields, models.FieldError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Error:  "Ошибка валидации",
			Fields: fields,
		})
		return
	}

	badRequest(c, "Некорректный формат запроса")
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
		Error:  message,
		Fields: []models.FieldError{},
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "обязательное поле"
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("минимальная длина - %s символов", fe.Param())
		}
		return fmt.Sprintf("минимальное значение - %s", fe.Param())
	case "max":
		return fmt.Sprintf("Максимальное значение - %s", fe.Param())
	case "gt":
		return fmt.Sprintf("должно быть больше %s", fe.Param())
	case "eqfield":
		return fmt.Sprintf("должно совпадать с полем %s", strings.ToLower(fe.Param()))
	default:
		return fmt.Sprintf("не проходит проверку '%s'", fe.Tag())
	}
}

func bindPassword(c *gin.Context, err error) {

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

}

func avatarError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, media.ErrTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
	case errors.Is(err, media.ErrUnsupportedType):
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
	case errors.Is(err, errInvalidAvatarURL):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		internalError(c, err)
	}
}

func internalError(c *gin.Context, err error) {
	log.Printf("Internal error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

func (h *UserHandler) readAvatar(c *gin.Context, fromJSON *string) (*string, string, error) {
	if c.ContentType() != "multipart/form-data" {
		if fromJSON != nil && !isHTTPURL(*fromJSON) {
			return nil, "", errInvalidAvatarURL
		}
		return fromJSON, "", nil
	}

	fileHeader, err := c.FormFile("avatar")
	if err == nil {
		path, err := h.storage.SaveImage(fileHeader)
		if err != nil {
			return nil, "", err
		}
		return &path, path, nil
	}
	if !errors.Is(err, http.ErrMissingFile) {
		return nil, "", err
	}

	if s, ok := c.GetPostForm("avatar"); ok {
		if !isHTTPURL(s) {
			return nil, "", errInvalidAvatarURL
		}
		return &s, "", nil
	}
	return nil, "", nil
}

func (handler *UserHandler) present(user *models.User) *models.User {
	if user.Avatar != nil && media.IsLocal(*user.Avatar) {
		full := handler.baseURL + *user.Avatar
		user.Avatar = &full
	}
	return user
}

func withToken(user *models.User) gin.H {
	return gin.H{
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"login":         user.Login,
		"avatar":        user.Avatar,
		"boss_id":       user.BossID,
		"date_of_birth": user.DateOfBirth,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"token":         user.Token,
	}
}

func validateDateOfBirth(s string) error {
	t, err := time.Parse("2006-01-02", s)

	if err != nil {
		return errors.New("Поля date_of_birth должно быть датой в формате YYYY-MM-DD")
	}

	if t.Before(minDateOfBirth) {
		return errors.New("Поля date_of_birth не может быть ранее 1900 года")
	}

	if t.After(time.Now()) {
		return errors.New("Поле date_of_birth не может быть в будущем")
	}

	return nil
}
