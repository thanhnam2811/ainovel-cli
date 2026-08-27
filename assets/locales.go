package assets

import (
	"embed"
	"fmt"
	"strings"

	"github.com/voocel/ainovel-cli/internal/tools"
)

// localeFS contains downstream localization overlays. Keeping translations under
// assets/locales avoids translating engine code and keeps upstream merges small.
// Missing localized files intentionally fall back to the upstream built-ins.
//
//go:embed locales
var localeFS embed.FS

const viLanguageDirective = `## Ngôn ngữ đầu ra

Mọi nội dung văn xuôi, giải thích, nhận xét và dữ liệu ngôn ngữ tự nhiên do bạn tạo phải dùng tiếng Việt tự nhiên, rõ nghĩa và nhất quán với bối cảnh tác phẩm.

Giữ nguyên tuyệt đối tên tool, tên field JSON, enum, protocol marker, placeholder và khóa máy đọc được. Không dịch các identifier kỹ thuật như novel_context, plan_chapter, draft_chapter, commit_chapter, save_review, {{VOICE}} hoặc các giá trị enum được schema yêu cầu.`

const viSimulationGuidance = `## Hướng dẫn mô phỏng phong cách

Khi novel_context chứa simulation_profile trong planning_memory hoặc working_memory, hãy xem đó là định hướng phong cách của tác phẩm hiện tại. Đọc và vận dụng các mục style, lexicon, plot_design, hook_design, pacing_density, reader_engagement và role_guidance.

Chỉ học cấu trúc, nhịp, cách đặt móc, mật độ thông tin và kỹ thuật giữ người đọc. Không sao chép câu chữ, nhân vật, địa danh, thiết lập riêng hoặc tình tiết cố định của nguồn tham chiếu. Nếu simulation_profile xung đột với yêu cầu trực tiếp của người dùng, ưu tiên yêu cầu của người dùng.`

// ApplyLocale overlays localized writing assets onto an already loaded upstream
// bundle. It stays style-agnostic for existing callers; production startup should
// use ApplyLocaleForStyle so genre-specific references can also be localized.
func ApplyLocale(b *Bundle, locale string, loadOpts ...LoadOptions) error {
	return ApplyLocaleForStyle(b, locale, "", loadOpts...)
}

// ApplyLocaleForStyle overlays localized assets while preserving the selected
// upstream style/genre. Missing localized files keep the newest upstream asset.
// When LoadOptions is provided, localized built-ins keep upstream precedence:
// localized builtin < ~/.ainovel/style < <book>/style.
func ApplyLocaleForStyle(b *Bundle, locale, style string, loadOpts ...LoadOptions) error {
	if b == nil {
		return fmt.Errorf("assets: nil bundle")
	}

	var opts LoadOptions
	if len(loadOpts) > 0 {
		opts = loadOpts[0]
	}

	switch normalizeLocale(locale) {
	case "":
		return nil
	case "vi":
		applyVietnameseLocale(b, opts, style)
		return nil
	default:
		return fmt.Errorf("assets: unsupported locale %q (supported: vi, zh)", locale)
	}
}

func normalizeLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	locale = strings.ReplaceAll(locale, "_", "-")
	switch locale {
	case "", "zh", "zh-cn", "zh-hans":
		return ""
	case "vi", "vi-vn":
		return "vi"
	default:
		return locale
	}
}

func applyVietnameseLocale(b *Bundle, opts LoadOptions, style string) {
	// Core prompts may be translated independently. A missing file keeps the
	// latest upstream prompt, so protocol changes arrive immediately on sync.
	b.Prompts.ArchitectShort = localizedCorePrompt("vi", "architect-short.md", b.Prompts.ArchitectShort)
	b.Prompts.ArchitectLong = localizedCorePrompt("vi", "architect-long.md", b.Prompts.ArchitectLong)
	b.Prompts.Writer = localizedCorePrompt("vi", "writer.md", b.Prompts.Writer)
	b.Prompts.Editor = localizedCorePrompt("vi", "editor.md", b.Prompts.Editor)

	b.Prompts.ImportSegment = localizedPrompt("vi", "import-segment.md", b.Prompts.ImportSegment)
	b.Prompts.ImportAnalyze = localizedPrompt("vi", "import-analyze.md", b.Prompts.ImportAnalyze)
	b.Prompts.ImportSynthesize = localizedPrompt("vi", "import-synthesize.md", b.Prompts.ImportSynthesize)
	b.Prompts.ImportRange = localizedPrompt("vi", "import-range.md", b.Prompts.ImportRange)
	b.Prompts.SimulationSource = localizedPrompt("vi", "simulation-source.md", b.Prompts.SimulationSource)
	b.Prompts.SimulationMerge = localizedPrompt("vi", "simulation-merge.md", b.Prompts.SimulationMerge)
	b.Prompts.RevisionAnalyze = localizedPrompt("vi", "revision-analyze.md", b.Prompts.RevisionAnalyze)
	b.Prompts.ArbiterPlanStart = localizedPrompt("vi", "arbiter-plan-start.md", b.Prompts.ArbiterPlanStart)
	b.Prompts.ArbiterIntervention = localizedPrompt("vi", "arbiter-intervention.md", b.Prompts.ArbiterIntervention)
	b.Prompts.ArbiterFailure = localizedPrompt("vi", "arbiter-failure.md", b.Prompts.ArbiterFailure)

	// Voice and anti-AI tone are appendable upstream assets. Rebuild them from
	// the localized builtin and then re-apply user overrides in upstream order.
	if raw, ok := readLocale("vi", "voice.md"); ok {
		b.Voice = resolveAppendable(raw, "voice.md", opts)
	}
	if raw, ok := readLocale("vi", "references/anti-ai-tone.md"); ok {
		b.References.AntiAITone = resolveAppendable(raw, "anti-ai-tone.md", opts)
	}

	// Styles use whole-file replacement. Apply localized built-ins first, then
	// re-run the exact upstream global/book overlays so user styles always win.
	if b.Styles == nil {
		b.Styles = make(map[string]string)
	}
	for name := range b.Styles {
		if raw, ok := readLocale("vi", "styles/"+name+".md"); ok {
			b.Styles[name] = raw
		}
	}
	overlayStyles(b.Styles, opts.HomeStyleDir)
	overlayStyles(b.Styles, opts.BookStyleDir)

	applyLocalizedReferences(&b.References, "vi")
	applyLocalizedGenreReferences(&b.References, "vi", style, opts)
}

