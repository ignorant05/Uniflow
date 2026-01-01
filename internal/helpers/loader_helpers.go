package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	constants "github.com/ignorant05/Uniflow/internal/constants/config"
)

// GetConfigDir returns full config dir path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("<?> Error: Failed to get user home directory...\n<?> Error: %w", err)
	}

	return filepath.Join(homeDir, constants.DEFAULT_CONFIG_DIR_PATH), nil
}

// GetConfigPath full config file path
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, constants.DEFAULT_CONFIG_FILE_PATH), nil

}
