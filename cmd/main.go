package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cwd, _ := os.Getwd()

	// initialize flags arguments
	const usage = `Flag di startnewproject:
  -h, --help		stampa le informazioni di aiuto
  -l, --list		lista i template disponibili
  -i, --info		informazioni su un template
  -d, --dir		percorso dove creare il progetto	[Default: percorso corrente]
  -t, --template	sceglie che template usare		[Obbligatorio]
	`
	flag.Usage = func() { fmt.Print(usage) }

	var path string
	flag.StringVar(&path, "dir", cwd, "percorso dove creare il progetto")
	flag.StringVar(&path, "d", cwd, "percorso dove creare il progetto")
	var template string
	flag.StringVar(&template, "template", "D", "template da usare")
	flag.StringVar(&template, "t", "D", "template da usare")
	var info string
	flag.StringVar(&info, "info", "D", "informazioni su un templat")
	flag.StringVar(&info, "i", "D", "informazioni su un templat")
	var list bool
	flag.BoolVar(&list, "list", false, "lista i template disponibili")
	flag.BoolVar(&list, "l", false, "lista i template disponibili")

	flag.Parse()

	if noFlags() {
		// initializing reader
		reader := bufio.NewReader(os.Stdin)

		// Getting the page link
		fmt.Print("Inserisci il template da utilizzare: ")
		template, _ = reader.ReadString('\n')

		// Getting the path where to store the downloads
		fmt.Printf("Inserisci il percorso dove creare il progetto [Vuoto per: \"%s\"]: ", path)
		path, _ = reader.ReadString('\n')
	} else if !noFlags() && template == "D" && !list && info == "D" {
		panic("Il flag --template è obbligatorio.")
	}

	// Input formatting
	template = strings.TrimSpace(template)
	path = strings.TrimSpace(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		panic(fmt.Sprintf("Errore durante la creazione della directory `%s`: %s", path, err))
	}

	if list {
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
	} else if info != "D" {
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
					fmt.Printf(" > %s",cmds)
				}
				fmt.Println()
			} else {
				fmt.Println(d.Name()+"/")
			}
			return nil 
		})
	} else {
		run(path, template)
	}
}

func run(path, template string) {
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

func noFlags() bool {
	found := true
	flag.Visit(func(f *flag.Flag) {
		found = false
	})
	return found
}
