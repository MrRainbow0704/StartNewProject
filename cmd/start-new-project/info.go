package main

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/MrRainbow0704/StartNewProject/templates"
)

func infoCommand(info string) {
	fmt.Println("Informazioni sul template:", info)
	fs.WalkDir(templates.Content, info, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if s == info {
			return nil
		}

		depth := strings.Count(s, "/") + strings.Count(s, "\\") - 1
		if depth > 0 {
			fmt.Print(strings.Repeat("│   ", depth-1))
			fmt.Print("├── ")
		}
		if !d.IsDir() {
			fmt.Print(d.Name())
			if strings.HasSuffix(d.Name(), ".command") {
				b, err := templates.Content.ReadFile(s)
				if err != nil {
					panic(err)
				}
				cmds := strings.Join(strings.Split(strings.TrimSpace(string(b)), "\n"), " && ")
				fmt.Printf(" > %s", cmds)
			}
			fmt.Println()
		} else {
			fmt.Println(d.Name() + "/")
		}
		return nil
	})
}
