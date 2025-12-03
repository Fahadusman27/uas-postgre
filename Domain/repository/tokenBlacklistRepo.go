package repository

import (
	"sync"
	"time"
)

// TokenBlacklistRepository interface untuk abstraksi data access
type TokenBlacklistRepository interface {
	Add(token string, expiresAt time.Time) error
	Exists(token string) (bool, error)
	Remove(token string) error
	Cleanup() error
}

// InMemoryTokenBlacklist implementasi repository dengan in-memory storage
// Untuk production, buat implementasi dengan Redis atau database
type InMemoryTokenBlacklist struct {
	storage map[string]time.Time
	mu      sync.RWMutex
}

// NewInMemoryTokenBlacklist membuat instance repository baru
func NewInMemoryTokenBlacklist() *InMemoryTokenBlacklist {
	return &InMemoryTokenBlacklist{
		storage: make(map[string]time.Time),
	}
}

// Add menambahkan token ke blacklist
func (r *InMemoryTokenBlacklist) Add(token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[token] = expiresAt
	return nil
}

// Exists mengecek apakah token ada di blacklist dan masih valid
func (r *InMemoryTokenBlacklist) Exists(token string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expiresAt, exists := r.storage[token]
	if !exists {
		return false, nil
	}

	// Token dianggap tidak ada jika sudah expired
	if time.Now().After(expiresAt) {
		return false, nil
	}

	return true, nil
}

// Remove menghapus token dari blacklist
func (r *InMemoryTokenBlacklist) Remove(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, token)
	return nil
}

// Cleanup menghapus semua token yang sudah expired
func (r *InMemoryTokenBlacklist) Cleanup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for token, expiresAt := range r.storage {
		if now.After(expiresAt) {
			delete(r.storage, token)
		}
	}
	return nil
}

// Global instance untuk backward compatibility
// TODO: Refactor untuk menggunakan dependency injection
var defaultBlacklist = NewInMemoryTokenBlacklist()

// AddTokenToBlacklist - wrapper function untuk backward compatibility
func AddTokenToBlacklist(token string, expiresAt time.Time) error {
	return defaultBlacklist.Add(token, expiresAt)
}

// IsTokenBlacklisted - wrapper function untuk backward compatibility
func IsTokenBlacklisted(token string) (bool, error) {
	return defaultBlacklist.Exists(token)
}

// CleanupExpiredTokens - wrapper function untuk backward compatibility
func CleanupExpiredTokens() error {
	return defaultBlacklist.Cleanup()
}
