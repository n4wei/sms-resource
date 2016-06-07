package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nickwei84/sms-resource/lib/awsclient"
	"github.com/nickwei84/sms-resource/out/application"
	"github.com/nickwei84/sms-resource/out/models"
)

func main() {
	var (
		config models.SMSConfig
		client application.SMSService
	)

	err := getStdinInput(&config)
	if err != nil {
		exitWithErr(err)
	}

	err = config.CheckInput()
	if err != nil {
		exitWithErr(err)
	}

	client = awsclient.NewAWSClient(config.Source.AWSAccessKeyID, config.Source.AWSSecretAccessKey)
	app := application.NewApplication(client, config)

	err = app.Run()
	if err != nil {
		exitWithErr(err)
	}

	stdoutOutput, err := generateStdoutOutput()
	if err != nil {
		exitWithErr(err)
	}

	fmt.Println(string(stdoutOutput))
}

func exitWithErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func getStdinInput(config *models.SMSConfig) error {
	stdinData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %v", err)
	}

	err = json.Unmarshal(stdinData, config)
	if err != nil {
		return fmt.Errorf("error parsing stdin as JSON: %v", err)
	}

	return nil
}

func generateStdoutOutput() ([]byte, error) {
	output := models.OutputJSON{
		Version: models.Time{
			Time: time.Now().UTC(),
		},
		Metadata: []models.MetadataItem{},
	}

	stdoutOutput, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("error marshalling output for stdout: %v", err)
	}

	return stdoutOutput, nil
}
