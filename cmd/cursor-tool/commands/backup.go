package commands

import (
	"fmt"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/logger"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of Cursor IDE storage",
	Long:  `Create a backup of the Cursor IDE storage.json file.`,
	RunE:  runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return fmt.Errorf("failed to get storage path: %w", err)
	}

	mgr := backup.New(5)
	logger.Info("Creating backup...")

	backupPath, err := mgr.Create(storagePath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	logger.Success("Backup created: %s", backupPath)
	return nil
}
