package main

import (
	"gitlab.com/tokend/course-certificates/ccp/internal/cli"
	"os"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