func applyLocalizedReferences(refs *tools.References, locale string) {
	if refs == nil {
		return
	}
	setLocalized(&refs.ChapterGuide, locale, "references/chapter-guide.md")
	setLocalized(&refs.HookTechniques, locale, "references/hook-techniques.md")
	setLocalized(&refs.QualityChecklist, locale, "references/quality-checklist.md")
	setLocalized(&refs.OutlineTemplate, locale, "references/outline-template.md")
	setLocalized(&refs.CharacterTemplate, locale, "references/character-template.md")
	setLocalized(&refs.ChapterTemplate, locale, "references/chapter-template.md")
	setLocalized(&refs.Consistency, locale, "references/consistency.md")
	setLocalized(&refs.ContentExpansion, locale, "references/content-expansion.md")
	setLocalized(&refs.DialogueWriting, locale, "references/dialogue-writing.md")
	setLocalized(&refs.LongformPlanning, locale, "references/longform-planning.md")
	setLocalized(&refs.Differentiation, locale, "references/differentiation.md")
	// AntiAITone is intentionally excluded here because it has appendable user
	// overrides and is rebuilt above with resolveAppendable.
}

func applyLocalizedGenreReferences(refs *tools.References, locale, style string, opts LoadOptions) {
	if refs == nil {
		return
	}
	style = strings.ToLower(strings.TrimSpace(style))
	if style == "" || style == "default" || !styleNameRe.MatchString(style) {
		return
	}

	base := "references/genres/" + style + "/"
	if raw, ok := readLocale(locale, base+"style-references.md"); ok {
		refs.StyleReference = raw
		// Upstream allows user replacement of genre style references. Reapply
		// those exact layers after the localized builtin so user intent wins.
		rel := "genres/" + style + "/style-references.md"
		for _, dir := range []string{opts.HomeStyleDir, opts.BookStyleDir} {
			if override := readOverride(dir, rel); override != "" {
				refs.StyleReference = override
			}
		}
	}
	if raw, ok := readLocale(locale, base+"arc-templates.md"); ok {
		refs.ArcTemplates = raw
	}
}

func setLocalized(dst *string, locale, path string) {
	if raw, ok := readLocale(locale, path); ok {
		*dst = raw
	}
}

func localizedCorePrompt(locale, name, fallback string) string {
	if raw, ok := readLocale(locale, "prompts/"+name); ok {
		if name == "writer.md" && !strings.Contains(raw, voicePlaceholder) {
			// Keep a usable upstream writer instead of accepting a broken overlay.
			return appendDirective(fallback, viLanguageDirective)
		}
		return appendDirective(raw+"\n\n"+viSimulationGuidance, viLanguageDirective)
	}
	return appendDirective(fallback, viLanguageDirective)
}

func localizedPrompt(locale, name, fallback string) string {
	if raw, ok := readLocale(locale, "prompts/"+name); ok {
		return appendDirective(raw, viLanguageDirective)
	}
	return appendDirective(fallback, viLanguageDirective)
}

func appendDirective(base, directive string) string {
	if strings.Contains(base, directive) {
		return base
	}
	return strings.TrimSpace(base) + "\n\n" + directive
}

func readLocale(locale, path string) (string, bool) {
	data, err := localeFS.ReadFile("locales/" + locale + "/" + path)
	if err != nil {
		return "", false
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return "", false
	}
	return raw, true
}
