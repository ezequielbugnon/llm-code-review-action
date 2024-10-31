package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"script/fetch"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	urlCallback := os.Getenv("URLCALLBACK")
	urlExecution := os.Getenv("URLEXECUTION")
	urlToken := os.Getenv("URLTOKEN")
	clientID := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")

	fileChanges := make(map[string]fetch.FileChanges)

	output, err := exec.Command("git", "diff", "--name-only", "HEAD^", "HEAD").Output()
	if err != nil {
		log.Println("Error to get all files changes: ", err)
		return
	}

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		currentContent, err := exec.Command("git", "show", "HEAD:"+file).Output()
		if err != nil {
			log.Println("Error to get atual file content", file, err)
			continue
		}

		changes, err := exec.Command("git", "diff", "--unified=0", "HEAD^", "HEAD", "--", file).Output()
		if err != nil {
			log.Println("Error to get file changes", file, err)
			continue
		}

		fileChanges[file] = fetch.FileChanges{
			Current: string(currentContent),
			Changes: string(changes),
		}
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
