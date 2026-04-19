package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flagInfo struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}

type commandInfo struct {
	Name        string        `json:"name"`
	Use         string        `json:"use"`
	Description string        `json:"description"`
	Flags       []flagInfo    `json:"flags,omitempty"`
	Commands    []commandInfo `json:"commands,omitempty"`
}

type schemaResult struct {
	CLI         string        `json:"cli"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	GlobalFlags []flagInfo    `json:"global_flags,omitempty"`
	Commands    []commandInfo `json:"commands"`
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print the full command tree as JSON",
	RunE:  runSchema,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	writeJSON(schemaResult{
		CLI:         "clef",
		Version:     version,
		Description: rootCmd.Short,
		GlobalFlags: extractFlags(rootCmd.PersistentFlags()),
		Commands:    buildCommands(rootCmd),
	})
	return nil
}

func buildCommands(c *cobra.Command) []commandInfo {
	var result []commandInfo
	for _, sub := range c.Commands() {
		if sub.Hidden || sub.Name() == "help" || sub.Name() == "completion" {
			continue
		}
		result = append(result, commandInfo{
			Name:        sub.Name(),
			Use:         sub.Use,
			Description: sub.Short,
			Flags:       extractFlags(sub.LocalFlags()),
			Commands:    buildCommands(sub),
		})
	}
	return result
}

func extractFlags(fs *pflag.FlagSet) []flagInfo {
	var result []flagInfo
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		fi := flagInfo{Name: "--" + f.Name, Description: f.Usage}
		if f.Shorthand != "" {
			fi.Shorthand = "-" + f.Shorthand
		}
		if f.DefValue != "" && f.DefValue != "false" {
			fi.Default = f.DefValue
		}
		result = append(result, fi)
	})
	return result
}
