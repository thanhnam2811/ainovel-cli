package tui

import "github.com/voocel/ainovel-cli/internal/localization"

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
