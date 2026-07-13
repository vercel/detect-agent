package detectagent

import (
	"os"
	"testing"
)

// clearAgentEnvs unsets all env vars that any agent spec checks so that a
// freshly-reset test environment doesn't detect a false positive from the
// real process environment.
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
		// Neutralize PATH and TERM_PROGRAM so real env can't leak.
		"PATH", "TERM_PROGRAM",
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
	// Set PATH to a neutral value so env_matches on PATH returns false.
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

func TestDetermineAgent_CustomAIAgent(t *testing.T) {
	t.Run("returns no agent when AI_AGENT not set", func(t *testing.T) {
		clearAgentEnvs(t)
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("detects custom agent from AI_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "custom-agent")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, "custom-agent")
	})
}

func TestDetermineAgent_GithubCopilot(t *testing.T) {
	t.Run("detects github copilot from AI_AGENT=github-copilot", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "github-copilot")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GITHUB_COPILOT"])
	})

	t.Run("detects github copilot from COPILOT_MODEL", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("COPILOT_MODEL", "gpt-5")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GITHUB_COPILOT"])
	})

	t.Run("detects github copilot from COPILOT_ALLOW_ALL", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("COPILOT_ALLOW_ALL", "true")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GITHUB_COPILOT"])
	})

	t.Run("detects github copilot from COPILOT_GITHUB_TOKEN", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("COPILOT_GITHUB_TOKEN", "ghp_xxx")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GITHUB_COPILOT"])
	})
}

func TestDetermineAgent_Cursor(t *testing.T) {
	t.Run("returns no agent when CURSOR_TRACE_ID not set", func(t *testing.T) {
		clearAgentEnvs(t)
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("detects cursor from CURSOR_TRACE_ID", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_TRACE_ID", "some-uuid")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CURSOR"])
	})
}

func TestDetermineAgent_CursorCLI(t *testing.T) {
	t.Run("detects cursor cli from CURSOR_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_AGENT", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CURSOR_CLI"])
	})

	t.Run("detects cursor cli from CURSOR_EXTENSION_HOST_ROLE=agent-exec", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_EXTENSION_HOST_ROLE", "agent-exec")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CURSOR_CLI"])
	})
}

func TestDetermineAgent_Gemini(t *testing.T) {
	t.Run("detects gemini from GEMINI_CLI", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("GEMINI_CLI", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GEMINI"])
	})
}

func TestDetermineAgent_Codex(t *testing.T) {
	for _, tc := range []struct{ env, val string }{
		{"CODEX_SANDBOX", "seatbelt"},
		{"CODEX_CI", "1"},
		{"CODEX_THREAD_ID", "thread-123"},
	} {
		tc := tc
		t.Run("detects codex from "+tc.env, func(t *testing.T) {
			clearAgentEnvs(t)
			t.Setenv(tc.env, tc.val)
			result, err := DetermineAgent()
			if err != nil {
				t.Fatal(err)
			}
			assertAgentResult(t, result, true, KnownAgents["CODEX"])
		})
	}
}

func TestDetermineAgent_Antigravity(t *testing.T) {
	t.Run("detects antigravity from ANTIGRAVITY_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("ANTIGRAVITY_AGENT", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["ANTIGRAVITY"])
	})
}

func TestDetermineAgent_AugmentCLI(t *testing.T) {
	t.Run("detects augment cli from AUGMENT_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AUGMENT_AGENT", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["AUGMENT_CLI"])
	})
}

func TestDetermineAgent_Opencode(t *testing.T) {
	for _, tc := range []struct{ env, val string }{
		{"OPENCODE_CLIENT", "opencode"},
		{"OPENCODE", "1"},
	} {
		tc := tc
		t.Run("detects opencode from "+tc.env, func(t *testing.T) {
			clearAgentEnvs(t)
			t.Setenv(tc.env, tc.val)
			result, err := DetermineAgent()
			if err != nil {
				t.Fatal(err)
			}
			assertAgentResult(t, result, true, KnownAgents["OPENCODE"])
		})
	}
}

func TestDetermineAgent_Goose(t *testing.T) {
	t.Run("detects goose from GOOSE_PROVIDER", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("GOOSE_PROVIDER", "anthropic")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["GOOSE"])
	})
}

func TestDetermineAgent_Junie(t *testing.T) {
	for _, tc := range []struct{ env, val string }{
		{"JUNIE_DATA", "/tmp/junie"},
		{"JUNIE_SHIM_PATH", "/tmp/junie/shim"},
	} {
		tc := tc
		t.Run("detects junie from "+tc.env, func(t *testing.T) {
			clearAgentEnvs(t)
			t.Setenv(tc.env, tc.val)
			result, err := DetermineAgent()
			if err != nil {
				t.Fatal(err)
			}
			assertAgentResult(t, result, true, KnownAgents["JUNIE"])
		})
	}
}

func TestDetermineAgent_PI(t *testing.T) {
	t.Run("returns no agent when PATH has no .pi/agent segment", func(t *testing.T) {
		clearAgentEnvs(t)
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("detects pi when PATH contains a .pi/agent segment", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("PATH", "/usr/bin:/home/me/.pi/agent/bin")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["PI"])
	})
}

func TestDetermineAgent_Kiro(t *testing.T) {
	t.Run("returns no agent when TERM_PROGRAM is not kiro", func(t *testing.T) {
		clearAgentEnvs(t)
		orig := isTTYFn
		isTTYFn = func() bool { return false }
		t.Cleanup(func() { isTTYFn = orig })

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("detects kiro when TERM_PROGRAM=kiro and no TTY", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("TERM_PROGRAM", "kiro")
		orig := isTTYFn
		isTTYFn = func() bool { return false }
		t.Cleanup(func() { isTTYFn = orig })

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["KIRO"])
	})

	t.Run("does not detect kiro when TERM_PROGRAM=kiro with a TTY (human at terminal)", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("TERM_PROGRAM", "kiro")
		orig := isTTYFn
		isTTYFn = func() bool { return true }
		t.Cleanup(func() { isTTYFn = orig })

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})
}

