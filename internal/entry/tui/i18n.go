package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/voocel/ainovel-cli/internal/localization"
)

func uiText(upstreamChinese, vietnamese string) string {
	if localization.IsVietnamese() {
		return vietnamese
	}
	return upstreamChinese
}

func localizedCommandDescription(name, upstream string) string {
	if !localization.IsVietnamese() {
		return upstream
	}
	if translated, ok := vietnameseCommandDescriptions[name]; ok {
		return translated
	}
	return upstream
}

func localizedStatusLabel(status, upstream string) string {
	if !localization.IsVietnamese() {
		return upstream
	}
	if translated, ok := vietnameseStatusLabels[status]; ok {
		return translated
	}
	return upstream
}

func localizedFieldLabel(upstream string) string {
	if !localization.IsVietnamese() {
		return upstream
	}
	if translated, ok := vietnameseFieldLabels[upstream]; ok {
		return translated
	}
	return upstream
}

type localizedTeaModel struct {
	inner tea.Model
}

func wrapLocalizedTUI(inner tea.Model) tea.Model {
	if !localization.IsVietnamese() {
		return inner
	}
	return localizedTeaModel{inner: inner}
}

func (m localizedTeaModel) Init() tea.Cmd { return m.inner.Init() }

func (m localizedTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	next, cmd := m.inner.Update(msg)
	if next != nil {
		m.inner = next
	}
	return m, cmd
}

func (m localizedTeaModel) View() string {
	return localizeRenderedTUI(m.inner.View())
}

