# Vietnamese reference-pack policy

Các reference trong thư mục này là bản địa hóa ngữ nghĩa của reference upstream, không phải bản dịch máy từng chữ.

Nguyên tắc:

- giữ nguyên mục đích và logic sáng tác của upstream;
- không thay schema/tool/protocol — reference không được phép định nghĩa contract mới;
- ví dụ câu chữ, dấu câu hội thoại và idiom được đổi sang tiếng Việt tự nhiên;
- các ngưỡng phụ thuộc đơn vị `字`/thói quen tiếng Trung không được đổi máy thành “từ tiếng Việt”; độ dài chương phải theo `working_memory.user_rules.preferences` và mật độ beat/cảnh;
- nếu một heuristic tiếng Trung không có calibration tiếng Việt thì ưu tiên tiêu chí ngữ nghĩa thay vì giả lập độ chính xác số học;
- khi upstream sửa reference gốc, review thay đổi về ý nghĩa rồi port sang bản Việt trong cùng phạm vi.
