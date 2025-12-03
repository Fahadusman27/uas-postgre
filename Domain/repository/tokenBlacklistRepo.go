package repository

import (
	"sync"
	"time"
)

// In-memory token blacklist (alternatif tanpa tabel database)
// Untuk production, gunakan Redis
var (
	tokenBlacklist = make(map[string]time.Time)
	blacklistMutex sync.RWMutex
)

func AddTokenToBlacklist(token string, expiresAt time.Time) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	tokenBlacklist[token] = expiresAt
	return nil
}

func IsTokenBlacklisted(token string) (bool, error) {
	blacklistMutex.RLock()
	defer blacklistMutex.RUnlock()

	expiresAt, exists := tokenBlacklist[token]
	if !exists {
		return false, nil
	}

	// Cek apakah token sudah expired
	if time.Now().After(expiresAt) {
		return false, nil
	}

	return true, nil
}

func CleanupExpiredTokens() error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	now := time.Now()
	for token, expiresAt := range tokenBlacklist {
		if now.After(expiresAt) {
			delete(tokenBlacklist, token)
		}
	}
	return nil
}
