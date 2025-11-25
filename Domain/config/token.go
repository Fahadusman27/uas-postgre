package config

import (
	"os"
	"strconv"
	"time"
)

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_me"
	}
	return secret
}

func GetJWTExpiry() time.Duration {
	h := os.Getenv("JWT_EXPIRE_HOURS")
	if h == "" {
		return 24 * time.Minute
	}
	hi, err := strconv.Atoi(h)
	if err != nil {
		return 24 * time.Minute
	}
	return time.Duration(hi) * time.Minute
}