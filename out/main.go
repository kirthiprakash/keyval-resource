package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"gstack.io/concourse/keyval-resource/models"
)

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <source-dir>")
	}

	sourceDir := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	data := make(map[string]string)

	if request.Params.Directory == "" {
		fatalNoErr("missing parameter 'directory'. Which artifact directory is to be scanned?")
	}
	artifactDir := filepath.Join(sourceDir, request.Params.Directory)

	// read in files
	files, err := ioutil.ReadDir(artifactDir)

	if err != nil {
		fatal("could not open directory", err)
		return
	}
	log(fmt.Sprintf("reading directory '%s'", artifactDir))

	for _, file := range files {
		fileName := file.Name()

		// don't supported nested maps
		if file.IsDir() {
			log(fmt.Sprintf("skipping directory %s", fileName))
			continue
		}

		inputFile := filepath.Join(artifactDir, fileName)
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			log(fmt.Sprintf("skipping file %s", fileName))
			continue
		}

		data[fileName] = string(content)
	}

	// override with request keys
	for key, val := range request.Params.Overrides {
		data[key] = val
	}

	data["UPDATED"] = time.Now().Format(time.RFC850)
	data["UUID"] = uuid.New().String()
	log("read " + strconv.Itoa(len(data)) + " keys from input file")

	json.NewEncoder(os.Stdout).Encode(models.OutResponse{
		Version: data,
	})

}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}