func localizeRenderedTUI(view string) string {
	if !localization.IsVietnamese() || view == "" {
		return view
	}

	plain := []string{
		"加载中...", "Đang tải...",
		"终端宽度不足，请至少扩展到 100 列", "Terminal quá hẹp, hãy mở rộng tối thiểu 100 cột",
		"输入剧情干预，例如：把感情线提前到第4章", "Nhập chỉ dẫn cốt truyện, ví dụ: đẩy tuyến tình cảm lên chương 4",
		"正在初始化创作...", "Đang khởi tạo quá trình sáng tác...",
		"正在暂停创作...", "Đang tạm dừng quá trình sáng tác...",
		"逐章验收等待中：输入修改意见，或 /next 放行下一章", "Đang chờ duyệt chương: nhập góp ý chỉnh sửa hoặc /next để cho phép chương tiếp theo",
		"创作已暂停，输入任意内容继续创作", "Sáng tác đã tạm dừng, nhập nội dung bất kỳ để tiếp tục",
		"运行中断，输入任意内容恢复创作", "Phiên chạy bị gián đoạn, nhập nội dung bất kỳ để tiếp tục",
		"快速开始", "Bắt đầu nhanh",
		"共创规划", "Cùng lập kế hoạch",
		"先与 AI 对话澄清，再开始创作", "Trao đổi với AI để làm rõ trước khi bắt đầu sáng tác",
		"一句话直接开始写", "Bắt đầu viết trực tiếp từ một câu yêu cầu",
		"先输入你的核心想法，Enter 开始与 AI 共创", "Nhập ý tưởng cốt lõi, Enter để bắt đầu cùng AI",
		"输入一句小说需求，Enter 直接开始创作", "Nhập một câu yêu cầu về truyện, Enter để bắt đầu ngay",
		"AI 正在整理你的要求...", "AI đang tổng hợp yêu cầu của bạn...",
		"继续补充，或按 Ctrl+S 应用方向并继续创作", "Tiếp tục bổ sung hoặc Ctrl+S để áp dụng hướng đi và viết tiếp",
		"继续补充，或按 Ctrl+S 开始创作", "Tiếp tục bổ sung hoặc Ctrl+S để bắt đầu sáng tác",
		"继续补充你的要求，Enter 发送给 AI", "Tiếp tục bổ sung yêu cầu, Enter để gửi cho AI",
		"启动模式", "Chế độ khởi động",
		"AI 建议（按 1/2/3 组合，可编辑后发送）：", "Gợi ý từ AI (nhấn 1/2/3 để kết hợp; có thể sửa trước khi gửi):",
		"先把需求聊清楚，再开始创作", "Làm rõ yêu cầu trước khi bắt đầu sáng tác",
		"阶段共创", "Cùng lập kế hoạch giai đoạn",
		"规划后续走向，再继续创作", "Lập kế hoạch hướng tiếp theo rồi tiếp tục sáng tác",
		"已暂停创作，进入阶段共创 —— AI 会结合当前故事进度，和你一起规划接下来的走向。", "Đã tạm dừng sáng tác để cùng lập kế hoạch giai đoạn — AI sẽ dựa trên tiến độ hiện tại để cùng bạn định hướng phần tiếp theo.",
		"AI 回复中 · ↑↓ 滚对话 · 滚轮滚指令 · Esc 退出", "AI đang trả lời · ↑↓ cuộn hội thoại · con lăn cuộn chỉ dẫn · Esc thoát",
		"Enter 发送 · ↑↓ 滚对话 · 滚轮滚指令 · Esc 退出", "Enter gửi · ↑↓ cuộn hội thoại · con lăn cuộn chỉ dẫn · Esc thoát",
		"Enter 继续补充 · Ctrl+S 开始创作 · ↑↓ 滚对话 · 滚轮滚指令 · Esc 退出", "Enter bổ sung · Ctrl+S bắt đầu sáng tác · ↑↓ cuộn hội thoại · con lăn cuộn chỉ dẫn · Esc thoát",
		"Enter 继续补充 · Ctrl+S 应用并继续 · ↑↓ 滚对话 · 滚轮滚指令 · Esc 退出", "Enter bổ sung · Ctrl+S áp dụng và tiếp tục · ↑↓ cuộn hội thoại · con lăn cuộn chỉ dẫn · Esc thoát",
		"等待 AI 回复 · Esc 退出共创", "Đang chờ AI trả lời · Esc thoát chế độ cùng lập kế hoạch",
		"Tab 切换启动模式 · 输入 / 搜索命令 · Enter 直接开始创作 · Esc 清空输入", "Tab đổi chế độ khởi động · nhập / để tìm lệnh · Enter bắt đầu ngay · Esc xóa input",
		"Tab 切换启动模式 · 输入 / 搜索命令 · Enter 开始共创对话 · Esc 清空输入", "Tab đổi chế độ khởi động · nhập / để tìm lệnh · Enter bắt đầu hội thoại lập kế hoạch · Esc xóa input",
		"正在暂停创作 · 请等待当前轮次结束", "Đang tạm dừng sáng tác · vui lòng chờ lượt hiện tại kết thúc",
		"输入 / 搜索命令 · Enter 继续创作 · Esc 清空输入", "Nhập / để tìm lệnh · Enter tiếp tục sáng tác · Esc xóa input",
		"输入 / 搜索命令 · 点击/Tab 切换面板 · ↑↓ 滚动 · End 跳底 · Ctrl+L 清屏 · Esc 暂停 · Enter 发送", "Nhập / để tìm lệnh · click/Tab đổi panel · ↑↓ cuộn · End xuống cuối · Ctrl+L xóa màn hình · Esc tạm dừng · Enter gửi",
		"Press Ctrl+C again to exit", "Nhấn Ctrl+C lần nữa để thoát",
		"+ 新增 Provider…", "+ Thêm Provider…",
		"Provider 名称", "Tên Provider",
		"输入 API Key", "Nhập API Key",
		"输入新 Key，留空保留", "Nhập Key mới; để trống để giữ nguyên",
		"留空使用默认地址", "Để trống để dùng địa chỉ mặc định",
		"该 Provider 必须配置 API Key", "Provider này bắt buộc phải có API Key",
		"模型 ID", "Model ID",
		"请先添加模型", "Hãy thêm model trước",
		"测试连接", "Kiểm tra kết nối",
		"保存配置", "Lưu cấu hình",
		"默认地址", "Địa chỉ mặc định",
		"已清除", "Đã xóa",
		"未设置", "Chưa đặt",
		"当前没有可用 provider", "Hiện không có provider khả dụng",
		"/model 切换模型", "/model Chuyển model",
		"角色", "Vai trò",
		"推理强度", "Mức suy luận",
		"Tab 切字段   ←→ 切选项   Enter 应用   Esc 取消", "Tab đổi trường   ←→ đổi lựa chọn   Enter áp dụng   Esc hủy",
		"默认(继承)", "Mặc định (kế thừa)",
		"诊断报告不可用", "Không có báo cáo chẩn đoán",
		"已导出脱敏诊断（可贴到 GitHub issue）", "Đã xuất chẩn đoán đã ẩn dữ liệu nhạy cảm (có thể đính kèm GitHub issue)",
		"脱敏诊断导出失败：", "Xuất chẩn đoán đã ẩn dữ liệu nhạy cảm thất bại: ",
		"未发现问题", "Không phát hiện vấn đề",
		"可执行动作", "Hành động khả dụng",
		"正在生成诊断报告", "Đang tạo báo cáo chẩn đoán",
		"开始时间 ", "Bắt đầu ",
		"正在读取当前小说 output 产物并分析流程、质量、规划和上下文问题。项目较大时可能需要几秒。", "Đang đọc output hiện tại và phân tích luồng, chất lượng, kế hoạch và ngữ cảnh. Dự án lớn có thể mất vài giây.",
		"Esc 可先关闭面板，后台分析完成后下次打开会重新生成。", "Có thể nhấn Esc để đóng panel; lần mở sau báo cáo sẽ được tạo lại.",
		"导入外部小说", "Nhập truyện bên ngoài",
		"源文件 ", "Tệp nguồn ",
		"流程日志", "Nhật ký tiến trình",
		"导入失败", "Nhập thất bại",
		"Esc 关闭面板", "Esc đóng panel",
		"切分完成，等待你核对", "Đã chia đoạn, chờ bạn xác nhận",
		"导入已暂停，等待你的操作", "Đã tạm dừng nhập, chờ thao tác của bạn",
		"导入完成，Foundation 与章节已就绪", "Nhập hoàn tất, Foundation và các chương đã sẵn sàng",
		"Esc 取消导入", "Esc hủy nhập",
		"来源 ", "Nguồn ",
		"Esc 取消", "Esc hủy",
		"仿写画像处理失败", "Xử lý hồ sơ mô phỏng phong cách thất bại",
		"仿写画像已就绪，后续 Agent 会从 novel_context 读取", "Hồ sơ mô phỏng phong cách đã sẵn sàng; Agent sẽ đọc từ novel_context",
		"↑↓ 滚动 · Esc 取消/关闭", "↑↓ cuộn · Esc hủy/đóng",
		"仿写画像", "Hồ sơ mô phỏng phong cách",
		"生成仿写画像", "Tạo hồ sơ mô phỏng phong cách",
		"导入仿写画像", "Nhập hồ sơ mô phỏng phong cách",
		"用法：/simulate", "Cách dùng: /simulate",
		"用法：/importsim <profile.json>", "Cách dùng: /importsim <profile.json>",
	}
	view = strings.NewReplacer(plain...).Replace(view)

	for upstream, vietnamese := range vietnameseSectionTitles {
		view = strings.ReplaceAll(view, panelTitleStyle.Render(upstream), panelTitleStyle.Render(vietnamese))
	}
	for upstream, vietnamese := range vietnameseFieldValues {
		view = strings.ReplaceAll(view, fieldValueStyle.Render(upstream), fieldValueStyle.Render(vietnamese))
		view = strings.ReplaceAll(view, highlightValueStyle.Render(upstream), highlightValueStyle.Render(vietnamese))
	}
	return view
}

