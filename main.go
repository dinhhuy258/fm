package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/app"
	"github.com/dinhhuy258/fm/pkg/config"
)

var version = "unversioned"

func main() {
	showVersion := flag.Bool("version", false, "Print the current version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	config.LoadConfig()

	app := app.NewApp()

	err := app.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
