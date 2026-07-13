package detectagent

import (
	"os"
	"testing"
)

func TestEvaluateCondition_EnvSet(t *testing.T) {
	t.Run("true when variable has a non-empty value", func(t *testing.T) {
		t.Setenv("SOME_VAR", "1")
		got, err := EvaluateCondition(Condition{Type: "env_set", Name: "SOME_VAR"})
		assertBool(t, err, got, true)
	})

	t.Run("false when variable is unset", func(t *testing.T) {
		os.Unsetenv("SOME_VAR")
		got, err := EvaluateCondition(Condition{Type: "env_set", Name: "SOME_VAR"})
		assertBool(t, err, got, false)
	})

	t.Run("false when variable is set to empty string", func(t *testing.T) {
		t.Setenv("SOME_VAR", "")
		got, err := EvaluateCondition(Condition{Type: "env_set", Name: "SOME_VAR"})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_EnvValue(t *testing.T) {
	t.Run("true when value matches exactly", func(t *testing.T) {
		t.Setenv("ROLE", "agent-exec")
		got, err := EvaluateCondition(Condition{Type: "env_value", Name: "ROLE", Value: "agent-exec"})
		assertBool(t, err, got, true)
	})

	t.Run("false when value differs", func(t *testing.T) {
		t.Setenv("ROLE", "something-else")
		got, err := EvaluateCondition(Condition{Type: "env_value", Name: "ROLE", Value: "agent-exec"})
		assertBool(t, err, got, false)
	})

	t.Run("false when variable is unset", func(t *testing.T) {
		os.Unsetenv("ROLE")
		got, err := EvaluateCondition(Condition{Type: "env_value", Name: "ROLE", Value: "agent-exec"})
		assertBool(t, err, got, false)
	})

	t.Run("does not treat empty target as matching unset variable", func(t *testing.T) {
		// os.LookupEnv distinguishes unset from set-to-empty; an unset var
		// must not match value "".
		os.Unsetenv("ROLE")
		got, err := EvaluateCondition(Condition{Type: "env_value", Name: "ROLE", Value: ""})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_FileExists(t *testing.T) {
	t.Run("true when path exists", func(t *testing.T) {
		dir := t.TempDir()
		got, err := EvaluateCondition(Condition{Type: "file_exists", Path: dir})
		assertBool(t, err, got, true)
	})

	t.Run("false when path does not exist", func(t *testing.T) {
		got, err := EvaluateCondition(Condition{Type: "file_exists", Path: "/nonexistent/path/abc123"})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_EnvMatches(t *testing.T) {
	t.Run("true when value matches the pattern", func(t *testing.T) {
		t.Setenv("PATH", `/usr/bin:/home/me/.pi/agent/bin`)
		got, err := EvaluateCondition(Condition{Type: "env_matches", Name: "PATH", Pattern: `\.pi[\\/]agent`})
		assertBool(t, err, got, true)
	})

	t.Run("matches a backslash path separator too", func(t *testing.T) {
		t.Setenv("PATH", `C:\Users\me\.pi\agent\bin`)
		got, err := EvaluateCondition(Condition{Type: "env_matches", Name: "PATH", Pattern: `\.pi[\\/]agent`})
		assertBool(t, err, got, true)
	})

	t.Run("false when value does not match", func(t *testing.T) {
		t.Setenv("PATH", "/usr/bin:/usr/local/bin")
		got, err := EvaluateCondition(Condition{Type: "env_matches", Name: "PATH", Pattern: `\.pi[\\/]agent`})
		assertBool(t, err, got, false)
	})

	t.Run("false when variable is unset or empty", func(t *testing.T) {
		t.Setenv("TERM_PROGRAM", "")
		got, err := EvaluateCondition(Condition{Type: "env_matches", Name: "TERM_PROGRAM", Pattern: "kiro"})
		assertBool(t, err, got, false)
	})

	t.Run("false when pattern is invalid rather than erroring", func(t *testing.T) {
		t.Setenv("TERM_PROGRAM", "kiro")
		got, err := EvaluateCondition(Condition{Type: "env_matches", Name: "TERM_PROGRAM", Pattern: "("})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_NoTTY(t *testing.T) {
	t.Run("true when stdout is not a TTY", func(t *testing.T) {
		orig := isTTYFn
		isTTYFn = func() bool { return false }
		t.Cleanup(func() { isTTYFn = orig })

		got, err := EvaluateCondition(Condition{Type: "no_tty"})
		assertBool(t, err, got, true)
	})

	t.Run("false when stdout is a TTY", func(t *testing.T) {
		orig := isTTYFn
		isTTYFn = func() bool { return true }
		t.Cleanup(func() { isTTYFn = orig })

		got, err := EvaluateCondition(Condition{Type: "no_tty"})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_AnyOf(t *testing.T) {
	t.Run("true when at least one child is true", func(t *testing.T) {
		t.Setenv("B", "1")
		os.Unsetenv("A")
		got, err := EvaluateCondition(Condition{
			Type: "anyOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, true)
	})

	t.Run("false when all children are false", func(t *testing.T) {
		os.Unsetenv("A")
		os.Unsetenv("B")
		got, err := EvaluateCondition(Condition{
			Type: "anyOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, false)
	})

	t.Run("false for an empty condition list", func(t *testing.T) {
		got, err := EvaluateCondition(Condition{Type: "anyOf", Conditions: []Condition{}})
		assertBool(t, err, got, false)
	})

	t.Run("true when a later child is satisfied but an earlier one is not", func(t *testing.T) {
		t.Setenv("A", "")
		t.Setenv("B", "1")
		got, err := EvaluateCondition(Condition{
			Type: "anyOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, true)
	})
}

func TestEvaluateCondition_AllOf(t *testing.T) {
	t.Run("true when all children are true", func(t *testing.T) {
		t.Setenv("A", "1")
		t.Setenv("B", "1")
		got, err := EvaluateCondition(Condition{
			Type: "allOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, true)
	})

	t.Run("false when any child is false", func(t *testing.T) {
		t.Setenv("A", "1")
		os.Unsetenv("B")
		got, err := EvaluateCondition(Condition{
			Type: "allOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, false)
	})

	t.Run("true for an empty condition list (vacuous truth)", func(t *testing.T) {
		got, err := EvaluateCondition(Condition{Type: "allOf", Conditions: []Condition{}})
		assertBool(t, err, got, true)
	})

	t.Run("false when a later child is unsatisfied", func(t *testing.T) {
		t.Setenv("A", "1")
		t.Setenv("B", "")
		got, err := EvaluateCondition(Condition{
			Type: "allOf",
			Conditions: []Condition{
				{Type: "env_set", Name: "A"},
				{Type: "env_set", Name: "B"},
			},
		})
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_Nesting(t *testing.T) {
	coworkCond := Condition{
		Type: "allOf",
		Conditions: []Condition{
			{Type: "env_set", Name: "CLAUDE_CODE_IS_COWORK"},
			{
				Type: "anyOf",
				Conditions: []Condition{
					{Type: "env_set", Name: "CLAUDECODE"},
					{Type: "env_set", Name: "CLAUDE_CODE"},
				},
			},
		},
	}

	t.Run("evaluates allOf(env_set, anyOf(...)) tree", func(t *testing.T) {
		t.Setenv("CLAUDE_CODE_IS_COWORK", "1")
		t.Setenv("CLAUDECODE", "1")
		got, err := EvaluateCondition(coworkCond)
		assertBool(t, err, got, true)
	})

	t.Run("false when the nested anyOf has no satisfied child", func(t *testing.T) {
		t.Setenv("CLAUDE_CODE_IS_COWORK", "1")
		t.Setenv("CLAUDECODE", "")
		t.Setenv("CLAUDE_CODE", "")
		got, err := EvaluateCondition(coworkCond)
		assertBool(t, err, got, false)
	})
}

func TestEvaluateCondition_UnknownType(t *testing.T) {
	t.Run("returns false for an unrecognized type", func(t *testing.T) {
		got, err := EvaluateCondition(Condition{Type: "not_a_real_type"})
		assertBool(t, err, got, false)
	})
}

func assertBool(t *testing.T, err error, got bool, want bool) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
