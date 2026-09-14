package seed

import (
	"log"
	"time"

	"users-api-v1/src/repository"
)

func StartPeriodic(repo *repository.UserRepo, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := Seed(repo); err != nil {
				log.Printf("seed: Плановая переустановка данных не удалась: %v", err)
			}
		}
	}()
}
