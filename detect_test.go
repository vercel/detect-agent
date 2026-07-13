package detectagent

import (
	"encoding/json"
	"os"
	"testing"
)

func clearAgentEnvs(t *testing.T) {
	t.Helper()
	vars := []string{
		"AI_AGENT",
		"CURSOR_TRACE_ID", "CURSOR_AGENT", "CURSOR_EXTENSION_HOST_ROLE",
		"GEMINI_CLI",
		"CODEX_SANDBOX", "CODEX_CI", "CODEX_THREAD_ID",
		"ANTIGRAVITY_AGENT",
		"AUGMENT_AGENT",
		"OPENCODE_CLIENT", "OPENCODE",
		"GOOSE_PROVIDER",
		"JUNIE_DATA", "JUNIE_SHIM_PATH",
		"CLAUDECODE", "CLAUDE_CODE", "CLAUDE_CODE_IS_COWORK",
		"REPL_ID",
		"COPILOT_MODEL", "COPILOT_ALLOW_ALL", "COPILOT_GITHUB_TOKEN",
		"TERM_PROGRAM",
	}
	for _, v := range vars {
		orig, hadOrig := os.LookupEnv(v)
		os.Unsetenv(v)
		t.Cleanup(func() {
			if hadOrig {
				os.Setenv(v, orig)
			} else {
				os.Unsetenv(v)
			}
		})
	}
	t.Setenv("PATH", "/usr/bin")
}

func assertAgentResult(t *testing.T, got AgentResult, wantIsAgent bool, wantName string) {
	t.Helper()
	if got.IsAgent != wantIsAgent {
		t.Errorf("IsAgent: got %v, want %v", got.IsAgent, wantIsAgent)
	}
	if wantIsAgent {
		if got.Agent == nil {
			t.Fatal("Agent is nil but expected non-nil")
		}
		if got.Agent.Name != wantName {
			t.Errorf("Agent.Name: got %q, want %q", got.Agent.Name, wantName)
		}
	} else {
		if got.Agent != nil {
			t.Errorf("expected Agent nil, got %+v", got.Agent)
		}
	}
}

func TestDetermineAgent(t *testing.T) {
	type testCase struct {
		Name            string            `json:"name"`
		Env             map[string]string `json:"env"`
		TTY             *bool             `json:"tty"`
		Files           []string          `json:"files"`
		SkipGo          bool              `json:"skipGo"`
		ExpectedIsAgent bool              `json:"expectedIsAgent"`
		ExpectedAgentKey string           `json:"expectedAgentKey"`
		ExpectedName    string            `json:"expectedName"`
	}

	data, err := os.ReadFile("testcases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []testCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		tc := tc
		if tc.SkipGo {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			clearAgentEnvs(t)

			orig := isTTYFn
			if tc.TTY != nil {
				val := *tc.TTY
				isTTYFn = func() bool { return val }
			} else {
				isTTYFn = func() bool { return false }
			}
			t.Cleanup(func() { isTTYFn = orig })

			for k, v := range tc.Env {
				t.Setenv(k, v)
			}

			for _, path := range tc.Files {
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatalf("failed to create %s: %v", path, err)
				}
				t.Cleanup(func() { os.RemoveAll(path) })
			}

			result, err := DetermineAgent()
			if err != nil {
				t.Fatal(err)
			}

			wantName := tc.ExpectedName
			if tc.ExpectedAgentKey != "" {
				wantName = KnownAgents[tc.ExpectedAgentKey]
			}
			assertAgentResult(t, result, tc.ExpectedIsAgent, wantName)
		})
	}
}
