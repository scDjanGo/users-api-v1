package features

import (
	"log"
	"os"
	"strconv"
)

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetEnvInt(key string, def int64) int64 {
	v := os.Getenv(key)

	if v == "" {
		return def
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Printf("%s=%q не число, используем число по умолчании %d: %v", key, v, def, err)
		return def
	}
	return n
}
