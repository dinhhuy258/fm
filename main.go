package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dinhhuy258/fm/pkg/app"
)

var version = "unversioned"

func main() {
	showVersion := flag.Bool("version", false, "Print the current version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	fmt.Println("Running fm...")
	app, err := app.NewApp()
	if err != nil {
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		os.Exit(1)
	}
}
