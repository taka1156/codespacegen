package main

import (
	"codespacegen/internal/app"
	"fmt"
	"os"
)

func main() {
	app := app.NewApp()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
