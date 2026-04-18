// Package cmd wires the Cobra command tree for the clef CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is overridden by goreleaser ldflags at build time.
var version = "dev"

// Exit codes.
const (
	ExitSuccess = 0
	ExitError   = 1
)

var rootCmd = &cobra.Command{
	Use:           "clef",
	Short:         "Manage .claude/settings.local.json per project/worktree",
	Long:          "clef reads and writes .claude/settings.local.json in the current directory.\n\nIt lets developers set Claude Code model and effort level per project or worktree\nwithout touching global settings.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.Version = version
}

// Execute runs the root command and returns an exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if e, ok := err.(*exitCodeError); ok {
			return e.code
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		return ExitError
	}
	return ExitSuccess
}

// exitCodeError lets subcommands signal the exit code back to Execute without
// printing the message twice (the subcommand already wrote to stderr).
type exitCodeError struct {
	code int
	msg  string
}

func (e *exitCodeError) Error() string { return e.msg }

// userError prints msg to stderr and returns an exitCodeError with ExitError.
func userError(msg string) error {
	fmt.Fprintln(os.Stderr, "clef:", msg)
	return &exitCodeError{code: ExitError, msg: msg}
}
