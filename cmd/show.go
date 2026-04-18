package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sapihav/clef-cli/internal/settings"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Print .claude/settings.local.json content",
	Long:  "Print the current content of .claude/settings.local.json as pretty JSON.",
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	path := settings.FilePath(".")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("no settings file found at %s\n", path)
			return nil
		}
		return userError("read settings: " + err.Error())
	}

	// Re-marshal for pretty printing regardless of original formatting.
	var m interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return userError("parse settings: " + err.Error())
	}
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return userError("format settings: " + err.Error())
	}
	fmt.Println(string(out))
	return nil
}
