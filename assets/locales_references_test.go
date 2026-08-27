package assets

import (
	"strings"
	"testing"
)

func TestVietnameseReferencePackIsLocalized(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocale(&b, "vi"); err != nil {
		t.Fatal(err)
	}

	checks := []struct {
		name string
		text string
		want []string
	}{
		{"chapter-guide", b.References.ChapterGuide, []string{"# Hướng dẫn viết chương", "chapter_contract", "user_rules.preferences"}},
		{"hook-techniques", b.References.HookTechniques, []string{"# Kỹ thuật tạo hook", "cliffhanger", "Foreshadow"}},
		{"quality-checklist", b.References.QualityChecklist, []string{"# Checklist chất lượng chương", "user_rules", "anti_ai_tone"}},
		{"outline-template", b.References.OutlineTemplate, []string{"# Template quy hoạch outline", "layered_outline", "estimated_chapters"}},
		{"character-template", b.References.CharacterTemplate, []string{"# Hồ sơ nhân vật", "Alias / danh hiệu"}},
		{"chapter-template", b.References.ChapterTemplate, []string{"# Chương [X]", "working_memory.user_rules.preferences"}},
		{"consistency", b.References.Consistency, []string{"# Cơ chế bảo đảm continuity", "novel_context"}},
		{"content-expansion", b.References.ContentExpansion, []string{"# Kỹ thuật mở rộng nội dung", "user_rules.preferences"}},
		{"dialogue-writing", b.References.DialogueWriting, []string{"# Quy chuẩn viết hội thoại", "Dấu câu và format tiếng Việt", "Subtext"}},
		{"longform-planning", b.References.LongformPlanning, []string{"# Tham chiếu quy hoạch truyện dài", "Story engine", "estimated_chapters"}},
		{"differentiation", b.References.Differentiation, []string{"# Tham chiếu thiết kế khác biệt hóa", "Năm chiều khác biệt hóa"}},
		{"anti-ai-tone", b.References.AntiAITone, []string{"# Tiêu chí giảm"}},
	}

	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			for _, token := range tc.want {
				if !strings.Contains(tc.text, token) {
					t.Errorf("localized %s missing token %q", tc.name, token)
				}
			}
		})
	}
}
