package commands

import (
	"fmt"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/logger"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Cursor IDE storage from latest backup",
	Long:  `Restore the Cursor IDE storage.json file from the latest backup.`,
	RunE:  runRestore,
}

func runRestore(cmd *cobra.Command, args []string) error {
	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return fmt.Errorf("failed to get storage path: %w", err)
	}

	mgr := backup.New(5)
	logger.Info("Restoring from latest backup...")

	restoredFrom, err := mgr.Restore(storagePath)
	if err != nil {
		return fmt.Errorf("failed to restore: %w", err)
	}

	logger.Success("Restored from: %s", restoredFrom)
	return nil
}
