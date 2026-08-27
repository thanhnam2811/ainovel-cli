package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/voocel/ainovel-cli/internal/host"
)

func renderTopBar(snap host.UISnapshot, width int, spinnerFrame, version string) string {
	bookTitle := snap.BookTitle
	if bookTitle == "" {
		bookTitle = uiText("未定书名", "Chưa đặt tên")
	}

	var infoParts []string
	if version != "" {
		infoParts = append(infoParts, "ainovel-cli "+version)
	}
	if snap.Provider != "" {
		infoParts = append(infoParts, snap.Provider)
	}
	if snap.ModelName != "" {
		if w := formatContextWindow(snap.ModelContextWindow); w != "" {
			infoParts = append(infoParts, snap.ModelName+"("+w+")")
		} else {
			infoParts = append(infoParts, snap.ModelName)
		}
	}
	if snap.Style != "" && snap.Style != "default" {
		infoParts = append(infoParts, snap.Style)
	}
	leftText := strings.Join(infoParts, " · ")

	label := snap.StatusLabel
	if label == "" {
		label = "READY"
	}
	color, ok := statusColors[label]
	if !ok {
		color = colorDim
	}
	disp, ok := statusDisplay[label]
	if !ok {
		disp = struct {
			icon  string
			label string
		}{"○", strings.ToLower(label)}
	}
	icon := disp.icon
	if snap.IsRunning && spinnerFrame != "" {
		icon = spinnerFrame
	}
	var status string
	if icon != "" {
		status = statusIconStyle.Foreground(color).Render(icon) + " " + statusLabelStyle.Render(localizedStatusLabel(label, disp.label))
	} else {
		status = statusLabelStyle.Render(localizedStatusLabel(label, disp.label))
	}

	innerW := max(12, width-2)
	titleText := truncate(bookTitle, max(8, innerW/3))
	centerW := max(16, lipgloss.Width(titleText)+6)
	if centerW > innerW-24 {
		centerW = max(8, innerW-24)
	}
	sideTotal := innerW - centerW
	if sideTotal < 0 {
		sideTotal = 0
		centerW = innerW
	}
	leftW := sideTotal / 2
	rightW := innerW - centerW - leftW

	leftCell := lipgloss.NewStyle().Width(leftW).AlignHorizontal(lipgloss.Left).Foreground(colorDim).Render(truncate(leftText, leftW))
	centerCell := lipgloss.NewStyle().Width(centerW).AlignHorizontal(lipgloss.Center).Bold(true).Foreground(bodyTextColor).Render(titleText)
	rightCell := lipgloss.NewStyle().Width(rightW).AlignHorizontal(lipgloss.Right).Render(status)

	content := leftCell + centerCell + rightCell
	return topBarStyle.Width(width).Border(baseBorder, false, false, true, false).BorderForeground(colorDim).Render(content)
}

