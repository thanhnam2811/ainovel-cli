package assets

import (
	"os"
	"path/filepath"
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

func TestApplyLocalePreservesStyleOverridePrecedence(t *testing.T) {
	home := t.TempDir()
	book := t.TempDir()
	mustWriteTestFile(t, filepath.Join(home, "voice.md"), "GLOBAL VOICE")
	mustWriteTestFile(t, filepath.Join(book, "voice.md"), "BOOK VOICE")
	mustWriteTestFile(t, filepath.Join(home, "anti-ai-tone.md"), "GLOBAL ANTI AI")
	mustWriteTestFile(t, filepath.Join(book, "anti-ai-tone.md"), "BOOK ANTI AI")
	mustWriteTestFile(t, filepath.Join(home, "styles", "default.md"), "GLOBAL STYLE")
	mustWriteTestFile(t, filepath.Join(book, "styles", "default.md"), "BOOK STYLE")

	opts := LoadOptions{HomeStyleDir: home, BookStyleDir: book}
	b := Load("default", opts)
	if err := ApplyLocale(&b, "vi", opts); err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(b.Voice, "## Tiêu chuẩn viết") || !strings.Contains(b.Voice, "GLOBAL VOICE") || !strings.Contains(b.Voice, "BOOK VOICE") {
		t.Fatalf("voice precedence was not preserved: %q", b.Voice)
	}
	if !strings.HasPrefix(b.References.AntiAITone, "# Tiêu chí giảm") || !strings.Contains(b.References.AntiAITone, "GLOBAL ANTI AI") || !strings.Contains(b.References.AntiAITone, "BOOK ANTI AI") {
		t.Fatalf("anti-AI precedence was not preserved: %q", b.References.AntiAITone)
	}
	if got := strings.TrimSpace(b.Styles["default"]); got != "BOOK STYLE" {
		t.Fatalf("book style must win, got %q", got)
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

func mustWriteTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
