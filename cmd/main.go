package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var cwd string

func init() {
	cwd, _ = os.Getwd()
}

func main() {
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
	run(path, template, info, list)
}

func run(path, template, info string, list bool) {
	if list {
		listCommand()
	} else if info != "D" {
		infoCommand(info)
	} else {
		createCommand(path, template)
	}
}

func noFlags() bool {
	found := true
	flag.Visit(func(f *flag.Flag) {
		found = false
	})
	return found
}
