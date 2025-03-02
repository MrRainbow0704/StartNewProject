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
  -d, --dir		percorso dove creare il progetto		[Default: percorso corrente]
  -t, --template		sceglie che template usare		[Obbligatorio]
  --no-git		non inizializza la repository git
	`
	flag.Usage = func() { fmt.Print(usage) }

	var path string
	flag.StringVar(&path, "dir", cwd, "percorso dove salvare i file")
	flag.StringVar(&path, "d", cwd, "percorso dove salvare i file")
	var template string
	flag.StringVar(&template, "template", "D", "template da usare")
	flag.StringVar(&template, "t", "D", "template da usare")
	var git bool
	flag.BoolVar(&git, "no-git", false, "non inizializza la repository git")
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
	} else if !noFlags() && template == "D" && !list {
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
		templs, err := os.ReadDir("../templates")
		if err != nil {
			panic(err)
		}

		for _, t := range templs {
			if t.IsDir() {
				fmt.Println(" - ", t.Name())
			}
		}
	} else {
		run(path, template, git)
	}
}

func run(path, template string, git bool) {
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
	if git {
		if err := exec.Command("git", "init", path).Run(); err != nil {
			panic(err)
		}
	}

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
