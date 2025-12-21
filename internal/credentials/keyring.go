package credentials

import (
	"fmt"

	constants "github.com/ignorant05/Uniflow/internal/constants/credentials"
	"github.com/zalando/go-keyring"
)

func Store(key, val string) error {
	if err := keyring.Set(constants.SERVICE, key, val); err != nil {
		return fmt.Errorf("<?> Error: Failed to store credentials in keyring\nError: %w\n", err)
	}

	return nil
}

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

func List() ([]string, error) {
	var availableKeys []string
	for _, key := range constants.CommonKeys {
		if _, err := Get(key); err == nil {
			availableKeys = append(availableKeys, key)
		}
	}

	return availableKeys, nil
}

func StoreForPlatform(profile, platform, token string) error {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Store(key, token)
}

func GetForPlatform(profile, platform string) (string, error) {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Get(key)
}

func DeleteForPlatform(profile, platform string) error {
	key := fmt.Sprintf("%s.%s.token", profile, platform)
	return Delete(key)
}
