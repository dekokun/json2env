package main

import (
	"context"
	"flag"
	"log"
	"os"

	json2env2 "github.com/dekokun/json2env"
)

func main() {
	log.SetFlags(0)
	err := json2env2.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr, os.Stdin)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}
