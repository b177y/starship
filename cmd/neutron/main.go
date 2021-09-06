package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Set log level
	log.SetLevel(log.DebugLevel)

	// define -help command flag
	printUsage := flag.Bool("help", false, "Print command line usage")

	// define subcommands
	joinCommand := flag.NewFlagSet("join", flag.ExitOnError)
	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)

	// define flags for join subcommand
	joinQaddrFlag := joinCommand.String("quasar", "", "Quasar address")
	joinNameFlag := joinCommand.String("name", "", "Node name")
	joinNetnameFlag := joinCommand.String("network", "", "Name of network to join")

	// define flags for update subcommand
	updateNetnameFlag := updateCommand.String("network", "", "Name of network to update")

	// parse basic flags (for help option)
	flag.Parse()
	if *printUsage {
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("join or update subcommand is required.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "join":
		joinCommand.Parse(os.Args[2:])
	case "update":
		updateCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if joinCommand.Parsed() {
		// run initialise (join) function
		initialise(*joinQaddrFlag,
			*joinNameFlag,
			*joinNetnameFlag,
		)
	} else if updateCommand.Parsed() {
		// run update function
		update(*updateNetnameFlag)
	}
}
