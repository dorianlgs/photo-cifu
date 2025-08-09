package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	WorkflowDB struct {
		Name string
	}
	Gallery struct {
		MaxFileSize int64 // in bytes
		MaxImages   int
	}
	Workflow struct {
		DefaultTimeout int // in seconds
	}
}

// New creates a new configuration with defaults and environment overrides
func New() *Config {
	cfg := &Config{}
	
	// Set defaults
	cfg.WorkflowDB.Name = "workflow.db"
	cfg.Gallery.MaxFileSize = 100 * 1024 * 1024 // 100MB
	cfg.Gallery.MaxImages = 100
	cfg.Workflow.DefaultTimeout = 300 // 5 minutes

	// Override with environment variables if present
	if dbName := os.Getenv("WORKFLOW_DB_NAME"); dbName != "" {
		cfg.WorkflowDB.Name = dbName
	}

	if maxSize := os.Getenv("GALLERY_MAX_FILE_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			cfg.Gallery.MaxFileSize = size
		}
	}

	if maxImages := os.Getenv("GALLERY_MAX_IMAGES"); maxImages != "" {
		if count, err := strconv.Atoi(maxImages); err == nil {
			cfg.Gallery.MaxImages = count
		}
	}

	if timeout := os.Getenv("WORKFLOW_DEFAULT_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Workflow.DefaultTimeout = t
		}
	}

	return cfg
}