var vietnameseSectionTitles = map[string]string{
	"概览":   "Tổng quan",
	"运行角色": "Vai trò",
	"返工":   "Làm lại",
	"干预":   "Can thiệp",
	"验收停靠": "Chờ duyệt",
	"用量":   "Sử dụng",
	"缓存":   "Cache",
}

var vietnameseFieldValues = map[string]string{
	"前提":   "Tiền đề",
	"大纲":   "Dàn ý",
	"写作":   "Sáng tác",
	"完成":   "Hoàn tất",
	"初始化":  "Khởi tạo",
	"运行中":  "Đang chạy",
	"暂停中":  "Đang tạm dừng",
	"已暂停":  "Đã tạm dừng",
	"已完成":  "Đã hoàn tất",
	"空闲":   "Rảnh",
	"评审":   "Đánh giá",
	"重写":   "Viết lại",
	"打磨":   "Biên tập",
	"干预":   "Can thiệp",
	"逐章验收": "Duyệt từng chương",
	"自动":   "Tự động",
}

var vietnameseFieldLabels = map[string]string{
	"运行态": "Trạng thái",
	"阶段":  "Giai đoạn",
	"流程":  "Luồng",
	"推进":  "Tiến hành",
	"已完成": "Đã xong",
	"已规划": "Đã lên kế hoạch",
	"进度":  "Tiến độ",
	"字数":  "Số từ",
	"当前":  "Hiện tại",
	"待恢复": "Chờ khôi phục",
	"队列":  "Hàng đợi",
	"原因":  "Lý do",
	"待处理": "Chờ xử lý",
	"等待":  "Đang chờ",
	"输入":  "Input",
	"输出":  "Output",
	"费用":  "Chi phí",
	"节省":  "Tiết kiệm",
	"预算":  "Ngân sách",
}