func TestDetermineAgent_Claude(t *testing.T) {
	t.Run("detects claude from CLAUDE_CODE", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDE_CODE", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CLAUDE"])
	})

	t.Run("detects claude from CLAUDECODE", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDECODE", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CLAUDE"])
	})
}

func TestDetermineAgent_Cowork(t *testing.T) {
	t.Run("detects claude (not cowork) when CLAUDE_CODE_IS_COWORK not set", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDECODE", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CLAUDE"])
	})

	t.Run("detects cowork when CLAUDE_CODE_IS_COWORK and CLAUDECODE both set", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDECODE", "1")
		t.Setenv("CLAUDE_CODE_IS_COWORK", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["COWORK"])
	})

	t.Run("detects cowork when CLAUDE_CODE_IS_COWORK and CLAUDE_CODE both set", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDE_CODE", "1")
		t.Setenv("CLAUDE_CODE_IS_COWORK", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["COWORK"])
	})

	t.Run("returns no agent when CLAUDE_CODE_IS_COWORK set without CLAUDECODE or CLAUDE_CODE", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CLAUDE_CODE_IS_COWORK", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})
}

func TestDetermineAgent_Devin(t *testing.T) {
	t.Run("returns no agent when /opt/.devin does not exist", func(t *testing.T) {
		clearAgentEnvs(t)
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})
}

func TestDetermineAgent_Replit(t *testing.T) {
	t.Run("detects replit from REPL_ID", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("REPL_ID", "1")
		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["REPLIT"])
	})
}

func TestDetermineAgent_Priority(t *testing.T) {
	t.Run("AI_AGENT takes highest priority", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "custom-priority")
		t.Setenv("CURSOR_TRACE_ID", "some-uuid")
		t.Setenv("CURSOR_AGENT", "1")
		t.Setenv("GEMINI_CLI", "1")
		t.Setenv("CODEX_SANDBOX", "seatbelt")
		t.Setenv("ANTIGRAVITY_AGENT", "1")
		t.Setenv("AUGMENT_AGENT", "1")
		t.Setenv("OPENCODE_CLIENT", "opencode")
		t.Setenv("CLAUDE_CODE", "1")
		t.Setenv("REPL_ID", "1")
		t.Setenv("COPILOT_MODEL", "gpt-5")
		t.Setenv("COPILOT_ALLOW_ALL", "true")
		t.Setenv("COPILOT_GITHUB_TOKEN", "ghp_xxx")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, "custom-priority")
	})

	t.Run("CURSOR_TRACE_ID takes priority over other agents (except AI_AGENT)", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_TRACE_ID", "some-uuid")
		t.Setenv("CURSOR_AGENT", "1")
		t.Setenv("GEMINI_CLI", "1")
		t.Setenv("CODEX_SANDBOX", "seatbelt")
		t.Setenv("ANTIGRAVITY_AGENT", "1")
		t.Setenv("AUGMENT_AGENT", "1")
		t.Setenv("OPENCODE_CLIENT", "opencode")
		t.Setenv("CLAUDE_CODE", "1")
		t.Setenv("REPL_ID", "1")
		t.Setenv("COPILOT_MODEL", "gpt-5")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CURSOR"])
	})

	t.Run("CURSOR_AGENT takes priority over remaining agents", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_AGENT", "1")
		t.Setenv("GEMINI_CLI", "1")
		t.Setenv("CODEX_SANDBOX", "seatbelt")
		t.Setenv("ANTIGRAVITY_AGENT", "1")
		t.Setenv("AUGMENT_AGENT", "1")
		t.Setenv("OPENCODE_CLIENT", "opencode")
		t.Setenv("CLAUDE_CODE", "1")
		t.Setenv("REPL_ID", "1")
		t.Setenv("COPILOT_MODEL", "gpt-5")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, KnownAgents["CURSOR_CLI"])
	})
}

func TestDetermineAgent_EdgeCases(t *testing.T) {
	t.Run("handles empty string values for environment variables", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "")
		t.Setenv("CURSOR_TRACE_ID", "")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("handles whitespace-only values for AI_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "   ")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, false, "")
	})

	t.Run("handles special characters in AI_AGENT value", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "my-custom-agent@v1.0")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, "my-custom-agent@v1.0")
	})

	t.Run("trims leading and trailing whitespace from AI_AGENT", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "  custom-agent  ")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		assertAgentResult(t, result, true, "custom-agent")
	})

	t.Run("isAgent is true when agent detected", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("AI_AGENT", "test-agent")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		if !result.IsAgent {
			t.Error("expected IsAgent true")
		}
	})

	t.Run("agent details available when detected", func(t *testing.T) {
		clearAgentEnvs(t)
		t.Setenv("CURSOR_TRACE_ID", "some-id")

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		if result.IsAgent && result.Agent.Name != KnownAgents["CURSOR"] {
			t.Errorf("got %q, want %q", result.Agent.Name, KnownAgents["CURSOR"])
		}
	})

	t.Run("agent is nil when not detected", func(t *testing.T) {
		clearAgentEnvs(t)

		result, err := DetermineAgent()
		if err != nil {
			t.Fatal(err)
		}
		if result.IsAgent {
			t.Error("expected IsAgent false")
		}
		if result.Agent != nil {
			t.Errorf("expected Agent nil, got %+v", result.Agent)
		}
	})
}
