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
		"AGENT_CONTEXT_OUT",
		"AI_AGENT",
		"ANTIGRAVITY_AGENT",
		"ANTIGRAVITY_CLI_ALIAS",
		"APPLICATION_NAME",
		"APP_NAME",
		"AUGMENT_AGENT",
		"CLAUDECODE",
		"CLAUDE_CODE",
		"CLAUDE_CODE_IS_COWORK",
		"CLINE_ACTIVE",
		"CODEBUDDY_CONFIG_DIR",
		"CODEX_CI",
		"CODEX_SANDBOX",
		"CODEX_SANDBOX_NETWORK_DISABLED",
		"CODEX_THREAD_ID",
		"CONFIG_DIR",
		"COPILOT_ALLOW_ALL",
		"COPILOT_GITHUB_TOKEN",
		"COPILOT_MODEL",
		"CRUSH",
		"CURSOR_AGENT",
		"CURSOR_EXTENSION_HOST_ROLE",
		"CURSOR_TRACE_ID",
		"DOUBAO_OFFICE_AGENT_NAME",
		"DOUBAO_OFFICE_ANSWER_MESSAGE_ID",
		"DOUBAO_OFFICE_APP_ID",
		"DOUBAO_OFFICE_CLI_CITATION",
		"DOUBAO_OFFICE_CLI_FORWARD_PROXY",
		"DOUBAO_OFFICE_CLI_TOKEN",
		"DOUBAO_OFFICE_EDITION",
		"DOUBAO_OFFICE_FORWARD_PROXY",
		"DOUBAO_OFFICE_LOG_ID",
		"DOUBAO_OFFICE_MARKET",
		"DOUBAO_OFFICE_PLATFORM_APP_ID",
		"GEMINI_CLI",
		"GOOSE_PROVIDER",
		"GOOSE_TERMINAL",
		"GROK_AGENT",
		"GROK_PLUGIN_DATA",
		"GROK_PLUGIN_ROOT",
		"HATCHLING_ID",
		"HERMES_AGENT",
		"HERMES_SESSION_ID",
		"IS_PACKAGED",
		"JARVIS_IS_ASSIGNED",
		"JARVIS_TOOL_CALL_ID",
		"JUNIE_DATA",
		"JUNIE_SHIM_PATH",
		"KILOCODE_FEATURE",
		"KIMI_CODE_HOME",
		"KIMI_PLUGIN_ROOT",
		"MANUS_ADDON_PREVIEW_PUBLIC_ORIGIN",
		"MANUS_ADDON_REGISTRY_ENABLED",
		"MANUS_ADDON_RUNTIME_PUBLIC_ORIGIN",
		"MANUS_ADDON_RUNTIME_SPECS_REVISION",
		"MANUS_ADDON_RUNTIME_SPEC_PATH",
		"MANUS_ADDON_SERVER_BRANCH",
		"MANUS_CONFIG_DIR",
		"MANUS_GROUP_ENABLED",
		"MANUS_HOME",
		"MANUS_LONG_DESKTOP_WORKER",
		"OPENCLAW_SHELL",
		"OPENCODE",
		"OPENCODE_CLIENT",
		"QODERCN_SERVER_ENDPOINT",
		"QODERWORK_EXECUTION_ENV",
		"QODERWORK_IS_CN",
		"QODERWORK_SOURCE_CHAT_ID",
		"QODER_SITE",
		"QODER_WORK_INTEGRATION_PRODUCT",
		"QWENWORK_BASE_URL",
		"REGION",
		"REPL_ID",
		"SAND_AGENT",
		"SAND_MANIFEST_V2",
		"SAND_SUPERVISOR_ENABLED",
		"SESSION_ID",
		"TERM_PROGRAM",
		"TIER",
		"TIMEZONE",
		"TRAE_AI_SHELL_ID",
		"VTCODE",
		"WORKBUDDY_APPLICATION_NAME",
		"WORKBUDDY_APP_NAME",
		"WORKBUDDY_CONFIG_DIR",
		"WORKBUDDY_IS_PACKAGED",
		"ZED_TERM",
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
