// clef manages .claude/settings.local.json in the current working directory.
//
// See `clef --help` for usage. All command wiring lives in ./cmd and
// settings logic in ./internal/settings.
package main

import (
	"os"

	"github.com/sapihav/clef-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
