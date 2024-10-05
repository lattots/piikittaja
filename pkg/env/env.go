package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// LoadVariables loads environment variables from /assets/.env
func LoadVariables() error {
	projectRoot, err := GetProjectRoot()
	envPath := filepath.Join(projectRoot, "assets", ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}

// GetProjectRoot returns the full path to the project root folder.
// This is useful when trying to access certain files in the project
func GetProjectRoot() (string, error) {
	rawPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	projectRoot := strings.TrimSuffix(filepath.Dir(rawPath), "/bin")
	return projectRoot, nil
}
