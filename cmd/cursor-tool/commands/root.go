package commands

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "cursor-tool",
	Short:   "Cursor IDE Machine ID reset utility",
	Long:    `cursor-tool resets Cursor IDE telemetry IDs (machineId, devDeviceId, macMachineId) across Windows, macOS and Linux.`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(statusCmd)
}
