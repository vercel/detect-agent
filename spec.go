package detectagent

import (
	_ "embed"
	"encoding/json"
)

//go:embed agents.json
var agentsJSON []byte

type agentsFile struct {
	AIAgentVar string      `json:"aiAgentVar"`
	Agents     []agentJSON `json:"agents"`
}

type agentJSON struct {
	Key   string    `json:"key"`
	Name  string    `json:"name"`
	Match Condition `json:"match"`
}

var (
	aiAgentVar string
	agentSpecs []AgentSpec
	// KnownAgents maps UPPER_SNAKE keys to canonical agent names, e.g.
	// KnownAgents["CURSOR"] == "cursor". Derived from agents.json at init.
	KnownAgents map[string]string
)

func init() {
	var f agentsFile
	if err := json.Unmarshal(agentsJSON, &f); err != nil {
		panic("detect-agent: failed to parse agents.json: " + err.Error())
	}

	aiAgentVar = f.AIAgentVar
	agentSpecs = make([]AgentSpec, len(f.Agents))
	KnownAgents = make(map[string]string, len(f.Agents))

	for i, a := range f.Agents {
		agentSpecs[i] = AgentSpec{Key: a.Key, Name: a.Name, Match: a.Match}
		KnownAgents[a.Key] = a.Name
	}
}
