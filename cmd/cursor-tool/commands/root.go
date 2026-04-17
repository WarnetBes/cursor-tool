package commands

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

var RootCmd = &cobra.Command{
	Use:     "cursor-tool",
	Short:   "Cursor IDE Machine ID reset utility",
	Long:    `cursor-tool resets Cursor IDE telemetry IDs (machineId, devDeviceId, macMachineId) across Windows, macOS and Linux.`,
	Version: Version,
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(resetCmd)
	RootCmd.AddCommand(backupCmd)
	RootCmd.AddCommand(restoreCmd)
	RootCmd.AddCommand(statusCmd)
}
