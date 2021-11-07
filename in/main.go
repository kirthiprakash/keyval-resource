package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gstack.io/concourse/keyval-resource/models"
)

var (
	destination string
)

func createFile(name string, contents string) {
	filePath := filepath.Join(destination, name)

	file, err := os.Create(filePath)
	if err != nil {
		fatal("creating output file", err)
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", contents)

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

	for key, val := range request.Version {
		createFile(key, val)
	}
	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: request.Version,
	})
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
