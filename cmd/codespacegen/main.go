package main

import (
	"os"
	"codespacegen/internal/app"
)

func main() {
	app := app.NewApp()
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
