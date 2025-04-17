package main

import (
	"github.com/NibiruChain/pricefeeder/cmd"
)

// Entry point for the pricefeeder application.
// Delegates execution to the cmd package which handles CLI and app initialization.
func main() {
	cmd.Execute()
}
