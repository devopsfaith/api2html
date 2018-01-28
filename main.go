package main

import (
	"log"

	"github.com/devopsfaith/api2html/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println("error:", err.Error())
	}
}
