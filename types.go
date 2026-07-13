package detectagent

import "errors"

// ErrAgentNotFound is returned by DetermineAgent when the process is not
// running inside a known AI agent environment.
var ErrAgentNotFound = errors.New("agent not found")

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
