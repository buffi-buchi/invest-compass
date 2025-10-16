package main

import (
	"os"

	"github.com/buffi-buchi/invest-compass/backend/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		os.Exit(1)
	}
}
