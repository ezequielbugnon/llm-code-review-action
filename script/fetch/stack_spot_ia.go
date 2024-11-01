package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func NewStackSpotAgent(urlCallback string, urlExecution string, urlToken string, clientID string, clientSecret string) *StackSpoTAgent {
	return &StackSpoTAgent{
		UrlCallback:  urlCallback,
		UrlExecution: urlExecution,
		UrlToken:     urlToken,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *StackSpoTAgent) GetDataFromEndpoint(inputData InputData) (review string, err error) {
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()
	errCh := make(chan error, 1)

	token, err := s.getToken()
	if err != nil {
		return review, fmt.Errorf("error obtaining the token: %w", err)
	}

	callback, err := s.createExecution(token, inputData)
	if err != nil {
		return review, fmt.Errorf("error creating the execution %w", err)
	}

	go func() {
		for range ticker.C {
			result, err := s.getCallback(token, callback)
			if err != nil {
				errCh <- err
				return
			}

			log.Println(result.Steps)

			if result.Progress.Status == "COMPLETED" {
				if len(result.Steps) == 0 {
					errCh <- fmt.Errorf("there are no steps in the result")
					return
				}

				review = result.Steps[0].StepResult.Answer

				log.Println("The process has finished.")
				errCh <- nil
				return
			}

			log.Println("The process is not finished yet, waiting for code review response")
		}
	}()

	if err := <-errCh; err != nil {
		return "", err
	}

	return review, nil
}

func (s *StackSpoTAgent) getToken() (string, error) {
	var result map[string]interface{}
	data := fmt.Sprintf("client_id=%s&grant_type=client_credentials&client_secret=%s", s.clientID, s.clientSecret)

	req, err := http.NewRequest("POST", s.UrlToken, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

func (s *StackSpoTAgent) createExecution(token string, inputData InputData) (callback string, err error) {

	jsonData, err := json.Marshal(inputData)
	if err != nil {
		return callback, err
	}

	req, err := http.NewRequest("POST", s.UrlExecution, bytes.NewBuffer(jsonData))
	if err != nil {
		return callback, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return callback, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return callback, fmt.Errorf("failed to create execution: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return callback, err
	}

	if err := json.Unmarshal(body, &callback); err != nil {
		return callback, err
	}

	return callback, nil
}

func (s *StackSpoTAgent) getCallback(token, callbackID string) (result Result, err error) {
	url := fmt.Sprintf(s.UrlCallback+"%s", callbackID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to get callback: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	return result, nil
}
