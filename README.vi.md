# ainovel-cli — bản Việt hóa upstream-friendly

Fork này bám sát `voocel/ainovel-cli` và thêm lớp Việt hóa dạng **overlay**, thay vì dịch hàng loạt source Go.

Mục tiêu: giữ engine/state machine/checkpoint/import pipeline/bugfix upstream gần nguyên vẹn, nhưng toàn bộ lớp prompt + reference tác động trực tiếp tới agent sáng tác bằng tiếng Việt tự nhiên.

## Cách hoạt động

Runtime nạp asset upstream trước rồi `assets.ApplyLocale` phủ asset Việt trong:

```text
assets/locales/vi/
├── prompts/
├── references/
├── styles/
└── voice.md
```

Tên tool, JSON field, enum, protocol marker và placeholder như `novel_context`, `commit_chapter`, `required_beats`, `{{VOICE}}` luôn giữ nguyên.

Nếu tương lai upstream thêm prompt mới hoặc một overlay phải tạm gỡ vì chưa port kịp protocol, loader giữ prompt upstream mới nhất và chỉ thêm chỉ thị output tiếng Việt. **Prompt upstream mới nhưng chưa dịch an toàn tốt hơn một bản dịch stale đã lệch contract.**

## Locale mặc định

```bash
ainovel-cli
```

Fork mặc định là tiếng Việt. Có thể chạy nguyên asset upstream để đối chiếu/debug:

```bash
AINOVEL_LOCALE=zh ainovel-cli
```

`vi` và `vi-VN` đều được hỗ trợ.

## Phạm vi Việt hóa

### Runtime prompts — đầy đủ

- Architect Short / Architect Long
- Writer / Editor
- import-segment / import-analyze / import-synthesize / import-range
- revision-analyze
- simulation-source / simulation-merge
- arbiter-plan-start / arbiter-intervention / arbiter-failure

### Reference pack mà `loadReferences()` thật sự nạp — đầy đủ

- chapter-guide
- hook-techniques
- quality-checklist
- outline-template
- character-template
- chapter-template
- consistency
- content-expansion
- dialogue-writing
- longform-planning
- differentiation
- anti-ai-tone

Các asset upstream hiện **không được runtime loader dùng** như `character-building.md` hay `plot-structures.md` không được dịch chỉ để tăng coverage giả.

### Style

- writing voice: đã Việt hóa
- default writing style: đã Việt hóa
- user global/book style override vẫn giữ nguyên precedence upstream

## Adaptation tiếng Việt, không dịch máy heuristic tiếng Trung

Reference được bản địa hóa theo **ý nghĩa và chức năng**, không thay máy từng chữ.

Một số thứ phụ thuộc tiếng Trung được chủ động điều chỉnh:

- ngưỡng `字数` không được đổi máy thành “số từ tiếng Việt”;
- độ dài chương ưu tiên `working_memory.user_rules.preferences` và lượng beat/cảnh;
- ví dụ hội thoại/dấu câu được đổi sang tiếng Việt tự nhiên;
- heuristic tiếng Trung chưa có calibration tiếng Việt không được giả precision bằng threshold mới.

Chi tiết: `assets/locales/vi/references/README.md`.

## Guard cho protocol

Smoke test khóa các token/invariant máy đọc thay vì khóa wording bản dịch:

- Architect: `save_foundation`, `audit_foundation`, `layered_outline`, `append_volume`, `complete_book`, `expand_arc`, `final_volume`…
- Editor: `save_review`, `requires_change`, `rule_violations`, `accept`, `polish`, `rewrite`.
- Import: enum/value domain của `hook_type`, `dominant_strand`, `planning_tier`, `story_status`, boundary/foreshadow action.
- Arbiter: `answer → rules → hold → reopen → dispatch`, `rewrites_drained`, `phase = complete`, `dynamic_planning`, tập agent hợp lệ.
- Reference pack: test xác nhận mọi field runtime-loaded thật sự nhận bản Việt.

Checklist port prompt nằm trong `assets/locales/vi/prompts/README.md` và `CONTRACT.md`.

## Override của người dùng vẫn thắng

Với `voice.md`, `anti-ai-tone.md` và `styles/*.md`:

```text
Vietnamese builtin < ~/.ainovel/style < <book>/style
```

Locale chỉ thay builtin. Cấu hình global và cấu hình riêng từng sách vẫn thắng như upstream.

## Rules tiếng Việt

`~/.ainovel/rules/*.md` và `./.ainovel/rules/*.md` vốn là natural-language input nên có thể viết tiếng Việt trực tiếp.

Mechanical baseline upstream cho sáo/từ mòn tiếng Trung có threshold từ corpus thực tế. Fork **không dịch máy threshold đó sang tiếng Việt**; baseline cơ học tiếng Việt chỉ nên thêm sau khi có corpus generated thật để đo frequency/false-positive.

## Cài đặt và self-update

Installer và `update` đều trỏ release của `thanhnam2811/ainovel-cli`, không tải binary `voocel` đè mất localization.

```bash
curl -fsSL https://raw.githubusercontent.com/thanhnam2811/ainovel-cli/main/scripts/install.sh | sh
```

Trước release riêng đầu tiên, build trực tiếp từ source là chắc chắn nhất.

## Sync upstream

Không đổi module path `github.com/voocel/ainovel-cli`.

Khi upstream đổi prompt/reference đã dịch:

1. đọc diff upstream;
2. port thay đổi **behavior/ý nghĩa**, không chỉ wording;
3. giữ machine-readable contract;
4. cập nhật smoke test nếu upstream thêm invariant;
5. nếu chưa port an toàn, tạm fallback upstream thay vì phát hành bản Việt stale.

Diff downstream nên tiếp tục tập trung dưới `assets/locales/` và một hook nhỏ ở startup/release boundary.

## Validation

Trước merge/release:

```bash
go test ./...
```

Test locale ít nhất phải bảo đảm:

- bản Việt được load thật;
- `{{VOICE}}` không mất;
- prompt không rơi protocol/routing token;
- reference runtime-loaded nhận bản Việt;
- fallback vẫn hoạt động;
- `AINOVEL_LOCALE=zh` là no-op trên asset upstream;
- apply locale idempotent;
- global/book override thắng localized builtin.
