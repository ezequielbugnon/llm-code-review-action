package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type FileChanges struct {
	Current string `json:"current"`
	Changes string `json:"changes"`
}

type InputData struct {
	InputData map[string]FileChanges `json:"input_data"`
}

type StackSpoTAgent struct {
	url          string
	urlPost      string
	clientID     string
	clientSecret string
}

type StackSpoTAgentResponse struct {
	City        string
	Description string
	Enable      bool
}

type Progress struct {
	Duration            int       `json:"duration"`
	End                 time.Time `json:"end"`
	ExecutionPercentage float64   `json:"execution_percentage"`
	Start               time.Time `json:"start"`
	Status              string    `json:"status"`
}

type StepResult struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

type Step struct {
	ExecutionOrder int        `json:"execution_order"`
	StepName       string     `json:"step_name"`
	StepResult     StepResult `json:"step_result"`
	Type           string     `json:"type"`
}

type Result struct {
	ConversationID   string   `json:"conversation_id"`
	ExecutionID      string   `json:"execution_id"`
	Progress         Progress `json:"progress"`
	QuickCommandSlug string   `json:"quick_command_slug"`
	Result           string   `json:"result"`
	Steps            []Step   `json:"steps"`
}

type City struct {
	City                   string `json:"city"`
	Enable                 string `json:"enable"`
	PorcentagemDeSeguranca int    `json:"porcentagem-de-segurança"`
	Descricao              string `json:"descrição"`
}

func main() {
	fileChanges := make(map[string]FileChanges)

	output, err := exec.Command("git", "diff", "--name-only", "HEAD^", "HEAD").Output()
	if err != nil {
		log.Println("Error al obtener archivos cambiados: ", err)
		return
	}

	log.Println("obtener", string(output))

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		currentContent, err := exec.Command("git", "show", "HEAD:"+file).Output()
		if err != nil {
			log.Println("Error al obtener contenido actual de ", file, err)
			continue
		}

		log.Println("current", string(currentContent))

		changes, err := exec.Command("git", "diff", "--unified=0", "HEAD^", "HEAD", "--", file).Output()
		if err != nil {
			log.Println("Error al obtener cambios de ", file, err)
			continue
		}

		log.Println("changes", string(changes))

		fileChanges[file] = FileChanges{
			Current: string(currentContent),
			Changes: string(changes),
		}
	}

	inputData := InputData{
		InputData: fileChanges,
	}

	IAStackSpot := NewStackSpotAgent("https://genai-code-buddy-api.stackspot.com/v1/quick-commands/callback/", "https://genai-code-buddy-api.stackspot.com/v1/quick-commands/create-execution/safe-explorer-api", "7f4a1870-4f93-4fe1-bd2e-146b270f269b", "ZkCg9OdKwKI2turlQkE23F5r7G2UW0wu49r4NzCvV930jIWSMP5h44ASSu38Db6F")

	review, err := IAStackSpot.GetDataFromEndpoint(inputData)
	if err != nil {
		log.Println("Error getFromDataEndpoint ", err)
	}

	log.Println("getFromDataEndpoint ", review)

	fmt.Println("hi")
}

func NewStackSpotAgent(url string, urlPost string, clientID string, clientSecret string) *StackSpoTAgent {
	return &StackSpoTAgent{
		url:          url,
		urlPost:      urlPost,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *StackSpoTAgent) GetDataFromEndpoint(inputData InputData) (review string, err error) {
	ticker := time.NewTicker(time.Second * 2)
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

			if result.Progress.Status == "COMPLETED" {
				if len(result.Steps) == 0 {
					errCh <- fmt.Errorf("there are no steps in the result")
					return
				}

				answer := result.Steps[0].StepResult.Answer
				//re := regexp.MustCompile(`(?s)\{.*\}`)
				//jsonStr := re.FindString(answer)
				//jsonStr = strings.TrimSpace(jsonStr)

				/*err = json.Unmarshal([]byte(answer), &review)
				if err != nil {
					errCh <- fmt.Errorf("error deserializing the response: %w", err)
					return
				}*/
				review = answer

				fmt.Println("The process has finished.")
				errCh <- nil
				return
			}

			fmt.Println("The process is not finished yet, waiting for the next interval...")
		}
	}()

	if err := <-errCh; err != nil {
		return "", err
	}

	return review, nil
}

func (s *StackSpoTAgent) getToken() (string, error) {
	var result map[string]interface{}
	url := "https://idm.stackspot.com/itau/oidc/oauth/token"
	data := fmt.Sprintf("client_id=%s&grant_type=client_credentials&client_secret=%s", s.clientID, s.clientSecret)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
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

type Data struct {
	CitiesData CitiesData `json:"input_data"`
}

type CitiesData struct {
	Cities []string `json:"cities"`
}

func (s *StackSpoTAgent) createExecution(token string, inputData InputData) (callback string, err error) {

	log.Println(inputData)

	data := Data{
		CitiesData: CitiesData{
			Cities: []string{"sao paulo"},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return callback, err
	}

	req, err := http.NewRequest("POST", s.urlPost, bytes.NewBuffer(jsonData))
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
	url := fmt.Sprintf(s.url+"%s", callbackID)
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
