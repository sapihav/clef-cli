package cmd

import (
	"encoding/json"
	"fmt"
)

// flagDryRun is shared by all mutating commands.
var flagDryRun bool

// writeJSON marshals v as indented JSON and prints it to stdout.
func writeJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
