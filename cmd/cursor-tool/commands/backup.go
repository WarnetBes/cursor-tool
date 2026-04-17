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
	Long:  `Create a backup of the Cursor IDE storage.json file and registry settings.`,
	RunE:  runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	log := logger.New()

	paths, err := platform.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get platform paths: %w", err)
	}

	log.Info("Creating backup...")
	backupPath, err := backup.CreateWithPath(paths)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	log.Success("Backup created: %s", backupPath)
	return nil
}
