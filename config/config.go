package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// Config holds the application configuration.
type Config struct {
	DefaultTarget string
	DefaultPorts  string
	ScanType      models.ScanType
	Timeout       int
	Workers       int
	OutputFormat  string
	OutputFile    string
	Debug         bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DefaultTarget: getEnv("PORTHUNTER_TARGET", "192.168.1.130"),
		DefaultPorts:  getEnv("PORTHUNTER_PORTS", "22,80,443,3306,3389,8080,8443"),
		ScanType:      models.ScanType(getEnv("PORTHUNTER_SCAN_TYPE", string(models.ScanTypeTCPConnect))),
		Timeout:       getEnvInt("PORTHUNTER_TIMEOUT", 3000),
		Workers:       getEnvInt("PORTHUNTER_WORKERS", 100),
		OutputFormat:  getEnv("PORTHUNTER_OUTPUT", "table"),
		OutputFile:    getEnv("PORTHUNTER_OUTPUT_FILE", ""),
		Debug:         getEnvBool("PORTHUNTER_DEBUG", false),
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Workers <= 0 {
		return fmt.Errorf("workers must be a positive number")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be a positive number")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}
