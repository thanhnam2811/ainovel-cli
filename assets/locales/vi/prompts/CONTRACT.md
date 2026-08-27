# Invariants cần giữ khi Việt hóa prompt

Các invariant sau không được dịch hoặc đổi nghĩa theo cách làm host/tool layer hiểu khác đi:

- tool name;
- JSON field name;
- enum value;
- `type=...` và các operation name;
- protocol marker / placeholder;
- điều kiện hoàn tất và điều kiện retry;
- phạm vi được phép sửa dữ liệu;
- thứ tự workflow khi upstream mô tả như hard constraint.

Nếu upstream thêm invariant mới, cập nhật bản Việt và test contract trong cùng PR.
