package main

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/cli"
	"log"
	"os"
)

func main() {
	if err := cli.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
