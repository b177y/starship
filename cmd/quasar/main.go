package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Set log level
	log.SetLevel(log.DebugLevel)

	// define -help command flag
	printUsage := flag.Bool("help", false, "Print command line usage")

	serveConfig := flag.String("config",
		"/etc/quasar/config.yml",
		"Quasar config file path",
	)
	serveListenAddress := flag.String("host",
		"",
		"Quasar server listen address",
	)
	serveListenPort := flag.Int("port",
		0,
		"Quasar server listen port",
	)
	flag.Parse()

	if *printUsage {
		flag.Usage()
		os.Exit(0)
	}

	runServe(*serveConfig,
		*serveListenAddress,
		*serveListenPort,
	)
}
