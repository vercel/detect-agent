package detectagent

import (
	_ "embed"
	"encoding/json"
	"errors"
	"testing"
)

//go:embed testcases.json
var testCasesJSON []byte

func clearAgentEnvs(t *testing.T) {
	t.Helper()
	for _, v := range []string{
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
	} {
		t.Setenv(v, "")
	}
	t.Setenv("PATH", "/usr/bin")
}

func assertAgentResult(t *testing.T, got *AgentDetails, err error, wantIsAgent bool, wantName string) {
	t.Helper()
	if wantIsAgent {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected AgentDetails, got nil")
		}
		if got.Name != wantName {
			t.Errorf("Name: got %q, want %q", got.Name, wantName)
		}
	} else {
		if !errors.Is(err, ErrAgentNotFound) {
			t.Errorf("expected ErrAgentNotFound, got err=%v agent=%v", err, got)
		}
	}
}

func TestDetermineAgent(t *testing.T) {
	type testCase struct {
		Name             string            `json:"name"`
		Env              map[string]string `json:"env"`
		TTY              bool              `json:"tty"`
		Files            []string          `json:"files"`
		SkipGo           bool              `json:"skipGo"`
		ExpectedIsAgent  bool              `json:"expectedIsAgent"`
		ExpectedAgentKey string            `json:"expectedAgentKey"`
		ExpectedName     string            `json:"expectedName"`
	}

	var cases []testCase
	if err := json.Unmarshal(testCasesJSON, &cases); err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		if tc.SkipGo {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			clearAgentEnvs(t)

			orig := isTTYFn
			isTTYFn = func() bool { return tc.TTY }
			t.Cleanup(func() { isTTYFn = orig })

			if len(tc.Files) > 0 {
				t.Fatal("file-system fixtures are not supported in Go tests; set skipGo: true")
			}

			for k, v := range tc.Env {
				t.Setenv(k, v)
			}

			agent, err := Detect()

			wantName := tc.ExpectedName
			if tc.ExpectedAgentKey != "" {
				wantName = KnownAgents[tc.ExpectedAgentKey]
			}
			assertAgentResult(t, agent, err, tc.ExpectedIsAgent, wantName)
		})
	}
}
