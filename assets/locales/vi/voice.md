## Tiêu chuẩn viết

Đây là các tiêu chí chất lượng, không phải checklist để đánh dấu một cách máy móc. Trước hết chương phải tự nhiên và đứng vững như một phần của câu chuyện; sau đó mới xét việc đáp ứng đầy đủ các tiêu chí.

- Mở đầu sớm thiết lập xung đột, bí ẩn, ham muốn hoặc cảm giác bất thường; hạn chế hồi tưởng trừu tượng.
- Dùng hành động, hội thoại và chi tiết cảm giác để đẩy diễn biến; hạn chế kể tóm tắt hoặc giải thích lại.
- Lời thoại phải phản ánh khác biệt về thân phận, tính cách, ẩn ý và mục đích hành động; tránh biến nhân vật thành người giảng đạo.
- Thể hiện cảm xúc qua phản ứng cơ thể và lựa chọn, không dán nhãn trực tiếp cho cảm xúc.
- Biến chuyển quan hệ phải được kích hoạt bởi sự kiện; không để hai người từ xa lạ nhảy thẳng sang tin tưởng tuyệt đối chỉ trong một chương nếu chưa có đủ nguyên nhân.
- Tiết lộ bí mật theo từng lớp; không giải thích sớm các nút thắt lớn mà outline chưa yêu cầu mở.
- Hook cuối chương có thể là nguy cơ, một lựa chọn, dư âm cảm xúc, biến chuyển quan hệ hoặc mục tiêu chưa hoàn tất; không cần chương nào cũng kết bằng cliffhanger phô trương.
- **Giảm “mùi AI”**: khi viết, tránh toàn bộ các mẫu trong `reference_pack.references.anti_ai_tone` thuộc năm nhóm cấu trúc, từ ngữ, miêu tả, hội thoại và nhịp. Những từ mòn, câu khuôn và ngưỡng có thể kiểm tra cơ học nằm trong `working_memory.user_rules.structured` và sẽ bị kiểm tra bắt buộc lúc commit.
- **Đa dạng cấu trúc câu**: `episodic_memory.style_stats` nếu có là thống kê bằng code trên phần bạn đã viết — một tấm gương phản chiếu thói quen lặp của chính bạn. Chủ động giảm các mẫu có tần suất cao trong chương hiện tại. Các nguồn lặp phổ biến gồm câu đối lập chỉnh nghĩa kiểu “không phải… mà là…”, một đơn vị đo thời gian bị dùng liên tục, hoặc nhiều phép so sánh cùng khuôn. Luân phiên cách kết chương giữa câu ngắn dứt, dư âm hội thoại, hình ảnh còn đọng, câu hỏi bỏ ngỏ; tránh chương nào cũng mở bằng mô-típ thời gian như “đêm xuống”, “sáng hôm sau”, “khi tỉnh dậy”.
- **Không kể lại chương trước**: các summary, foreshadowing và trạng thái trong `episodic_memory` là dữ liệu để kiểm tra tính nối tiếp, không phải nguyên liệu cần chép lại vào chương mới. Thông tin đã được nói ở chương trước chỉ nên được chạm lại khi tình huống mới thực sự cần và từ một góc nhìn mới. Tránh kiểu “tóm tắt tập trước”; câu bị lặp nguyên văn qua nhiều chương sẽ bị `style_stats.repeated_sentences` ghi nhận.