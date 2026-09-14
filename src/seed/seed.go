package seed

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"users-api-v1/src/models"
	"users-api-v1/src/repository"
)

const (
	usersCount      = 10
	defaultPassword = "Student123"
)

var firstNames = []string{"Алия", "Бекзат", "Динара", "Ерлан", "Жанна", "Канат", "Мадина", "Нурлан", "Сауле", "Тимур"}
var lastNames = []string{"Ахметова", "Бектасов", "Валиева", "Габитов", "Досжанова", "Ержанов", "Жумабекова", "Каримов", "Смагулова", "Токтаров"}

func Seed(repo *repository.UserRepo) error {
	if err := repo.DeleteAll(); err != nil {
		return fmt.Errorf("Очистка таблицы users: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Хеширование пароля по умолчанию: %w", err)
	}

	passwordHash := string(hash)

	for i := 1; i <= usersCount; i++ {
		dob := time.Date(1970+rand.Intn(50), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), 0, 0, 0, 0, time.UTC).Format("2006-01-02")

		req := models.CreateUserRequest{
			FirstName:   firstNames[rand.Intn(len(firstNames))],
			LastName:    lastNames[rand.Intn(len(lastNames))],
			Login:       fmt.Sprintf("student%d", i),
			Password:    defaultPassword,
			DateOfBirth: dob,
		}

		if rand.Intn(2) == 0 {
			avatar := fmt.Sprintf("https://i.pravatar.cc/150?img=%d", rand.Intn(70)+1)
			req.Avatar = &avatar
		}

		if _, err := repo.Create(req, passwordHash); err != nil {
			return fmt.Errorf("Создание пользователя %s: %w", req.Login, err)
		}
	}

	log.Printf("seed: таблица users очищена, создано %d демо-пользователей (login: student1...%d), пароль: %s", usersCount, usersCount, defaultPassword)
	return nil
}
