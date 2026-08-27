package assets

import (
	"strings"
	"testing"
)

func TestApplyLocaleVietnamese(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "vi-VN"); err != nil {
		t.Fatalf("ApplyLocale(vi-VN): %v", err)
	}

	if !strings.Contains(b.Prompts.Writer, "{{VOICE}}") {
		t.Fatal("Vietnamese writer prompt lost {{VOICE}} placeholder")
	}
	if !strings.Contains(b.Prompts.Writer, "## Ngôn ngữ đầu ra") {
		t.Fatal("Vietnamese writer prompt missing language directive")
	}
	if !strings.Contains(b.Prompts.Writer, "## Quy trình thực thi") {
		t.Fatal("Vietnamese writer overlay was not loaded")
	}
	if !strings.Contains(b.Voice, "## Tiêu chuẩn viết") {
		t.Fatal("Vietnamese voice overlay was not loaded")
	}
	if !strings.Contains(b.References.AntiAITone, "# Tiêu chí giảm") {
		t.Fatal("Vietnamese anti-AI reference was not loaded")
	}
	if !strings.Contains(b.Styles["default"], "## Phong cách viết mặc định") {
		t.Fatal("Vietnamese default style was not loaded")
	}

	// Prompts without a translated overlay must still preserve the latest
	// upstream contract while being instructed to produce Vietnamese output.
	if !strings.Contains(b.Prompts.ImportSegment, "## Ngôn ngữ đầu ra") {
		t.Fatal("fallback prompt missing Vietnamese output directive")
	}
}

func TestApplyLocaleUpstreamModeIsNoOp(t *testing.T) {
	b := Load("default", LoadOptions{})
	beforeWriter := b.Prompts.Writer
	beforeVoice := b.Voice

	if err := ApplyLocale(&b, "zh"); err != nil {
		t.Fatalf("ApplyLocale(zh): %v", err)
	}
	if b.Prompts.Writer != beforeWriter || b.Voice != beforeVoice {
		t.Fatal("zh locale must preserve upstream assets byte-for-byte")
	}
}

func TestApplyLocaleRejectsUnknownLocale(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "en"); err == nil {
		t.Fatal("expected unsupported locale error")
	}
}

func TestApplyLocaleVietnameseIsIdempotent(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "vi"); err != nil {
		t.Fatal(err)
	}
	first := b.Prompts.ImportSegment
	if err := ApplyLocale(&b, "vi"); err != nil {
		t.Fatal(err)
	}
	if b.Prompts.ImportSegment != first {
		t.Fatal("applying Vietnamese locale twice changed a fallback prompt")
	}
}
