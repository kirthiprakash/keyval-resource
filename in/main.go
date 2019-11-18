package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"bufio"
	"fmt"
	"github.com/moredhel/keyval-resource/models"
)

var (
	destination string
)

func create_file(key string, value string) {
	output := filepath.Join(destination, key)

	file, err := os.Create(output)
	if err != nil {
		fatal("creating output file", err)
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", value)

	err = w.Flush()

	if err != nil {
		fatal("writing output file", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <destination>")
	}

	destination = os.Args[1]

	log("creating destination dir " + destination)
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	for k, v := range request.Version {
		create_file(k, v)
	}
	log("Done")
}

func fatal(doing string, err error) {
	fmt.Fprintln(os.Stderr, "error "+doing+": "+err.Error())
	os.Exit(1)
}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}
