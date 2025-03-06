package main

import (
	"fmt"
	"os"
)

func listCommand() {
	fmt.Println("Lista dei template disponibili:")
	templs, err := os.ReadDir(cwd + "/templates/")
	if err != nil {
		panic(err)
	}

	for _, t := range templs {
		if t.IsDir() {
			fmt.Println(" - ", t.Name())
		}
	}
}