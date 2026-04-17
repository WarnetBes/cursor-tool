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
	Short: "Restore Cursor IDE storage from backup",
	Long:  `Restore the Cursor IDE storage.json file from a previously created backup.`,
	RunE:  runRestore,
}

func init() {
	restoreCmd.Flags().StringP("file", "f", "", "Backup file to restore from (default: latest backup)")
}

func runRestore(cmd *cobra.Command, args []string) error {
	log := logger.New()

	backupFile, _ := cmd.Flags().GetString("file")

	paths, err := platform.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get platform paths: %w", err)
	}

	if backupFile == "" {
		log.Info("No backup file specified, using latest backup...")
		backupFile, err = backup.LatestBackup(paths)
		if err != nil {
			return fmt.Errorf("failed to find latest backup: %w", err)
		}
	}

	log.Info("Restoring from: %s", backupFile)
	if err := backup.Restore(backupFile, paths); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	log.Success("Backup restored successfully!")
	return nil
}
