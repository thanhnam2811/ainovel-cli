package localization

import (
	"fmt"
	"strings"
	"unicode"
)

const blockedHanLabel = "[đã chặn chữ Hán]"

// Select returns the localized runtime label while preserving the upstream
// Chinese compatibility mode.
func Select(vietnamese, chinese string) string {
	if IsVietnamese() {
		return vietnamese
	}
	return chinese
}

// ContainsHan reports whether text contains a Han-script code point. In the
// Vietnamese locale, generated prose and its natural-language metadata must
// not silently fall back to Chinese.
func ContainsHan(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// RequireVietnamese is the persistence-boundary guard for model-authored
// output. The zh compatibility locale intentionally preserves upstream
// behavior.
func RequireVietnamese(field string, values ...string) error {
	if !IsVietnamese() {
		return nil
	}
	for _, value := range values {
		if ContainsHan(value) {
			return fmt.Errorf("%s chứa chữ Hán; hãy viết lại hoàn toàn bằng tiếng Việt", field)
		}
	}
	return nil
}

// ForDisplay is the last-resort runtime boundary. Known UI labels are
// localized at their source; any unexpected Han text from a provider or a
// stale upstream message is redacted before it reaches TUI/headless output.
func ForDisplay(text string) string {
	if !IsVietnamese() || !ContainsHan(text) {
		return text
	}
	var out strings.Builder
	inHan := false
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			if !inHan {
				out.WriteString(blockedHanLabel)
			}
			inHan = true
			continue
		}
		inHan = false
		out.WriteRune(r)
	}
	return out.String()
}
