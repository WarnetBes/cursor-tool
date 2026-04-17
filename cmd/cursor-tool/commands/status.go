package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/WarnetBes/cursor-tool/internal/platform"
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

	paths, err := platform.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get platform paths: %w", err)
	}

	data, err := os.ReadFile(paths.StorageJSON)
	if err != nil {
		return fmt.Errorf("failed to read storage.json: %w", err)
	}

	var storage map[string]interface{}
	if err := json.Unmarshal(data, &storage); err != nil {
		return fmt.Errorf("failed to parse storage.json: %w", err)
	}

	keys := []string{"telemetry.machineId", "telemetry.devDeviceId", "telemetry.macMachineId", "telemetry.sqmId"}

	if jsonOutput {
		result := make(map[string]interface{})
		for _, k := range keys {
			if v, ok := storage[k]; ok {
				result[k] = v
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Println("Current Cursor IDE Machine IDs:")
	fmt.Println("================================")
	for _, k := range keys {
		if v, ok := storage[k]; ok {
			fmt.Fprintf(os.Stdout, "  %-30s %v\n", k+":", v)
		} else {
			fmt.Fprintf(os.Stdout, "  %-30s (not set)\n", k+":")
		}
	}

	return nil
}
