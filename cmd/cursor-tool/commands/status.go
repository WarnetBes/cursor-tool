package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/WarnetBes/cursor-tool/internal/platform"
	"github.com/WarnetBes/cursor-tool/internal/storage"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current Cursor IDE Machine ID status",
	Long:  `Display the current Machine ID values stored in Cursor IDE storage.json.`,
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func runStatus(cmd *cobra.Command, args []string) error {
	jsonOutput, _ := cmd.Flags().GetBool("json")

	storagePath, err := platform.GetStoragePath()
	if err != nil {
		return fmt.Errorf("failed to get storage path: %w", err)
	}

	ids, err := storage.ReadCurrentIDs(storagePath)
	if err != nil {
		return fmt.Errorf("failed to read current IDs: %w", err)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(ids)
	}

	fmt.Println("Current Cursor IDE Machine IDs:")
	fmt.Println("================================")
	for _, k := range storage.TelemetryFields {
		v := ids[k]
		if v == "" {
			v = "(not set)"
		}
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", k+":", v)
	}
	return nil
}
