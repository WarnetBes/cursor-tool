package commands

import (
	"fmt"
	"os"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/integrity"
	"github.com/WarnetBes/cursor-tool/internal/logger"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/WarnetBes/cursor-tool/internal/storage"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Cursor IDE Machine ID",
	Long:  `Reset the Cursor IDE Machine ID by generating new UUIDs and updating storage.json.`,
	RunE:  runReset,
}

func init() {
	resetCmd.Flags().BoolP("no-backup", "n", false, "Skip backup before reset")
}

func runReset(cmd *cobra.Command, args []string) error {
	noBackup, _ := cmd.Flags().GetBool("no-backup")

	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return fmt.Errorf("failed to get storage path: %w", err)
	}

	mgr := backup.New(5)
	if noBackup {
		mgr = backup.New(0)
	}

	result, err := storage.ModifyStorageIDs(storagePath, mgr)
	if err != nil {
		return fmt.Errorf("failed to reset IDs: %w", err)
	}

	if err := integrity.WriteHMAC(storagePath); err != nil {
		logger.Warn("Failed to write HMAC: %v", err)
	}

	logger.Success("Machine ID reset successfully!")
	if result.BackupPath != "" {
		fmt.Fprintf(os.Stdout, "Backup: %s\n", result.BackupPath)
	}
	fmt.Fprintln(os.Stdout, "\nNew IDs:")
	for k, v := range result.After {
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", k+":", v)
	}
	return nil
}
