package detectagent

// Condition is a node in a condition tree loaded from agents.json.
type Condition struct {
	Type       string      `json:"type"`
	Name       string      `json:"name,omitempty"`
	Value      string      `json:"value,omitempty"`
	Pattern    string      `json:"pattern,omitempty"`
	Path       string      `json:"path,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

type AgentSpec struct {
	Key   string
	Name  string
	Match Condition
}

type AgentDetails struct {
	Name string
}

type AgentResult struct {
	IsAgent bool
	Agent   *AgentDetails
}
