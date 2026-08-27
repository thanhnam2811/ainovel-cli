package eval

import (
	"strings"
	"testing"

	"github.com/voocel/ainovel-cli/assets"
)

func TestLoadEvalBundleDefaultsToVietnamese(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "")
	bundle, err := loadEvalBundle("default")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(bundle.Prompts.Writer, "## Ngôn ngữ đầu ra") {
		t.Fatalf("eval baseline did not receive Vietnamese locale")
	}
	if !strings.Contains(bundle.References.ChapterGuide, "# Hướng dẫn viết chương") {
		t.Fatalf("eval reference pack did not receive Vietnamese locale")
	}
}

func TestLoadEvalBundleZhMatchesUpstreamBuiltin(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "zh")
	got, err := loadEvalBundle("fantasy")
	if err != nil {
		t.Fatal(err)
	}
	want := assets.Load("fantasy", assets.LoadOptions{})
	if got.Prompts.Writer != want.Prompts.Writer {
		t.Fatalf("zh eval mode changed upstream writer prompt")
	}
	if got.References.StyleReference != want.References.StyleReference {
		t.Fatalf("zh eval mode changed upstream genre reference")
	}
}

func TestEvalVariantRemainsFinalPromptLayer(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "")
	bundle, err := loadEvalBundle("default")
	if err != nil {
		t.Fatal(err)
	}
	const variant = "VARIANT_WRITER_PROMPT"
	if err := applyVariant(&bundle, map[string]string{"writer.md": variant}); err != nil {
		t.Fatal(err)
	}
	want := assets.WithSimulationGuidance(variant, "writer")
	if bundle.Prompts.Writer != want {
		t.Fatalf("variant must override localized baseline, got %q, want %q", bundle.Prompts.Writer, want)
	}
}
