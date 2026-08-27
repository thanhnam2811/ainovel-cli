package localization

import (
	"os"
	"strings"
)

// Current returns the runtime natural-language locale used by downstream
// LLM-facing helpers that are not part of the embedded asset Bundle.
// The Vietnamese fork defaults to vi; zh keeps upstream behavior for control
// runs and debugging. Unknown values are preserved so callers can safely
// choose their upstream fallback instead of guessing.
func Current() string {
	return Normalize(os.Getenv("AINOVEL_LOCALE"))
}

// Normalize canonicalizes locale aliases used by this fork.
func Normalize(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	locale = strings.ReplaceAll(locale, "_", "-")
	switch locale {
	case "", "vi", "vi-vn":
		return "vi"
	case "zh", "zh-cn", "zh-hans":
		return "zh"
	default:
		return locale
	}
}

func IsVietnamese() bool { return Current() == "vi" }
