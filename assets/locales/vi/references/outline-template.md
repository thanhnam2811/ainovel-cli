# Template quy hoạch outline

Template này không ép mọi tác phẩm vào một độ dài cố định. Mục tiêu là xác định **scale của tác phẩm trước**, rồi chọn độ hạt outline phù hợp.

## Bước 1: Xác định cấp độ dài

### Truyện ngắn / một volume

- Phù hợp khi: một xung đột chính, một mục tiêu, ít nhân vật, kết cục tập trung.
- Scale tham khảo: khoảng 8–25 chương.
- Format khuyến nghị: `outline` phẳng.

### Truyện trung bình / nhiều stage

- Phù hợp khi: có tăng tiến theo giai đoạn, vài tuyến phụ, quan hệ nhân vật thay đổi.
- Scale tham khảo: khoảng 25–60 chương.
- Format: `outline` phẳng hoặc phân tầng nhẹ.

### Truyện dài serial / web novel

- Phù hợp khi: thể loại có không gian tăng tiến lâu dài, tension quan hệ dài hạn, nhiều stage goal, world mở rộng được, bí ẩn dài hoặc growth line dài.
- Scale chỉ mang tính tham khảo; không khóa tổng chương từ đầu.
- Format khuyến nghị: `layered_outline`.

## Bước 2: Khi nào ưu tiên layered outline

Nếu có từ hai dấu hiệu sau trở lên, ưu tiên `layered_outline`:

- world cần mở dần, không thể info-dump một lần;
- main tăng trưởng qua nhiều stage chứ không một cú nhảy;
- quan hệ biến đổi xuyên nhiều giai đoạn;
- mid/late game có loại xung đột khác đầu truyện;
- có nhiều lần đổi map/thế lực/thân phận/goal;
- tác phẩm rõ ràng mang logic serial thương mại hơn là một câu chuyện khép trong một volume.

## Bước 3: Truyện dài không làm “bảng kê chương toàn sách” ngay từ đầu

Thứ tự quy hoạch nên là:

1. lời hứa/điểm bán và khác biệt hóa;
2. story engine dài hạn;
3. theme/chức năng cấp volume;
4. goal và chuyển ngoặt cấp arc;
5. event + hook cấp chương khi arc sắp được viết.

Cách sai:

- viết trước một số chương đầu rồi kéo giãn cưỡng ép;
- volume nào cũng “gặp địch → mạnh hơn → đổi map”;
- chỉ tăng progression mà không tăng complexity quan hệ;
- tiêu hết bí mật lớn quá sớm rồi hậu kỳ chỉ còn lặp công thức.

## Template outline phẳng — truyện ngắn/trung bình

```json
[
  {
    "chapter": 1,
    "title": "tiêu đề chương",
    "core_event": "sự kiện cốt lõi của chương",
    "hook": "hook cuối chương",
    "scenes": ["cảnh 1", "cảnh 2", "cảnh 3"]
  }
]
```

## Template layered outline — truyện dài, rolling expansion theo volume/arc

Ban đầu chỉ chi tiết đủ gần để viết. Volume/arc tương lai giữ skeleton và mở dần theo story fact.

```json
[
  {
    "index": 1,
    "title": "tiêu đề volume 1",
    "theme": "xung đột/theme mới mà volume này mang vào",
    "arcs": [
      {
        "index": 1,
        "title": "arc đầu đã mở",
        "goal": "goal cục bộ, lực cản và chuyển ngoặt",
        "chapters": [
          {
            "chapter": 1,
            "title": "tiêu đề chương",
            "core_event": "sự kiện cốt lõi",
            "hook": "hook cuối chương",
            "scenes": ["cảnh 1", "cảnh 2"]
          }
        ]
      },
      {
        "index": 2,
        "title": "arc skeleton",
        "goal": "tóm tắt goal arc",
        "estimated_chapters": 12,
        "chapters": []
      }
    ]
  },
  {
    "index": 2,
    "title": "tiêu đề volume 2",
    "theme": "theme volume 2",
    "arcs": [
      {"index": 1, "title": "arc", "goal": "goal", "estimated_chapters": 15, "chapters": []},
      {"index": 2, "title": "arc", "goal": "goal", "estimated_chapters": 10, "chapters": []}
    ]
  }
]
```

`estimated_chapters` là ước lượng nhịp, không phải hợp đồng tổng chương. Khi mở arc, phải điều chỉnh theo story fact và độ dài chương người dùng muốn.

## Checklist cấp volume

Mỗi volume phải trả lời:

- Nó thêm world information gì thật sự hữu ích?
- Nó nâng cấp hoặc đổi loại xung đột cốt lõi thế nào?
- Main đạt được gì và mất gì?
- Quan hệ chính thay đổi ra sao?
- Sau volume này, vì sao truyện **phải** bước sang trạng thái mới thay vì quay lại loop cũ?

## Checklist cấp arc

- Goal rõ là gì?
- Lực cản đến từ ai, rule nào hoặc cái giá nào?
- Chuyển ngoặt chính là gì?
- Kết thúc arc, trạng thái nào đã thay đổi không đảo ngược?

## Checklist cấp chương

- Mỗi chương phục vụ goal của arc.
- Mỗi chương có ít nhất một thay đổi/sự kiện không thể xóa mà vẫn giữ nguyên tiến triển.
- Hook phải đa dạng, không lặp một kiểu “phát hiện bí mật”.
- Chương đầu không chỉ giới thiệu world; phải đồng thời đẩy nhân vật hoặc xung đột.
- Lượng `core_event`/`scenes` phải khớp `working_memory.user_rules.preferences`, không ép fixed plot density vào mọi độ dài.
