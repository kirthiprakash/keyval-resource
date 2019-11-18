package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"fmt"
	"github.com/google/uuid"
	"github.com/moredhel/keyval-resource/models"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <destination>")
	}

	destination := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	data := make(map[string]string)

	// read in files
	err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		fileName := info.Name()

		// don't supported nested maps
		if info.IsDir() {
			log(fmt.Sprintf("skipping directory %s", fileName))
			return nil
		}

		inputFile := filepath.Join(destination, fileName)
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return err
		}

		data[fileName] = string(content)

		return nil
	})

	if err != nil {
		fatal("could not open directory", err)
	}

	// override with request keys
	for k, v := range request.Params {
		data[k] = v
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
