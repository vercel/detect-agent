package detectagent

import (
	"errors"
	"os"
	"strings"
)

// ErrAgentNotFound is returned by DetermineAgent when the process is not
// running inside a known AI agent environment.
var ErrAgentNotFound = errors.New("agent not found")

type AgentDetails struct {
	Name string
}

// DetermineAgent inspects the environment and returns which AI agent is
// running, if any. AI_AGENT takes highest priority; after that the specs in
// agentSpecs are evaluated in order and the first match wins.
// Returns ErrAgentNotFound when not running inside a known agent environment.
func DetermineAgent() (*AgentDetails, error) {
	if name := resolveAIAgentVar(); name != "" {
		return &AgentDetails{Name: name}, nil
	}

	for _, spec := range agentSpecs {
		matched, err := EvaluateCondition(spec.Match)
		if err != nil {
			return nil, err
		}
		if matched {
			return &AgentDetails{Name: spec.Name}, nil
		}
	}

	return nil, ErrAgentNotFound
}

func resolveAIAgentVar() string {
	raw := os.Getenv(aiAgentVar)
	return strings.TrimSpace(raw)
}
