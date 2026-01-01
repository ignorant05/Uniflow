package credentials

import (
	"fmt"

	constants "github.com/ignorant05/Uniflow/internal/constants/credentials"
	"github.com/zalando/go-keyring"
)

// Store stores credentials
//
// Parameters:
//   - key: key string
//   - val: value string
//
// Example usage:
//
//	err := Store("name", "ignorant05")
func Store(key, val string) error {
	if err := keyring.Set(constants.SERVICE, key, val); err != nil {
		return fmt.Errorf("<?> Error: Failed to store credentials in keyring\nError: %w", err)
	}

	return nil
}

// gets value from key if exists
//
// Parameters:
//   - key: key string
//
// Example usage:
//
//	val, err := Get("name")
func Get(key string) (string, error) {
	val, err := keyring.Get(constants.SERVICE, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("<?> Error: Credential '%s' not found in keyring", key)
		}
		return "", fmt.Errorf("<?> Error: Failed to retrieve credential '%s' from keyring", key)
	}

	return val, nil
}

// Delete deletes value from key if exists
//
// Parameters:
//   - key: key string
//
// Example usage:
//
//	err := Delete("name")
func Delete(key string) error {
	err := keyring.Delete(constants.SERVICE, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("<?> Error: Credential '%s' not found in keyring", key)
		}
		return fmt.Errorf("<?> Error: Failed to delete credential '%s' from keyring", key)
	}

	return nil
}

// List lists key, value pairs
//
// Parameters:
//   - key: key string
//
// Example usage:
//
//	list, err := List()
func List() ([]string, error) {
	var availableKeys []string
	for _, key := range constants.CommonKeys {
		if _, err := Get(key); err == nil {
			availableKeys = append(availableKeys, key)
		}
	}

	return availableKeys, nil
}

// StoreForPlatform stores token with formatted key
//
// Parameters:
//   - profile: profile name
//   - platform: platform name
//   - token: token string
//
// Example usage:
//
//	err := StoreForPlatform("default", "ignorant05", "gibbris as token")
func StoreForPlatform(profile, platform, token string) error {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Store(key, token)
}

// GetForPlatform gets token with formatted key
//
// Parameters:
//   - profile: profile name
//   - platform: platform name
//
// Example usage:
//
//	err := GetForPlatform("default", "ignorant05")
func GetForPlatform(profile, platform string) (string, error) {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Get(key)
}

// DeleteForPlatform deteles token with formatted key
//
// Parameters:
//   - profile: profile name
//   - platform: platform name
//
// Example usage:
//
//	err := DeleteForPlatform("default", "ignorant05")

func DeleteForPlatform(profile, platform string) error {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Delete(key)
}
