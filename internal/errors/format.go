package errors

import (
	"fmt"
	"strings"
)

func (e *Error) Error() string {
	switch e.Kind {
	case KindRequired:
		// Example:
		// CFG error: DATABASE_URL is required but was not set
		// field: DBURL
		return formatLines(
			fmt.Sprintf("CFG error: %s is required but was not set", e.Var),
			fmt.Sprintf("field: %s", e.Field),
		)

	case KindParse:
		// Example:
		// CFG error: PORT is invalid
		// field: Port
		// expected: int
		// received: "abc"
		return formatLines(
			fmt.Sprintf("CFG error: %s is invalid", e.Var),
			fmt.Sprintf("field: %s", e.Field),
			fmt.Sprintf("expected: %s", e.Expected),
			fmt.Sprintf(`received: %q`, e.Value),
		)

	case KindUnsupported:
		// Example:
		// CFG error: HOSTS has unsupported type
		// field: Hosts
		// expected: []int
		return formatLines(
			fmt.Sprintf("CFG error: %s has unsupported type", e.Var),
			fmt.Sprintf("field: %s", e.Field),
			fmt.Sprintf("expected: %s", e.Expected),
		)

	case KindInvalidTarget:
		// Keep this simple; it's a programmer error (wrong usage).
		if e.Cause != nil {
			return "CFG error: invalid target: " + e.Cause.Error()
		}
		return "CFG error: invalid target"

	default:
		// Fallback that still looks decent.
		msg := fmt.Sprintf("CFG error: %s", string(e.Kind))
		if e.Var != "" {
			msg += " (" + e.Var + ")"
		}
		return msg
	}
}

func formatLines(lines ...string) string {
	var out []string
	for _, s := range lines {
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return strings.Join(out, "\n")
}
