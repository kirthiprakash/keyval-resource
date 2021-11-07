package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gstack.io/concourse/keyval-resource/models"
)

func main() {
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}

	response := models.CheckResponse{}
	if len(request.Version) > 0 {
		response = models.CheckResponse{request.Version}
	}
	json.NewEncoder(os.Stdout).Encode(response)
}
