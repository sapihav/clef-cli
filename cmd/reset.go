package cmd

import (
	"fmt"
	"os"

	"github.com/sapihav/clef-cli/internal/settings"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove model and effortLevel from .claude/settings.local.json",
	Long: `Remove model and effortLevel keys from .claude/settings.local.json.
All other keys are preserved. If the file does not exist, this is a no-op.`,
	RunE: runReset,
}

func init() {
	resetCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Print what would be written without writing")
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	path := settings.FilePath(".")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if flagJSON {
			writeJSON(map[string]interface{}{"removed": []string{}, "file": path})
		} else {
			fmt.Printf("no settings file at %s — nothing to reset\n", path)
		}
		return nil
	}

	data, err := settings.Load(".")
	if err != nil {
		return sysError("read settings: " + err.Error())
	}

	delete(data, "model")
	delete(data, "effortLevel")

	if flagDryRun {
		writeJSON(data)
		return nil
	}

	if err := settings.Save(".", data); err != nil {
		return sysError("write settings: " + err.Error())
	}

	if flagJSON {
		writeJSON(map[string]interface{}{"removed": []string{"model", "effortLevel"}, "file": path})
	} else {
		fmt.Printf("model and effortLevel removed from %s\n", path)
	}
	return nil
}
