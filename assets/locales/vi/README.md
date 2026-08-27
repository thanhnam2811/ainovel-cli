# Vietnamese locale assets

Các file trong thư mục này là lớp dịch downstream cho `voocel/ainovel-cli`.

Nguyên tắc bảo trì:

- giữ nguyên tên tool, JSON field, enum, protocol marker và placeholder;
- không đổi schema chỉ để câu chữ tiếng Việt tự nhiên hơn;
- khi upstream sửa một prompt đã được dịch đầy đủ, phải review diff upstream trước khi phát hành bản Việt tiếp theo;
- prompt chưa dịch tiếp tục dùng fallback upstream + chỉ thị output tiếng Việt;
- ưu tiên dịch các prompt ảnh hưởng trực tiếp tới hành vi agent trước, không dịch source Go hàng loạt.
