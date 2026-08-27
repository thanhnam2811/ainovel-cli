# ainovel-cli — Vietnamese upstream-friendly fork

Bản fork tiếng Việt của [`voocel/ainovel-cli`](https://github.com/voocel/ainovel-cli), giữ engine upstream gần nguyên vẹn và phủ lớp Việt hóa ở creative/semantic asset boundary: prompt, reference, voice và style.

> **Trạng thái release:** fork hiện chưa phát hành GitHub Release riêng. Để tránh cài nhầm binary upstream không có lớp Việt hóa, **chưa dùng installer/release command** cho tới khi release đầu tiên của fork được publish sau khi CI xanh.

## Điểm khác của fork này

- mặc định `AINOVEL_LOCALE=vi`;
- toàn bộ runtime prompt hiện tại đã có overlay tiếng Việt;
- runtime reference pack và các preset `default`, `fantasy`, `romance`, `suspense` đã được bản địa hóa;
- tool name, JSON field, enum và protocol marker vẫn giữ nguyên để không phá engine;
- `AINOVEL_LOCALE=zh` chạy nguyên asset upstream để đối chiếu/debug;
- global/book style override vẫn thắng localized builtin;
- module path vẫn giữ `github.com/voocel/ainovel-cli` để giảm conflict khi sync upstream.

**Phạm vi có chủ ý:** đây là bản địa hóa lớp sáng tác/ngữ nghĩa, không phải bản dịch toàn bộ source/UI. Các label/log/task nội bộ thuộc deterministic engine/TUI vẫn có thể dùng tiếng Trung upstream nếu việc dịch chúng không cần thiết cho correctness. Những boundary có ảnh hưởng chức năng — parser, eval, installer, updater, release, Docker — được test riêng trong fork.

Chi tiết thiết kế, phạm vi localization và policy sync: **[README.vi.md](README.vi.md)**.

## Cài và chạy hiện tại

Fork chưa có release binary riêng, vì vậy build từ source là đường cài an toàn:

```bash
git clone https://github.com/thanhnam2811/ainovel-cli.git
cd ainovel-cli

GOWORK=off go build -o ainovel-cli ./cmd/ainovel-cli
./ainovel-cli
```

Repo dùng toolchain khai báo trong `go.mod`; nếu bản Go bootstrap hỗ trợ toolchain auto-download, nó sẽ lấy đúng version cần thiết.

Không chạy:

```bash
go install github.com/voocel/ainovel-cli/cmd/ainovel-cli@latest
```

lúc muốn dùng bản Việt — lệnh đó cài **upstream**.

Sau khi fork có GitHub Release riêng, `scripts/install.sh` và lệnh `ainovel-cli update` đã được cấu hình để dùng release của `thanhnam2811/ainovel-cli`.

## Dùng locale upstream để đối chiếu

```bash
AINOVEL_LOCALE=zh ./ainovel-cli
```

Mặc định không cần set biến môi trường:

```bash
./ainovel-cli
```

sẽ dùng creative/semantic asset tiếng Việt.

## Kiến trúc và tài liệu upstream

Fork không copy/Việt hóa toàn bộ comment và tài liệu engine vì việc đó làm tăng maintenance debt. Kiến trúc cốt lõi, flow, store/checkpoint, import/revision/eval và lịch sử thiết kế vẫn theo upstream:

- [Upstream repository](https://github.com/voocel/ainovel-cli)
- [Upstream README](https://github.com/voocel/ainovel-cli/blob/main/README.md)
- [`docs/`](docs/)

Downstream diff nên tiếp tục tập trung ở `assets/locales/`, test contract và các boundary thật sự cần biết locale.

## Validation trước release

```bash
gofmt -l .
GOWORK=off go vet ./...
GOWORK=off go test -buildvcs=false -count=1 ./...
GOWORK=off go test -race -buildvcs=false -count=1 ./internal/host ./internal/store ./internal/tools
sh -n scripts/install.sh
sh -n scripts/check-fork-boundaries.sh
sh scripts/check-fork-boundaries.sh
```

Không publish release nếu các check trên chưa xanh.
