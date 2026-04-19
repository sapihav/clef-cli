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
	effortCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Print what would be written without writing")
	rootCmd.AddCommand(effortCmd)
}

func runEffort(cmd *cobra.Command, args []string) error {
	v := args[0]
	if !validEfforts[v] {
		return userError(fmt.Sprintf("invalid effort %q (want low|medium|high|xhigh)", v))
	}

	data, err := settings.Load(".")
	if err != nil {
		return sysError("read settings: " + err.Error())
	}
	data["effortLevel"] = v

	if flagDryRun {
		writeJSON(data)
		return nil
	}

	if err := settings.Save(".", data); err != nil {
		return sysError("write settings: " + err.Error())
	}

	if flagJSON {
		writeJSON(map[string]interface{}{"effortLevel": v, "file": settings.FilePath(".")})
	} else {
		fmt.Printf("effort=%s → %s\n", v, settings.FilePath("."))
	}
	return nil
}
