package main

import (
	"codespacegen/internal/app"
	"os"
)

func main() {
	app := app.NewApp()
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
