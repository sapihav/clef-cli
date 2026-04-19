// Package cmd wires the Cobra command tree for the clef CLI.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is overridden by goreleaser ldflags at build time.
var version = "dev"

// Exit codes.
const (
	ExitSuccess     = 0
	ExitError       = 1 // user / config error
	ExitSystemError = 2 // read / write failure
)

var (
	flagJSON       bool
	flagJSONErrors bool
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
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output result as JSON")
	rootCmd.PersistentFlags().BoolVar(&flagJSONErrors, "json-errors", false, "Output errors as JSON on stderr")
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

// userError prints msg to stderr and returns exit code 1.
func userError(msg string) error {
	writeErr(msg, ExitError)
	return &exitCodeError{code: ExitError, msg: msg}
}

// sysError prints msg to stderr and returns exit code 2.
func sysError(msg string) error {
	writeErr(msg, ExitSystemError)
	return &exitCodeError{code: ExitSystemError, msg: msg}
}

func writeErr(msg string, code int) {
	if flagJSONErrors {
		b, _ := json.Marshal(map[string]interface{}{"error": msg, "code": code})
		fmt.Fprintln(os.Stderr, string(b))
	} else {
		fmt.Fprintln(os.Stderr, "clef:", msg)
	}
}
