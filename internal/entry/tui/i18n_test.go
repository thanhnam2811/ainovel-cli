package tui

import "testing"

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