var vietnameseStatusLabels = map[string]string{
	"READY":      "sẵn sàng",
	"RUNNING":    "đang chạy",
	"PAUSED":     "đã tạm dừng",
	"PAUSING":    "đang tạm dừng",
	"COMPLETED":  "hoàn tất",
	"ERROR":      "lỗi",
	"REVIEW":     "đang duyệt",
	"REWRITING":  "đang viết lại",
	"POLISHING":  "đang biên tập",
	"STEERING":   "đang can thiệp",
	"IMPORTING":  "đang nhập",
	"SIMULATING": "đang mô phỏng",
}

var vietnameseCommandDescriptions = map[string]string{
	"help":      "Xem danh sách lệnh",
	"model":     "Chuyển model và mức suy luận cho từng vai trò",
	"config":    "Thêm hoặc chỉnh sửa Provider, model và cửa sổ ngữ cảnh",
	"diag":      "Chẩn đoán tình trạng quá trình sáng tác",
	"review":    "Bật/tắt chế độ duyệt từng chương",
	"next":      "Cho phép viết chương mới sau khi duyệt",
	"start":     "Tạo truyện mới từ tệp thiết lập hoặc dàn ý",
	"import":    "Nhập truyện ngoài theo ngữ nghĩa; không có tham số sẽ tiếp tục lần nhập chưa hoàn tất",
	"reopen":    "Mở lại truyện đã hoàn thành để viết tiếp",
	"cocreate":  "Tạm dừng sáng tác để cùng lập kế hoạch cho giai đoạn tiếp theo",
	"simulate":  "Đọc ./simulate để tạo hoặc cập nhật hồ sơ mô phỏng phong cách",
	"importsim": "Nhập hồ sơ mô phỏng phong cách có sẵn và hợp nhất theo dấu vân tay ngữ liệu",
	"sync":      "Kiểm tra hoặc chấp nhận chỉnh sửa thủ công ở các chương đã hoàn thành",
	"export":    "Xuất các chương đã hoàn thành",
	"quit":      "Thoát chương trình",
}
