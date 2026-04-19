package cmd

import (
	"fmt"

	"github.com/sapihav/clef-cli/internal/settings"
	"github.com/spf13/cobra"
)

var validEfforts = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
	"xhigh":  true,
}

var effortCmd = &cobra.Command{
	Use:   "effort <value>",
	Short: "Set effort level in .claude/settings.local.json",
	Args:  cobra.ExactArgs(1),
	RunE:  runEffort,
}

func init() {
	rootCmd.AddCommand(effortCmd)
}

func runEffort(cmd *cobra.Command, args []string) error {
	v := args[0]
	if !validEfforts[v] {
		return userError(fmt.Sprintf("invalid effort %q (want low|medium|high|xhigh)", v))
	}

	data, err := settings.Load(".")
	if err != nil {
		return userError("read settings: " + err.Error())
	}
	data["effortLevel"] = v
	if err := settings.Save(".", data); err != nil {
		return userError("write settings: " + err.Error())
	}

	fmt.Printf("effort=%s → %s\n", v, settings.FilePath("."))
	return nil
}
