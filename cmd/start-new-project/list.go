package main

import (
	"fmt"

	"github.com/MrRainbow0704/StartNewProject/templates"
)

func listCommand() {
	fmt.Println("Lista dei template disponibili:")

	templ, err := templates.Content.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, t := range templ {
		if t.IsDir() {
			fmt.Println(" - ", t.Name())
		}
	}
}
