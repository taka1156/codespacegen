package main

import (
	"fmt"
	"os"

	"github.com/taka1156/codespacegen/internal/app"
)

func main() {
	launchApp := app.NewApp()
	if err := launchApp.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
