# ainovel-cli — Vietnamese upstream-friendly fork

Bản fork tiếng Việt của [`voocel/ainovel-cli`](https://github.com/voocel/ainovel-cli), giữ engine upstream gần nguyên vẹn và phủ lớp Việt hóa ở prompt/reference/style boundary.

> **Trạng thái release:** fork hiện chưa phát hành GitHub Release riêng. Để tránh cài nhầm binary upstream không có lớp Việt hóa, **chưa dùng installer/release command** cho tới khi release đầu tiên của fork được publish sau khi CI xanh.

## Điểm khác của fork này

- mặc định `AINOVEL_LOCALE=vi`;
- toàn bộ runtime prompt hiện tại đã có overlay tiếng Việt;
- runtime reference pack và các preset `default`, `fantasy`, `romance`, `suspense` đã được bản địa hóa;
- tool name, JSON field, enum và protocol marker vẫn giữ nguyên để không phá engine;
- `AINOVEL_LOCALE=zh` chạy nguyên asset upstream để đối chiếu/debug;
- global/book style override vẫn thắng localized builtin;
- module path vẫn giữ `github.com/voocel/ainovel-cli` để giảm conflict khi sync upstream.

Chi tiết thiết kế, phạm vi localization và policy sync: **[README.vi.md](README.vi.md)**.

## Cài và chạy hiện tại

Fork chưa có release binary riêng, vì vậy build từ source là đường cài an toàn:

```bash
git clone https://github.com/thanhnam2811/ainovel-cli.git
cd ainovel-cli

go build -o ainovel-cli ./cmd/ainovel-cli
./ainovel-cli
```

Repo đang theo toolchain được khai báo trong `go.mod`; nếu Go hỗ trợ toolchain auto-download, nó sẽ dùng đúng version cần thiết.

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

sẽ chạy locale tiếng Việt.

## Kiến trúc và tài liệu upstream

Fork không cố copy/Việt hóa toàn bộ comment và tài liệu engine vì việc đó làm tăng maintenance debt. Kiến trúc cốt lõi, flow, store/checkpoint, import/revision/eval và lịch sử thiết kế vẫn theo upstream:

- [Upstream repository](https://github.com/voocel/ainovel-cli)
- [Upstream README](https://github.com/voocel/ainovel-cli/blob/main/README.md)
- [`docs/`](docs/)

Downstream diff chủ yếu nên tiếp tục tập trung ở `assets/locales/`, test contract, và các boundary thật sự cần biết locale.

## Validation trước release

```bash
gofmt -l .
go vet ./...
go test -buildvcs=false -count=1 ./...
go test -race -count=1 ./internal/host ./internal/store ./internal/tools
bash -n scripts/install.sh
```

Không publish release nếu các check trên chưa xanh.
