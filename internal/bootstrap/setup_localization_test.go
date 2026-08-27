package bootstrap

import (
	"strings"
	"testing"
)

func TestSetupCopyDefaultsToVietnamese(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "")
	copy := activeSetupCopy()
	if copy.noConfig != "Không tìm thấy tệp cấu hình, bắt đầu thiết lập..." {
		t.Fatalf("unexpected default setup locale: %q", copy.noConfig)
	}
	if strings.Contains(copy.providerTitle, "选择") || strings.Contains(copy.selectHelp, "确认") {
		t.Fatalf("Vietnamese setup copy leaked Chinese UI: %#v", copy)
	}
}

func TestSetupCopyPreservesChineseControlLocale(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "zh-CN")
	copy := activeSetupCopy()
	if copy.noConfig != "未检测到配置文件，开始初始化设置..." {
		t.Fatalf("unexpected Chinese control copy: %q", copy.noConfig)
	}
	if copy.providerTitle != "[1/4] 选择 Provider" {
		t.Fatalf("unexpected Chinese provider title: %q", copy.providerTitle)
	}
}

func TestSetupViewsUseActiveLocaleHelp(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")
	selectView := (setupSelectModel{title: activeSetupCopy().providerTitle, items: setupProviders}).View()
	if !strings.Contains(selectView, "Enter xác nhận") || strings.Contains(selectView, "Enter 确认") {
		t.Fatalf("selector help was not localized: %q", selectView)
	}

	inputView := (setupInputModel{label: activeSetupCopy().modelTitle, placeholder: activeSetupCopy().modelExample}).View()
	if !strings.Contains(inputView, "Esc hủy") || strings.Contains(inputView, "Esc 取消") {
		t.Fatalf("input help was not localized: %q", inputView)
	}
}

func TestAPITypeDisplayLabelsAreLocalizedWithoutChangingIDs(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")
	items := apiTypeOptions()
	wantNames := []string{"openai", "anthropic", "gemini"}
	for i, item := range items {
		if item.name != wantNames[i] {
			t.Fatalf("protocol identifier changed: got %q want %q", item.name, wantNames[i])
		}
		if !strings.Contains(item.label, "tương thích") {
			t.Fatalf("protocol display label was not localized: %q", item.label)
		}
	}
}
