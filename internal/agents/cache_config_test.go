package agents

import (
	"testing"

	"github.com/voocel/ainovel-cli/internal/bootstrap"
)

func TestCacheLastMessageForOpenAIResponsesUsesPromptCacheKeyOnly(t *testing.T) {
	cfg := bootstrap.Config{Providers: map[string]bootstrap.ProviderConfig{
		"zen": {Type: "openai", API: "responses"},
	}}

	if got := cacheLastMessageFor(cfg, "zen"); got != "" {
		t.Fatalf("OpenAI Responses must not receive block cache breakpoints, got %q", got)
	}
}

func TestCacheLastMessageForAnthropicKeepsRollingBreakpoint(t *testing.T) {
	cfg := bootstrap.Config{Providers: map[string]bootstrap.ProviderConfig{
		"claude": {Type: "anthropic"},
	}}

	if got := cacheLastMessageFor(cfg, "claude"); got != "ephemeral" {
		t.Fatalf("Anthropic should keep rolling cache breakpoints, got %q", got)
	}
}
