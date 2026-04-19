package cmd

import (
	"fmt"
	"strings"

	"github.com/sapihav/clef-cli/internal/settings"
	"github.com/spf13/cobra"
)

var (
	flagModel  string
	flagEffort string
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set model and/or effort level in .claude/settings.local.json",
	Long: `Set model and/or effort level in .claude/settings.local.json.

At least one flag must be provided. Existing keys are preserved.

Valid models:  sonnet, opus, haiku, opus4.7, sonnet4.6, haiku4.5, or any claude-* id
Valid efforts: low, medium, high, xhigh

Example:
  clef set --model opus4.7 --effort xhigh
  clef set --effort low`,
	RunE: runSet,
}

func init() {
	setCmd.Flags().StringVar(&flagModel, "model", "", "Claude model: sonnet|opus|haiku")
	setCmd.Flags().StringVar(&flagEffort, "effort", "", "Effort level: low|medium|high|xhigh")
	rootCmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	if flagModel == "" && flagEffort == "" {
		return userError("at least one of --model or --effort must be provided")
	}
	if flagModel != "" {
		if _, err := resolveModel(flagModel); err != nil {
			return userError("--model: " + err.Error())
		}
	}
	if flagEffort != "" && !validEfforts[flagEffort] {
		return userError(fmt.Sprintf("invalid --effort %q (want low|medium|high|xhigh)", flagEffort))
	}

	data, err := settings.Load(".")
	if err != nil {
		return userError("read settings: " + err.Error())
	}

	if flagModel != "" {
		resolved, _ := resolveModel(flagModel)
		data["model"] = resolved
	}
	if flagEffort != "" {
		data["effortLevel"] = flagEffort
	}

	if err := settings.Save(".", data); err != nil {
		return userError("write settings: " + err.Error())
	}

	var parts []string
	if flagModel != "" {
		parts = append(parts, "model="+flagModel)
	}
	if flagEffort != "" {
		parts = append(parts, "effort="+flagEffort)
	}
	fmt.Printf("%s → %s\n", strings.Join(parts, " "), settings.FilePath("."))
	return nil
}
