# ainovel-cli — bản Việt hóa upstream-friendly

Bản fork này bám sát `voocel/ainovel-cli` và thêm lớp Việt hóa theo kiểu **overlay**, thay vì dịch hàng loạt source Go.

Mục tiêu là giữ toàn bộ engine, state machine, checkpoint/resume, import pipeline và bugfix mới nhất từ upstream, đồng thời cho các agent sáng tác bằng tiếng Việt mà không biến mỗi lần sync upstream thành một đợt resolve conflict lớn.

## Cách hoạt động

Runtime vẫn nạp asset gốc từ upstream trước. Sau đó `assets.ApplyLocale` phủ các asset có bản Việt hóa trong:

```text
assets/locales/vi/
├── prompts/
├── references/
├── styles/
└── voice.md
```

Nếu một prompt chưa có bản Việt, chương trình **không lỗi và không giữ một bản dịch cũ đã lệch protocol**. Nó dùng prompt upstream hiện tại rồi thêm chỉ thị bắt buộc output tiếng Việt. Nhờ vậy bugfix/protocol update từ upstream vẫn có hiệu lực ngay khi sync.

Các tên tool, JSON field, enum, protocol marker và placeholder như `novel_context`, `commit_chapter`, `required_beats`, `{{VOICE}}` luôn được giữ nguyên.

## Locale mặc định

Fork này mặc định dùng tiếng Việt:

```bash
ainovel-cli
```

Có thể chạy nguyên asset upstream tiếng Trung để đối chiếu hoặc debug:

```bash
AINOVEL_LOCALE=zh ainovel-cli
```

`vi-VN` và `vi` đều được hỗ trợ nếu đặt `AINOVEL_LOCALE` thủ công.

## Phạm vi Việt hóa hiện tại

Các asset đã có bản Việt đầy đủ:

- Writer system prompt
- writing voice
- anti-AI-tone reference
- default writing style

Các core/function prompt khác hiện dùng fallback upstream + chỉ thị output tiếng Việt cho tới khi có bản dịch được review theo đúng protocol mới nhất.

Source Go, domain model, tool/schema identifier và comment nội bộ không được dịch chỉ để “trông Việt hơn”. Đây là chủ ý để giảm maintenance debt.

## Override của người dùng vẫn được ưu tiên

Lớp locale không được phép phá cơ chế tùy biến có sẵn của upstream. Với `voice.md`, `anti-ai-tone.md` và `styles/*.md`, thứ tự vẫn là:

```text
Vietnamese builtin < ~/.ainovel/style < <book>/style
```

Nghĩa là bản Việt chỉ thay builtin mặc định; cấu hình global và cấu hình riêng của từng sách vẫn thắng như upstream.

## Rules tiếng Việt

`~/.ainovel/rules/*.md` và `./.ainovel/rules/*.md` là dữ liệu do người dùng viết bằng ngôn ngữ tự nhiên, nên có thể viết tiếng Việt trực tiếp và không cần một “bản dịch rules” riêng.

Upstream còn có mechanical baseline nội bộ cho các câu sáo/từ mòn tiếng Trung với threshold được hiệu chỉnh từ chạy dài thực tế. Fork này **không dịch máy các threshold đó sang tiếng Việt** vì như vậy sẽ tạo ra heuristic giả chính xác. Anti-AI semantic reference tiếng Việt đã được bật ngay; mechanical baseline tiếng Việt chỉ nên bổ sung sau khi có corpus chạy thực tế để đo tần suất và false-positive.

## Cài đặt và self-update

Installer và lệnh `update` của fork đều trỏ về release của `thanhnam2811/ainovel-cli`, không tải binary từ `voocel/ainovel-cli` rồi làm mất lớp Việt hóa.

Sau khi fork có release riêng, có thể cài bằng:

```bash
curl -fsSL https://raw.githubusercontent.com/thanhnam2811/ainovel-cli/main/scripts/install.sh | sh
```

Trước release đầu tiên, build trực tiếp từ source là cách chắc chắn nhất.

## Nguyên tắc sync upstream

Không đổi module path `github.com/voocel/ainovel-cli`. Khi upstream có thay đổi, merge/rebase upstream như một fork bình thường; phần khác biệt của bản Việt chủ yếu nằm trong `assets/locales/` và một hook nhỏ ở startup/release boundary.

Khi upstream sửa protocol trong một prompt chưa có bản Việt hóa, fallback sẽ dùng ngay prompt upstream mới. Với prompt đã dịch đầy đủ, cần review diff upstream và cập nhật bản dịch tương ứng trước khi phát hành bản Việt tiếp theo.

## Phát triển bản dịch

Prompt core/function và reference có mapping tập trung trong `assets/locales.go`; style preset được overlay theo đúng tên file hiện có. Khi thêm bản dịch, đặt file dưới `assets/locales/vi/` theo đường dẫn tương ứng và thêm mapping tập trung nếu loại asset đó dùng struct field cố định. Không thêm `if locale == ...` rải rác trong engine.

Ví dụ:

```text
assets/prompts/editor.md
→ assets/locales/vi/prompts/editor.md
```

Trước khi merge, chạy:

```bash
go test ./...
```

Các test locale kiểm tra ít nhất:

- bản Việt được nạp thật;
- `{{VOICE}}` không bị mất;
- prompt chưa dịch vẫn nhận chỉ thị output tiếng Việt;
- `AINOVEL_LOCALE=zh` giữ asset upstream nguyên vẹn;
- apply locale nhiều lần không làm prompt phình lặp;
- global/book style override vẫn thắng localized builtin.