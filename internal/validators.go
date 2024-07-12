package cache

import (
	"errors"
	"strings"
	"time"
)

func validateExpiry(t time.Time) error {
	if t.Before(time.Now()) {
		return errors.New("expiry cannot be in the past")
	}
	return nil
}

func validateKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key cannot be empty")
	}
	return nil
}

func validateValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("value cannot be empty")
	}
	return nil
}
