package main

import (
	"os"

	"github.com/sv-go-tools/grade/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
