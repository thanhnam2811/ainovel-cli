package assets

import (
	"strings"
	"testing"
)

func TestVietnameseCoreAgentPromptsPreserveProtocol(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "vi"); err != nil {
		t.Fatal(err)
	}

	assertPromptTokens(t, "architect-short", b.Prompts.ArchitectShort,
		"save_book", "save_foundation", "revise_outline", "audit_foundation",
		"foundation_ready", "remaining", `type=\"outline\"`,
	)
	assertPromptTokens(t, "architect-long", b.Prompts.ArchitectLong,
		"save_book", "save_foundation", "revise_outline", "audit_foundation",
		"layered_outline", "update_compass", "append_volume", "complete_book",
		"expand_arc", "final_volume", "open_threads", "completion_signals",
	)
	assertPromptTokens(t, "editor", b.Prompts.Editor,
		"read_chapter", "save_review", "save_arc_summary", "save_volume_summary",
		"requires_change", "rule_violations", "forbidden_chars", "fatigue_words",
		"accept", "polish", "rewrite",
	)
}

func assertPromptTokens(t *testing.T, name, prompt string, tokens ...string) {
	t.Helper()
	if !strings.Contains(prompt, "## Ngôn ngữ đầu ra") {
		t.Fatalf("%s is missing Vietnamese output directive", name)
	}
	for _, token := range tokens {
		if !strings.Contains(prompt, token) {
			t.Errorf("%s lost protocol token %q", name, token)
		}
	}
}
