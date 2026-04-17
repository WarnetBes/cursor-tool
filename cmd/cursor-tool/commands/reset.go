package commands

import (
	"fmt"
	"os"

	"github.com/WarnetBes/cursor-tool/internal/backup"
	"github.com/WarnetBes/cursor-tool/internal/integrity"
	"github.com/WarnetBes/cursor-tool/internal/logger"
	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/WarnetBes/cursor-tool/internal/storage"
	"github.com/WarnetBes/cursor-tool/internal/uuid"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Cursor IDE Machine ID",
	Long:  `Reset the Cursor IDE Machine ID by generating new UUIDs and updating storage.json and Windows Registry (on Windows).`,
	RunE:  runReset,
}

func init() {
	resetCmd.Flags().BoolP("backup", "b", true, "Create backup before reset")
	resetCmd.Flags().BoolP("tui", "t", false, "Launch interactive TUI mode")
}

func runReset(cmd *cobra.Command, args []string) error {
	log := logger.New()

	backupFlag, _ := cmd.Flags().GetBool("backup")

	paths, err := platform.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get platform paths: %w", err)
	}

	if backupFlag {
		log.Info("Creating backup before reset...")
		if err := backup.Create(paths); err != nil {
			log.Warn("Failed to create backup: %v", err)
		}
	}

	newIDs, err := uuid.GenerateIDs()
	if err != nil {
		return fmt.Errorf("failed to generate new IDs: %w", err)
	}

	log.Info("Generated new Machine ID: %s", newIDs.MachineID)
	log.Info("Generated new Device ID: %s", newIDs.DeviceID)
	log.Info("Generated new Mac Machine ID: %s", newIDs.MacMachineID)

	if err := storage.Write(paths.StorageJSON, newIDs); err != nil {
		return fmt.Errorf("failed to write storage.json: %w", err)
	}

	if err := integrity.WriteHMAC(paths.StorageJSON); err != nil {
		log.Warn("Failed to write integrity HMAC: %v", err)
	}

	if err := platform.WriteRegistry(newIDs); err != nil {
		log.Warn("Failed to update registry (non-Windows or no permissions): %v", err)
	}

	log.Success("Machine ID reset successfully!")
	fmt.Fprintln(os.Stdout, "\nNew IDs:")
	fmt.Fprintf(os.Stdout, "  machineId:    %s\n", newIDs.MachineID)
	fmt.Fprintf(os.Stdout, "  devDeviceId:  %s\n", newIDs.DeviceID)
	fmt.Fprintf(os.Stdout, "  macMachineId: %s\n", newIDs.MacMachineID)

	return nil
}
