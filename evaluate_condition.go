package detectagent

import (
	"os"
	"regexp"
)

// Condition is a node in a condition tree loaded from agents.json.
type Condition struct {
	Type       string      `json:"type"`
	Name       string      `json:"name,omitempty"`
	Value      string      `json:"value,omitempty"`
	Pattern    string      `json:"pattern,omitempty"`
	Path       string      `json:"path,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// isTTYFn is swappable in tests.
var isTTYFn = func() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// EvaluateCondition evaluates a condition tree. anyOf/allOf are combinators;
// the rest are leaves.
func EvaluateCondition(c Condition) (bool, error) {
	switch c.Type {
	case "env_set":
		return os.Getenv(c.Name) != "", nil

	case "env_value":
		val, ok := os.LookupEnv(c.Name)
		if !ok {
			return false, nil
		}
		return val == c.Value, nil

	case "env_matches":
		val := os.Getenv(c.Name)
		if val == "" {
			return false, nil
		}
		re, err := regexp.Compile(c.Pattern)
		if err != nil {
			// A malformed pattern must not propagate — fail closed.
			return false, nil
		}
		return re.MatchString(val), nil

	case "no_tty":
		return !isTTYFn(), nil

	case "file_exists":
		_, err := os.Stat(c.Path)
		return err == nil, nil

	case "anyOf":
		for _, sub := range c.Conditions {
			ok, err := EvaluateCondition(sub)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil

	case "allOf":
		for _, sub := range c.Conditions {
			ok, err := EvaluateCondition(sub)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil

	default:
		return false, nil
	}
}
