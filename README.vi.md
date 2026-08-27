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

Toàn bộ prompt đang được runtime load đã có bản Việt đầy đủ:

- Architect Short
- Architect Long
- Writer
- Editor
- import-segment / import-analyze / import-synthesize / import-range
- revision-analyze
- simulation-source / simulation-merge
- arbiter-plan-start / arbiter-intervention / arbiter-failure

Các asset văn phong quan trọng đã có bản Việt:

- writing voice
- anti-AI-tone reference
- default writing style

Fallback upstream + chỉ thị output tiếng Việt vẫn được giữ trong loader như một **cơ chế an toàn cho tương lai**: nếu upstream thêm prompt mới hoặc một overlay bị chủ động gỡ vì chưa port kịp protocol, fork vẫn chạy trên contract upstream mới nhất.

Source Go, domain model, tool/schema identifier và comment nội bộ không được dịch chỉ để “trông Việt hơn”. Đây là chủ ý để giảm maintenance debt.

## Guard cho protocol

Bản dịch prompt có smoke test kiểm tra các token máy đọc quan trọng vẫn tồn tại. Ví dụ:

- Architect Short: `save_book`, `save_foundation`, `revise_outline`, `audit_foundation`, `foundation_ready`, `remaining`.
- Architect Long: `layered_outline`, `update_compass`, `append_volume`, `complete_book`, `expand_arc`, `final_volume`, `open_threads`, `completion_signals`.
- Editor: `save_review`, `requires_change`, `rule_violations`, `accept`, `polish`, `rewrite`.
- Import: enum/value domain như `hook_type`, `dominant_strand`, `planning_tier`, `story_status`, boundary kind và foreshadow action.
- Arbiter: thứ tự action `answer → rules → hold → reopen → dispatch`, các trạng thái `rewrites_drained`, `phase = complete`, `dynamic_planning`, cùng tập agent hợp lệ.

Test không khóa câu chữ dịch; mục tiêu là ngăn một lần chỉnh wording vô tình làm rơi protocol hoặc đổi routing behavior.

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

Khi upstream sửa protocol trong một prompt đã dịch đầy đủ, phải review diff upstream và port đầy đủ thay đổi behavior trước khi phát hành bản Việt tiếp theo. Nếu chưa port an toàn được, **ưu tiên gỡ/không phát hành overlay đó để runtime fallback về prompt upstream mới + chỉ thị tiếng Việt**, thay vì giữ một bản dịch stale.

## Phát triển bản dịch

Prompt core/function và reference có mapping tập trung trong `assets/locales.go`; style preset được overlay theo đúng tên file hiện có. Khi thêm bản dịch, đặt file dưới `assets/locales/vi/` theo đường dẫn tương ứng và thêm mapping tập trung nếu loại asset đó dùng struct field cố định. Không thêm `if locale == ...` rải rác trong engine.

Ví dụ:

```text
assets/prompts/editor.md
→ assets/locales/vi/prompts/editor.md
```

Checklist chi tiết nằm trong `assets/locales/vi/prompts/README.md` và `CONTRACT.md`.

Trước khi merge, chạy:

```bash
go test ./...
```

Các test locale kiểm tra ít nhất:

- bản Việt được nạp thật;
- `{{VOICE}}` không bị mất;
- core-agent không làm rơi token protocol quan trọng;
- function prompt không làm rơi enum/routing contract quan trọng;
- fallback prompt vẫn nhận chỉ thị output tiếng Việt;
- `AINOVEL_LOCALE=zh` giữ asset upstream nguyên vẹn;
- apply locale nhiều lần không làm prompt phình lặp;
- global/book style override vẫn thắng localized builtin.
