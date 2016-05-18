package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func handleErr(errMsg string) {
	fmt.Fprintln(os.Stderr, errMsg)
	os.Exit(1)
}

func main() {
	var output struct {
		Version interface{} `json:"version"`
	}

	stdinData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		handleErr(fmt.Sprintf("Error reading from stdin: %v", err))
	}

	err = json.Unmarshal(stdinData, &output)
	if err != nil {
		handleErr(fmt.Sprintf("Error unmarshalling JSON: %v", err))
	}

	if output.Version == nil {
		fmt.Fprintf(os.Stderr, "Error: version key pair is missing from stdin")
		os.Exit(1)
	}

	stdoutOutput, err := json.Marshal(output)
	if err != nil {
		handleErr(fmt.Sprintf("Error marshalling output for stdout: %v", err))
	}

	fmt.Printf("%s", []byte(stdoutOutput))
}
