package tui

import (
	"strings"
	"testing"
)

func TestSharedTUIDefaultsToVietnamese(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "")
	if got := uiText("命令帮助", "Trợ giúp lệnh"); got != "Trợ giúp lệnh" {
		t.Fatalf("uiText default locale = %q", got)
	}
	if got := localizedFieldLabel("运行态"); got != "Trạng thái" {
		t.Fatalf("field label = %q", got)
	}
	if got := localizedCommandDescription("help", "查看命令列表"); got != "Xem danh sách lệnh" {
		t.Fatalf("command description = %q", got)
	}
}

func TestSharedTUIChineseControlPreservesUpstream(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "zh")
	if got := uiText("命令帮助", "Trợ giúp lệnh"); got != "命令帮助" {
		t.Fatalf("uiText zh control = %q", got)
	}
	if got := localizedFieldLabel("运行态"); got != "运行态" {
		t.Fatalf("field label zh control = %q", got)
	}
	if got := localizedCommandDescription("help", "查看命令列表"); got != "查看命令列表" {
		t.Fatalf("command description zh control = %q", got)
	}
}

func TestRenderedTUILocalizesChromeButPreservesArbitraryStoryText(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "vi")
	story := "他决定快速开始调查，但这只是小说正文。"
	view := "加载中...\n" + story + "\n" + panelTitleStyle.Render("概览")
	got := localizeRenderedTUI(view)
	if !strings.Contains(got, "Đang tải...") {
		t.Fatalf("fixed chrome was not localized: %q", got)
	}
	if !strings.Contains(got, story) {
		t.Fatalf("arbitrary story text changed: %q", got)
	}
	if !strings.Contains(got, panelTitleStyle.Render("Tổng quan")) {
		t.Fatalf("styled section title was not localized: %q", got)
	}
}

func TestRenderedTUIChineseControlIsExact(t *testing.T) {
	t.Setenv("AINOVEL_LOCALE", "zh")
	view := "加载中...\n" + panelTitleStyle.Render("概览")
	if got := localizeRenderedTUI(view); got != view {
		t.Fatalf("zh control changed rendered TUI: got %q want %q", got, view)
	}
}
