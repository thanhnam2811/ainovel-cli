package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type helpState struct {
	viewport viewport.Model
}

func newHelpState(width, height int) *helpState {
	boxW, boxH := reportModalSize(width, height)
	contentW := paddedModalContentWidth(boxW)
	text := renderHelpText(contentW)

	vp := viewport.New(contentW, boxH-4)
	vp.SetContent(text)
	return &helpState{viewport: vp}
}

func renderHelpText(width int) string {
	titleStyle := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(colorAccent2).Bold(true)
	usageStyle := lipgloss.NewStyle().Foreground(colorMuted)
	descStyle := lipgloss.NewStyle().Foreground(bodyTextColor)
	hintStyle := lipgloss.NewStyle().Foreground(colorDim)

	var b strings.Builder
	b.WriteString(titleStyle.Render(uiText("命令帮助", "Trợ giúp lệnh")))
	b.WriteString("\n\n")

	for i, spec := range commandSpecs() {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(nameStyle.Render("/" + spec.Name))
		if len(spec.Aliases) > 0 {
			b.WriteString(usageStyle.Render("  alias: /" + strings.Join(spec.Aliases, " /")))
		}
		b.WriteString("\n")
		b.WriteString(usageStyle.Render("Usage: " + spec.Usage))
		b.WriteString("\n")
		b.WriteString(descStyle.Render(wrapText(spec.Description, width)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(titleStyle.Render(uiText("快捷键", "Phím tắt")))
	b.WriteString("\n\n")
	for _, line := range []string{
		uiText("输入 / 搜索命令", "Nhập / để tìm lệnh"),
		uiText("↑↓ 选择命令候选", "↑↓ chọn lệnh gợi ý"),
		uiText("Tab/Enter 接受补全", "Tab/Enter chấp nhận gợi ý"),
		uiText("Esc 关闭当前命令面板", "Esc đóng bảng lệnh hiện tại"),
		uiText("Ctrl+R 切换选中复制模式（关闭鼠标上报后可拖拽选中复制，再按一次恢复）", "Ctrl+R bật/tắt chế độ chọn để sao chép (tắt mouse reporting để kéo chọn; nhấn lại để khôi phục)"),
	} {
		b.WriteString(hintStyle.Render(line))
		b.WriteString("\n")
	}
	return b.String()
}

func renderHelpModal(width, height int, state *helpState) string {
	if state == nil {
		return ""
	}

	boxW, boxH := reportModalSize(width, height)
	contentW := paddedModalContentWidth(boxW)

	if state.viewport.Width != contentW {
		state.viewport.Width = contentW
	}
	if state.viewport.Height != boxH-4 {
		state.viewport.Height = boxH - 4
	}

	modal := renderPaddedModalFrame(
		boxW,
		boxH,
		uiText("命令帮助", "Trợ giúp lệnh"),
		uiText("  ↑↓ 滚动 · Esc 关闭", "  ↑↓ cuộn · Esc đóng"),
		strings.Split(state.viewport.View(), "\n"),
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.help == nil {
		return m, nil
	}
	switch msg.Type {
	case tea.KeyEsc:
		m.help = nil
		return m, m.textarea.Focus()
	case tea.KeyUp:
		m.help.viewport.ScrollUp(1)
		return m, nil
	case tea.KeyDown:
		m.help.viewport.ScrollDown(1)
		return m, nil
	case tea.KeyPgUp:
		m.help.viewport.HalfPageUp()
		return m, nil
	case tea.KeyPgDown:
		m.help.viewport.HalfPageDown()
		return m, nil
	default:
		return m, nil
	}
}
