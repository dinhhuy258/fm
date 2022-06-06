package main

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/lua"
)

var version = "unversioned"

func main() {
	l := lua.NewLua()
	err := l.LoadConfig("hello.lua")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	// showVersion := flag.Bool("version", false, "Print the current version")
	// flag.Parse()
	//
	// if *showVersion {
	// 	fmt.Println(version)
	// 	os.Exit(0)
	// }
	//
	// if err := config.LoadConfig(); err != nil {
	// 	log.Fatalf("failed to load application configuration %v", err)
	// }
	//
	// app, err := app.NewApp()
	// if err != nil {
	// 	log.Fatalf("failed to new app %v", err)
	// }
	//
	// _ = app.Run()
	// app.OnQuit()
}