func renderStatePanel(vp viewport.Model, width, height int, focused bool) string {
	borderColor := colorDim
	if focused {
		borderColor = colorAccent
	}
	style := lipgloss.NewStyle().Width(width).Height(height).MaxHeight(height).
		Border(baseBorder, false, true, false, false).BorderForeground(borderColor).Padding(1, 1, 0, 1)
	return style.Render(vp.View())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderDetailPanel(vp viewport.Model, width, height int, focused bool) string {
	borderColor := colorDim
	if focused {
		borderColor = colorAccent
	}
	style := lipgloss.NewStyle().Width(width).Height(height).MaxHeight(height).
		Border(baseBorder, false, false, false, true).BorderForeground(borderColor).Padding(0, 1)
	return style.Render(vp.View())
}

func renderWelcome(width, height int, errMsg string, mode startupMode, importHint string) string {
	title := lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render("A I N O V E L")
	subtitle := lipgloss.NewStyle().Foreground(colorMuted).Italic(true).Render("AI-Powered Novel Creation Engine")

	divW := 44
	if divW > width-8 {
		divW = width - 8
	}
	divider := lipgloss.NewStyle().Foreground(colorDim).Render(strings.Repeat("~", divW))

	features := []struct{ icon, label, desc string }{
		{">>", uiText("多模型协作", "Phối hợp đa model"), uiText("Architect 规划 / Writer 创作 / Editor 审阅", "Architect lập kế hoạch / Writer sáng tác / Editor biên tập")},
		{"::", uiText("断点恢复", "Khôi phục tiến trình"), uiText("崩溃或中断后从上次进度自动续写", "Tự tiếp tục từ tiến độ gần nhất sau khi lỗi hoặc gián đoạn")},
		{"<>", uiText("实时干预", "Can thiệp thời gian thực"), uiText("创作过程中随时调整剧情走向", "Điều chỉnh hướng cốt truyện bất cứ lúc nào khi đang sáng tác")},
		{"##", uiText("分层长篇", "Truyện dài phân tầng"), uiText("支持卷-弧-章分层结构的长篇创作", "Hỗ trợ cấu trúc Tập - Arc - Chương cho truyện dài")},
	}
	iconStyle := lipgloss.NewStyle().Foreground(colorAccent2).Bold(true)
	featLabelStyle := lipgloss.NewStyle().Foreground(bodyTextColor)
	descStyle := lipgloss.NewStyle().Foreground(colorDim)
	var featLines []string
	for _, f := range features {
		line := iconStyle.Render(f.icon) + " " + featLabelStyle.Render(f.label) + "  " + descStyle.Render(f.desc)
		featLines = append(featLines, line)
	}
	feats := strings.Join(featLines, "\n")

	prompt := lipgloss.NewStyle().Foreground(bodyTextColor).Render(uiText("在下方输入你的小说需求开始创作", "Nhập yêu cầu về truyện ở bên dưới để bắt đầu sáng tác"))
	modeLine := lipgloss.NewStyle().Foreground(colorMuted).Render(uiText("当前模式：", "Chế độ hiện tại: ") + mode.label() + " · " + mode.subtitle())

	examples := []string{
		uiText("写一部 12 章都市悬疑小说，主角是一名女法医", "Viết tiểu thuyết trinh thám đô thị 12 chương, nhân vật chính là nữ pháp y"),
		uiText("创作一部仙侠长篇，主角从凡人修炼至飞升", "Viết truyện tiên hiệp dài, nhân vật chính tu luyện từ phàm nhân đến phi thăng"),
		uiText("写一个科幻短篇，讲述 AI 觉醒后的伦理困境", "Viết truyện khoa học viễn tưởng ngắn về những vấn đề đạo đức sau khi AI thức tỉnh"),
	}
	exStyle := lipgloss.NewStyle().Foreground(colorAccent)
	dotStyle := lipgloss.NewStyle().Foreground(colorDim)
	var exLines []string
	for _, ex := range examples {
		exLines = append(exLines, dotStyle.Render("  . ")+exStyle.Render(ex))
	}
	exBlock := strings.Join(exLines, "\n")

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(subtitle)
	b.WriteString("\n\n")
	b.WriteString(divider)
	b.WriteString("\n\n")
	b.WriteString(feats)
	b.WriteString("\n\n")
	b.WriteString(divider)
	b.WriteString("\n\n")
	b.WriteString(modeLine)
	b.WriteString("\n\n")
	b.WriteString(prompt)
	b.WriteString("\n\n")
	b.WriteString(exBlock)
	b.WriteString("\n\n")
	if importHint != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(colorAccent2).Bold(true).Render("! " + importHint))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(colorDim).Render(uiText(
			"已有设定/大纲？/start <文件路径> 创建新书 · 已有小说存稿？/import <文件路径> 导入续写",
			"Có thiết lập/dàn ý? dùng /start <đường-dẫn-tệp> · Có bản thảo? dùng /import <đường-dẫn-tệp> để nhập và viết tiếp",
		)))
	}
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(colorDim).Italic(true).Render(uiText(
		"Tab 切换模式 · 快速开始下 Enter 直接创作 · 共创规划下 Enter 进入对话",
		"Tab đổi chế độ · Enter để bắt đầu nhanh · Ở chế độ cùng lập kế hoạch, Enter để vào hội thoại",
	)))

	if errMsg != "" {
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(colorError).Bold(true).Render("! " + errMsg))
	}

	return lipgloss.NewStyle().Width(width).Height(height).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).Render(b.String())
}
