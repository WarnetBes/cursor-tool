package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetStoragePath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(appdata,
			"Cursor", "User", "globalStorage", "storage.json"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home,
			"Library", "Application Support",
			"Cursor", "User", "globalStorage", "storage.json"), nil
	default: // linux
		cfg := os.Getenv("XDG_CONFIG_HOME")
		if cfg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			cfg = filepath.Join(home, ".config")
		}
		return filepath.Join(cfg,
			"Cursor", "User", "globalStorage", "storage.json"), nil
	}
}

func GetOS() string   { return runtime.GOOS }
func GetArch() string { return runtime.GOARCH }
