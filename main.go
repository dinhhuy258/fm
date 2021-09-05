package main

import (
	"flag"
	"fmt"
	"log"
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

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
