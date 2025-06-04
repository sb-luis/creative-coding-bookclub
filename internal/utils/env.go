package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnvFile loads environment variables from a .env file at the project root.
// This function is intended for development use only.
// It searches for the .env file starting from the current working directory
// and walking up the directory tree until it finds the project root.
func LoadEnvFile() error {
	envPath, err := findEnvFile()
	if err != nil {
		return fmt.Errorf("failed to find .env file: %w", err)
	}

	return loadEnvFromFile(envPath)
}

// findEnvFile searches for the .env file starting from the current directory
// and walking up the directory tree until it finds it or reaches the root.
func findEnvFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached the root directory
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf(".env file not found")
}

// loadEnvFromFile reads the .env file and sets environment variables.
func loadEnvFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format at line %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		// Only set the environment variable if it's not already set
		// This allows existing environment variables to take precedence
		if os.Getenv(key) == "" {
			err := os.Setenv(key, value)
			if err != nil {
				return fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}
