package cmd

import (
	"fmt"
	"strings"

	"github.com/sapihav/clef-cli/internal/settings"
	"github.com/spf13/cobra"
)

var shorthandModels = map[string]string{
	"opus4.7":   "claude-opus-4-7",
	"opus4.6":   "claude-opus-4-6",
	"sonnet4.6": "claude-sonnet-4-6",
	"haiku4.5":  "claude-haiku-4-5",
}

var familyAliases = map[string]bool{
	"sonnet": true,
	"opus":   true,
	"haiku":  true,
}

// resolveModel maps user input to the stored model value.
func resolveModel(input string) (string, error) {
	if familyAliases[input] {
		return input, nil
	}
	if full, ok := shorthandModels[input]; ok {
		return full, nil
	}
	if strings.HasPrefix(input, "claude-") {
		return input, nil
	}
	return "", fmt.Errorf("invalid model %q (want sonnet|opus|haiku, shorthand like opus4.7, or full claude-* id)", input)
}

var modelCmd = &cobra.Command{
	Use:   "model <value>",
	Short: "Set Claude model in .claude/settings.local.json",
	Args:  cobra.ExactArgs(1),
	RunE:  runModel,
}

func init() {
	rootCmd.AddCommand(modelCmd)
}

func runModel(cmd *cobra.Command, args []string) error {
	resolved, err := resolveModel(args[0])
	if err != nil {
		return userError(err.Error())
	}

	data, err := settings.Load(".")
	if err != nil {
		return userError("read settings: " + err.Error())
	}
	data["model"] = resolved
	if err := settings.Save(".", data); err != nil {
		return userError("write settings: " + err.Error())
	}

	fmt.Printf("model=%s → %s\n", resolved, settings.FilePath("."))
	return nil
}
