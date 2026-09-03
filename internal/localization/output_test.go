package localization

import (
	"strings"
	"testing"
)

func TestRequireVietnameseRejectsHanCharactersInVietnameseLocale(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")

	err := RequireVietnamese("nội dung chương", "Cô mở cửa hầm. 地下室很冷。")

	if err == nil || !strings.Contains(err.Error(), "chữ Hán") {
		t.Fatalf("expected a Vietnamese language violation, got %v", err)
	}
}

func TestRequireVietnameseAllowsVietnameseDiacritics(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")

	if err := RequireVietnamese("nội dung chương", "Pháo đài di động tiến vào đường hầm."); err != nil {
		t.Fatalf("Vietnamese prose was rejected: %v", err)
	}
}

func TestRequireVietnamesePreservesChineseLocaleCompatibility(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "zh")

	if err := RequireVietnamese("chapter", "地下室很冷。"); err != nil {
		t.Fatalf("zh compatibility mode was rejected: %v", err)
	}
}

func TestForDisplayRemovesHanFromVietnameseRuntimeOutput(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")

	got := ForDisplay("✻ 查询上下文 chapter: 1")

	if ContainsHan(got) || !strings.Contains(got, "đã chặn chữ Hán") {
		t.Fatalf("runtime output was not sanitized: %q", got)
	}
}
