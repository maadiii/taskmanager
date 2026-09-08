package config

import (
	"log"
	"os"
	"strconv"
)

func getEnvString(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("%q env does not exists", key)
	}

	return env
}

func getEnvInt(key string) int {
	env := getEnvString(key)

	value, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf("invalid int env: %q: %v", key, err)
	}

	return value
}
