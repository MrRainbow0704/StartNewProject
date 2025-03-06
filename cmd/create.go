package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createCommand(path, template string) {
	cwd, _ := os.Getwd()
	t, err := filepath.Abs(filepath.Join(cwd + "/templates/" + template))
	if err != nil {
		panic(err)
	}
	fmt.Println("Creazione del progetto in corso...")
	fmt.Println("Template:", template)
	if err := os.CopyFS(path, os.DirFS(t)); err != nil {
		panic(err)
	}
	filepath.WalkDir(path, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".command") {
			os.Remove(s)
		}
		return nil
	})

	reader := bufio.NewReader(os.Stdin)
	filepath.WalkDir(t, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".command") {
			b, err := os.ReadFile(s)
			if err != nil {
				panic(err)
			}
			cmds := strings.Split(strings.TrimSpace(string(b)), "\n")
			for _, cmdx := range cmds {
				cmd := strings.Split(cmdx, " ")
				if len(cmd) < 2 {
					cmd = append(cmd, "")
				}

				fmt.Printf("Eseguendo il comando `%s`\n\tcontenuto in `%s`.\n\tContinuare? [S/n] ", strings.Join(cmd, " "), strings.TrimPrefix(s, t))
				ok, _ := reader.ReadString('\n')
				ok = strings.TrimSpace(strings.ToLower(ok))
				if ok == "s" || ok == "y" {
					c := exec.Command(cmd[0], cmd[1:]...)
					c.Dir = filepath.Join(path, strings.TrimPrefix(strings.TrimSuffix(s, d.Name()), t))
					c.Stderr = os.Stderr
					c.Stdout = os.Stdout
					c.Stdin = os.Stdin
					if err := c.Run(); err != nil {
						panic(err)
					}
				} else {
					fmt.Println("Comando saltato.")
				}
			}
		}
		return nil
	})
}