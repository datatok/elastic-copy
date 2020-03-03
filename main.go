package main

import (
	"github.com/ebuildy/elastic-copy/commands"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	cmd := commands.NewRootCmd(os.Stdout, os.Args[1:])

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}