package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"script/fetch"
)

func main() {

	urlCallback := os.Getenv("URLCALLBACK")
	urlExecution := os.Getenv("URLEXECUTION")
	urlToken := os.Getenv("URLTOKEN")
	clientID := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")
	fileChangesJSON := os.Getenv("INPUT_FILECHANGES")

	var fileChanges map[string]fetch.FileChanges

	if err := json.Unmarshal([]byte(fileChangesJSON), &fileChanges); err != nil {
		fmt.Fprintf(os.Stderr, "Error to parser fileChanges: %v\n", err)
		os.Exit(1)
	}

	log.Println("Files sending to review", len(fileChanges))

	inputData := fetch.InputData{
		InputData: fileChanges,
	}

	IAStackSpot := fetch.NewStackSpotAgent(urlCallback, urlExecution, urlToken, clientID, clientSecret)

	review, err := IAStackSpot.GetDataFromEndpoint(inputData)
	if err != nil {
		log.Println("Error getFromDataEndpoint of StackSpot", err)
	}

	fmt.Println(review)
}
