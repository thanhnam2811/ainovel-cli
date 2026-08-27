package assets

import (
	"strings"
	"testing"
)

func TestVietnameseFunctionPromptsPreserveContracts(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "vi"); err != nil {
		t.Fatal(err)
	}

	assertFunctionPromptTokens(t, "import-segment", b.Prompts.ImportSegment,
		"owned_start", "owned_end", "units", "user_guidance",
		"unit_id", "chapter", "group", "front_matter", "back_matter", "uncertain", "boundaries",
	)
	assertFunctionPromptTokens(t, "import-analyze", b.Prompts.ImportAnalyze,
		"chapters", "hook_type", "crisis", "mystery", "desire", "emotion", "choice",
		"dominant_strand", "quest", "fire", "constellation", "foreshadow_updates", "plant", "advance", "resolve",
	)
	assertFunctionPromptTokens(t, "import-synthesize", b.Prompts.ImportSynthesize,
		"planning_tier", "short", "mid", "long", "story_status", "open", "closed", "uncertain",
		"compass.ending_direction", "synopsis", "premise", "structure",
	)
	assertFunctionPromptTokens(t, "import-range", b.Prompts.ImportRange,
		"start_chapter", "end_chapter", "plot", "characters", "world_facts", "opened_threads", "resolved_threads",
	)
	assertFunctionPromptTokens(t, "revision-analyze", b.Prompts.RevisionAnalyze,
		"facts", "revised_content", "changed_excerpt", "previous_facts", "style_delta",
		"story_changed", "outline_impact", "downstream_issues",
	)
	assertFunctionPromptTokens(t, "simulation-source", b.Prompts.SimulationSource,
		"cấu trúc", "nhịp", "không", // proves localized overlay rather than fallback-only behavior
	)
	assertFunctionPromptTokens(t, "simulation-merge", b.Prompts.SimulationMerge,
		"source_reports", "profile", "cấu trúc",
	)
	assertFunctionPromptTokens(t, "arbiter-plan-start", b.Prompts.ArbiterPlanStart,
		"requirement", "style", "architect_long", "architect_short", "save_foundation", "audit_foundation", "foundation_ready=true", "complete_book",
	)
	assertFunctionPromptTokens(t, "arbiter-intervention", b.Prompts.ArbiterIntervention,
		"intervention", "facts", "answer", "rules", "hold", "reopen", "dispatch",
		"has_advance_hold", "target_chapter", "rewrites_drained", "advance_mode", "requires_change=true",
		"dynamic_planning", "outlined_chapters", "update_compass", "append_volume", "expand_arc", "revise_outline",
		"architect_long", "architect_short", "editor", "writer", "phase = complete", "reopen_count", "recent_decisions",
	)
	assertFunctionPromptTokens(t, "arbiter-failure", b.Prompts.ArbiterFailure,
		"worker_failure", "deadlock", "reroute", "retry", "abort", "dispatch", "repeats",
		"foundation_missing", "architect_long", "architect_short", "writer", "editor",
	)
}

func assertFunctionPromptTokens(t *testing.T, name, prompt string, tokens ...string) {
	t.Helper()
	if !strings.Contains(prompt, "## Ngôn ngữ đầu ra") {
		t.Fatalf("%s missing Vietnamese output directive", name)
	}
	for _, token := range tokens {
		if !strings.Contains(prompt, token) {
			t.Errorf("%s lost contract/localization token %q", name, token)
		}
	}
}
