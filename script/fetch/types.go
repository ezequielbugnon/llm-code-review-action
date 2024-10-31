package fetch

import "time"

type FileChanges struct {
	Current string `json:"current"`
	Changes string `json:"changes"`
}

type InputData struct {
	InputData map[string]FileChanges `json:"input_data"`
}

type StackSpoTAgent struct {
	UrlCallback  string
	UrlExecution string
	UrlToken     string
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
