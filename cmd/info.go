package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func infoCommand(info string) {
	fmt.Println("Informazioni sul template:", info)
	t := filepath.Join(cwd, "templates", info)
	filepath.WalkDir(t, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		x := strings.TrimPrefix(s, t)
		depth := strings.Count(x, "/") + strings.Count(x, "\\") - 1
		if x == "" {
			return nil
		}

		if depth > 0 {
			fmt.Print(strings.Repeat("│   ", depth-1))
			fmt.Print("├── ")
		}
		if !d.IsDir() {
			fmt.Print(d.Name())
			if strings.HasSuffix(d.Name(), ".command") {
				b, err := os.ReadFile(s)
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