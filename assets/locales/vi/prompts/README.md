# Prompt localization contract

Bản dịch prompt phải giữ nguyên các contract mà host/tool layer phụ thuộc vào.

Checklist khi cập nhật từ upstream:

1. Đọc diff upstream của prompt gốc.
2. Port đầy đủ thay đổi hành vi, không chỉ dịch câu chữ.
3. Giữ nguyên tool names, field names, enum values, schema examples, protocol markers và placeholder.
4. Không đổi thứ tự workflow nếu upstream dùng thứ tự đó như invariant.
5. Cập nhật smoke test nếu upstream thêm protocol token quan trọng mới.
6. Nếu chưa kịp port an toàn, xóa/không thêm overlay của prompt đó để runtime fallback về prompt upstream + chỉ thị output tiếng Việt.
