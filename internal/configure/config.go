package configure

import (
	// Keeping fmt import for when we uncomment validation

	_ "fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jcc333/jkm/internal/log"
	"github.com/joho/godotenv"
)

// Global configuration for the application.
type Config struct {
	// IMAP server host.
	IMAPServer string

	// IMAP server port (993 by default).
	IMAPPort int

	// SMTP server host.
	SMTPServer string

	// SMTP server port (587 by default).
	SMTPPort int

	// The user's email address.
	EmailAddress string

	// The user's IMAP password.
	IMAPPassword string

	// The user's SMTP password.
	SMTPPassword string
}

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	log.Info("load configuration")
	home, err := os.UserHomeDir()
	if err == nil {
		_ = godotenv.Load(filepath.Join(home, ".env"))
		// An extra logger init in case JKM_LOGGING is set in the .env file
		isLogging := os.Getenv("JKM_LOGGING")
		log.Init(isLogging != "false" && isLogging != "")

	}
	_ = godotenv.Load(".env")
	cfg := &Config{
		IMAPPort: 993,
		SMTPPort: 587,
	}
	if val := os.Getenv("JKM_EMAIL"); val != "" {
		cfg.EmailAddress = val
	}
	if val := os.Getenv("JKM_SMTP_SERVER"); val != "" {
		cfg.SMTPServer = val
	}
	if val := os.Getenv("JKM_SMTP_PORT"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			log.Errorf("JKM_SMTP_PORT error: '%v'", err)
			return nil, err
		}
		cfg.SMTPPort = n
	}
	if val := os.Getenv("JKM_SMTP_PASSWORD"); val != "" {
		cfg.SMTPPassword = val
	}
	if val := os.Getenv("JKM_IMAP_SERVER"); val != "" {
		cfg.IMAPServer = val
	}
	if val := os.Getenv("JKM_IMAP_PORT"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			log.Errorf("JKM_IMAP_PORT error: '%v'", err)
			return nil, err
		}
		cfg.IMAPPort = n
	}
	if val := os.Getenv("JKM_IMAP_PASSWORD"); val != "" {
		cfg.IMAPPassword = val
	}
	return cfg, nil
}
