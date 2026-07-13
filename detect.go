package detectagent

import (
	"os"
	"strings"
)

// DetermineAgent inspects the environment and returns which AI agent is
// running, if any. AI_AGENT takes highest priority; after that the specs in
// agentSpecs are evaluated in order and the first match wins.
func DetermineAgent() (AgentResult, error) {
	if name := resolveAIAgentVar(); name != "" {
		return AgentResult{IsAgent: true, Agent: &AgentDetails{Name: name}}, nil
	}

	for _, spec := range agentSpecs {
		matched, err := EvaluateCondition(spec.Match)
		if err != nil {
			return AgentResult{}, err
		}
		if matched {
			return AgentResult{IsAgent: true, Agent: &AgentDetails{Name: spec.Name}}, nil
		}
	}

	return AgentResult{IsAgent: false}, nil
}

func resolveAIAgentVar() string {
	raw := os.Getenv(aiAgentVar)
	return strings.TrimSpace(raw)
}
