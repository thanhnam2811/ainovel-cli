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

Nếu upstream thêm hoặc thay đổi một prompt mà bản Việt chưa có file tương ứng, chương trình **không lỗi và không dùng bản dịch cũ**. Nó giữ prompt upstream mới nhất rồi thêm chỉ thị yêu cầu output tiếng Việt. Nhờ vậy bugfix/protocol update từ upstream có hiệu lực ngay cả khi bản dịch chưa kịp cập nhật.

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

## Phạm vi Việt hóa

Ưu tiên của lớp locale là những thứ ảnh hưởng trực tiếp đến chất lượng sáng tác:

- Writer system prompt
- writing voice
- anti-AI-tone reference
- default writing style
- các core/function prompt khác có cơ chế fallback an toàn về upstream

Source Go, domain model, tool/schema identifier và comment nội bộ không được dịch chỉ để “trông Việt hơn”. Đây là chủ ý để giảm maintenance debt.

## Nguyên tắc sync upstream

Không đổi module path `github.com/voocel/ainovel-cli`. Khi upstream có thay đổi, merge/rebase upstream như một fork bình thường; phần khác biệt của bản Việt chủ yếu nằm trong `assets/locales/` và một hook nhỏ ở startup.

Khi upstream sửa protocol trong một prompt chưa có bản Việt hóa, fallback sẽ dùng ngay prompt upstream mới. Với prompt đã dịch đầy đủ, cần review diff upstream và cập nhật bản dịch tương ứng trước khi phát hành bản Việt tiếp theo.

## Phát triển bản dịch

Thêm file đúng đường dẫn upstream dưới `assets/locales/vi/`. Ví dụ:

```text
assets/prompts/editor.md
→ assets/locales/vi/prompts/editor.md
```

Loader sẽ tự ưu tiên file Việt nếu tồn tại. Không thêm `if locale == ...` rải rác trong engine.

Trước khi merge, chạy:

```bash
go test ./...
```

Các test locale kiểm tra ít nhất:

- bản Việt được nạp thật;
- `{{VOICE}}` không bị mất;
- prompt chưa dịch vẫn nhận chỉ thị output tiếng Việt;
- `AINOVEL_LOCALE=zh` giữ asset upstream nguyên vẹn;
- apply locale nhiều lần không làm prompt phình lặp.