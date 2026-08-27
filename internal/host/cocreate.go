package host

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/voocel/agentcore"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/localization"
	"github.com/voocel/ainovel-cli/internal/store"
)

// Keep the upstream Chinese prompts intact for AINOVEL_LOCALE=zh. The
// Vietnamese fork selects the localized equivalents at process start.
var (
	coCreateSystemPrompt      = selectCoCreatePrompt(zhCoCreateSystemPrompt, viCoCreateSystemPrompt)
	stageCoCreateSystemPrompt = selectCoCreatePrompt(zhStageCoCreateSystemPrompt, viStageCoCreateSystemPrompt)
)

func selectCoCreatePrompt(upstream, vietnamese string) string {
	if localization.IsVietnamese() {
		return vietnamese
	}
	return upstream
}

const zhCoCreateSystemPrompt = `你是一个小说共创助手。你的任务不是直接开始写小说，而是通过多轮简短对话帮助用户澄清创作需求，并持续整理出一段可直接交给创作引擎的中文创作指令。

每一轮回复严格按以下 XML 格式输出，包含四个标签，依次出现，每个标签都必须有正确的开闭标签：

<reply>
给用户看的中文自然回复：先回应用户的输入，再最多提出 1 到 2 个当前最关键的问题。如果信息已足够开始创作，告诉用户可以按 Ctrl+S 开始。
</reply>

<draft>
当前完整的创作指令草稿，使用 Markdown：直接从二级标题开始，例如 "## 主题"、"## 关键要素"、"## 待澄清信息"；用项目符号列出要点。每一轮都要在已有结论上**累积更新**，吸收用户最新意图；即使本轮没有新增也要把完整草稿原样再写一次——不要省略、不要写"（保持上一轮）"之类的占位。
</draft>
` + zhCoCreateProtocolTail

const zhStageCoCreateSystemPrompt = `你是一个小说"阶段共创"助手。这本小说已经写了一部分（进度见下方"当前故事状态"）。用户暂停下来，想和你一起规划"后续阶段"的走向，再继续创作。

你的任务不是续写正文，而是通过多轮简短对话帮用户想清楚后面这一段（接下来若干章 / 下一弧 / 下一卷）要往哪走，并持续整理出一段"后续方向 brief"，供创作引擎据此推进。

铁律：所有建议必须与"当前故事状态"里已发生的剧情、人物、伏笔一致，绝不推翻或忽略已写内容；只规划"后续怎么走"，不重新设计整本书。

每一轮回复严格按以下 XML 格式输出，包含四个标签，依次出现，每个标签都必须有正确的开闭标签：

<reply>
给用户看的中文自然回复：先回应用户的输入，再最多提出 1 到 2 个当前最关键的问题。如果后续方向已足够清晰，告诉用户可以按 Ctrl+S 把方向交给创作引擎、继续创作。
</reply>

<draft>
当前完整的"后续方向 brief"，使用 Markdown：直接从二级标题开始，例如 "## 后续走向"、"## 关键转折"、"## 要收的伏笔"、"## 节奏与篇幅"；用项目符号列出要点。每一轮都要在已有结论上**累积更新**，吸收用户最新意图；即使本轮没有新增也要把完整 brief 原样再写一次——不要省略、不要写"（保持上一轮）"之类的占位。
</draft>
` + zhCoCreateProtocolTail

const zhCoCreateProtocolTail = `
<ready>false</ready>

<suggestions>
1-3 条"用户接下来可能想说的话"，每行一条以 "- " 开头。这是用户卡壳时的引导，
按数字键填入输入框，用户可再编辑后发送。

要求：
- 站在用户口吻，像用户对你说的话，不要写成助手反问。
- 每条不超过 25 字，多样化句式，避免千篇一律。
- 给倾向 / 选择 / 补充意图，不要一句话替用户写完整设定。
</suggestions>

输出规范：
- 必须使用四个 XML 标签：<reply> / <draft> / <ready> / <suggestions>，每个都必须完整开闭。
- 标签名只能小写英文，不要改写成 <REPLY> / <REWRITE> / <回复> 等任何变体。
- 标签外不要添加任何说明、思考或代码围栏。
- <draft> 内允许多行 Markdown，直接换行书写，不需要任何转义。
- <ready> 只写 true 或 false。信息已足够时填 true。
- <ready>true</ready> 时 <suggestions> 可以为空（保留空标签 <suggestions></suggestions> 即可）。`

const viCoCreateSystemPrompt = `Bạn là trợ lý đồng sáng tác tiểu thuyết. Nhiệm vụ của bạn không phải bắt đầu viết truyện ngay, mà là trao đổi ngắn qua nhiều lượt để giúp người dùng làm rõ yêu cầu và liên tục tổng hợp thành một chỉ dẫn sáng tác tiếng Việt có thể đưa thẳng cho engine.

Mỗi lượt trả lời phải dùng đúng bốn thẻ XML dưới đây, theo đúng thứ tự và có đủ thẻ mở/đóng:

<reply>
Phản hồi tự nhiên bằng tiếng Việt: trước hết đáp lại ý người dùng, sau đó hỏi tối đa 1-2 câu quan trọng nhất ở thời điểm hiện tại. Nếu thông tin đã đủ để bắt đầu sáng tác, nói rõ người dùng có thể nhấn Ctrl+S để bắt đầu.
</reply>

<draft>
Bản chỉ dẫn sáng tác đầy đủ hiện tại bằng Markdown, bắt đầu trực tiếp từ heading cấp 2, ví dụ "## Chủ đề", "## Yếu tố then chốt", "## Thông tin còn cần làm rõ". Dùng bullet cho các ý chính. Mỗi lượt phải cập nhật tích lũy trên toàn bộ kết luận trước đó và hấp thụ ý mới nhất của người dùng; dù lượt này không có thông tin mới, vẫn phải xuất lại đầy đủ bản draft, không dùng placeholder kiểu "giữ nguyên như trước".
</draft>
` + viCoCreateProtocolTail

const viStageCoCreateSystemPrompt = `Bạn là trợ lý "đồng sáng tác theo giai đoạn" cho một tiểu thuyết đã được viết một phần. Người dùng đang tạm dừng để cùng bạn định hướng giai đoạn tiếp theo rồi mới tiếp tục sáng tác.

Nhiệm vụ không phải viết tiếp正文, mà là trao đổi ngắn qua nhiều lượt để giúp người dùng chốt hướng đi cho một số chương tiếp theo / arc tiếp theo / volume tiếp theo, đồng thời duy trì một brief đầy đủ để engine dùng khi tiếp tục.

Nguyên tắc bắt buộc: mọi đề xuất phải nhất quán với những sự kiện, nhân vật và foreshadow đã xảy ra trong phần "Trạng thái câu chuyện hiện tại". Không lật lại hoặc bỏ qua phần đã viết; chỉ hoạch định hướng đi tiếp theo, không thiết kế lại toàn bộ tác phẩm.

Mỗi lượt trả lời phải dùng đúng bốn thẻ XML dưới đây, theo đúng thứ tự và có đủ thẻ mở/đóng:

<reply>
Phản hồi tự nhiên bằng tiếng Việt: trước hết đáp lại ý người dùng, sau đó hỏi tối đa 1-2 câu quan trọng nhất. Nếu hướng đi tiếp theo đã đủ rõ, nói người dùng có thể nhấn Ctrl+S để gửi brief cho engine và tiếp tục sáng tác.
</reply>

<draft>
Brief đầy đủ của giai đoạn tiếp theo bằng Markdown, bắt đầu trực tiếp từ heading cấp 2, ví dụ "## Hướng đi tiếp theo", "## Bước ngoặt chính", "## Foreshadow cần thu hồi", "## Nhịp và dung lượng". Dùng bullet cho các ý chính. Mỗi lượt phải cập nhật tích lũy và xuất lại toàn bộ brief; không dùng placeholder kiểu "giữ nguyên như trước".
</draft>
` + viCoCreateProtocolTail

const viCoCreateProtocolTail = `
<ready>false</ready>

<suggestions>
Đưa ra 1-3 câu mà người dùng có thể muốn nói tiếp, mỗi câu một dòng bắt đầu bằng "- ". Đây là gợi ý khi người dùng bí ý và có thể được chèn vào ô nhập để chỉnh sửa trước khi gửi.

Yêu cầu:
- Viết ở ngôi của người dùng, như lời người dùng nói với trợ lý; không biến thành câu hỏi của trợ lý.
- Mỗi gợi ý ngắn gọn, tối đa khoảng 25 từ, đa dạng cách diễn đạt.
- Chỉ gợi ý khuynh hướng, lựa chọn hoặc ý bổ sung; không tự viết thay toàn bộ thiết lập.
</suggestions>

Quy cách output:
- Bắt buộc dùng đúng bốn thẻ XML: <reply> / <draft> / <ready> / <suggestions>, mỗi thẻ phải đóng đầy đủ.
- Tên thẻ chỉ dùng chữ thường tiếng Anh; không đổi thành biến thể khác.
- Không thêm giải thích, suy nghĩ hay code fence bên ngoài các thẻ.
- <draft> được phép chứa Markdown nhiều dòng, không cần escape newline.
- <ready> chỉ được là true hoặc false.
- Khi <ready>true</ready>, <suggestions> có thể rỗng nhưng vẫn phải giữ cặp thẻ <suggestions></suggestions>.`

// CoCreateProgressKind 标识流式回调的内容类型。
const (
	CoCreateProgressThinking = "thinking"
	CoCreateProgressReply    = "reply"
)

const (
	tagReply       = "reply"
	tagDraft       = "draft"
	tagReady       = "ready"
	tagSuggestions = "suggestions"
)

func coCreateStream(ctx context.Context, models *bootstrap.ModelSet, sessions *store.SessionStore, sysPrompt string, history []CoCreateMessage, onProgress func(kind, text string)) (reply CoCreateReply, err error) {
	if len(history) == 0 {
		return CoCreateReply{}, fmt.Errorf("cocreate history is empty")
	}

	model := models.ForRole("thinking")

	msgs := []agentcore.Message{agentcore.SystemMsg(sysPrompt)}
	for _, item := range history {
		content := strings.TrimSpace(item.Content)
		if content == "" {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(item.Role)) {
		case "assistant":
			msgs = append(msgs, assistantMsg(content))
		default:
			msgs = append(msgs, agentcore.UserMsg(content))
		}
	}

	var raw, thinking strings.Builder
	start := time.Now()
	defer func() {
		if sessions == nil {
			return
		}
		if logErr := sessions.LogCoCreate(coCreateLogEntry{
			Time:         time.Now(),
			DurationMS:   time.Since(start).Milliseconds(),
			InputHistory: history,
			RawResponse:  raw.String(),
			RawLen:       len([]rune(raw.String())),
			Thinking:     thinking.String(),
			ParsedReply:  reply.Message,
			ParsedDraft:  reply.Prompt,
			ParsedReady:  reply.Ready,
			ParsedSugs:   reply.Suggestions,
			Error:        errString(err),
		}); logErr != nil {
			slog.Warn("共创会话日志落盘失败", "module", "cocreate", "err", logErr)
		}
	}()

	streamCh, err := model.GenerateStream(ctx, msgs, nil, agentcore.WithMaxTokens(2048))
	if err != nil {
		return CoCreateReply{}, fmt.Errorf("cocreate generate: %w", err)
	}

	var streamed bool
	for ev := range streamCh {
		switch ev.Type {
		case agentcore.StreamEventThinkingDelta:
			thinking.WriteString(ev.Delta)
			if onProgress != nil {
				onProgress(CoCreateProgressThinking, thinking.String())
			}
		case agentcore.StreamEventTextDelta:
			streamed = true
			raw.WriteString(ev.Delta)
			if onProgress != nil {
				onProgress(CoCreateProgressReply, extractReplyPreview(raw.String()))
			}
		case agentcore.StreamEventDone:
			if !streamed {
				raw.WriteString(ev.Message.TextContent())
			}
		case agentcore.StreamEventError:
			if ev.Err != nil {
				return CoCreateReply{}, fmt.Errorf("cocreate generate: %w", ev.Err)
			}
			return CoCreateReply{}, fmt.Errorf("cocreate generate failed")
		}
	}

	rawText := raw.String()
	if strings.TrimSpace(rawText) == "" {
		if t := strings.TrimSpace(thinking.String()); t != "" {
			rawText = t
		}
	}
	reply, err = parseCoCreateResponse(rawText)
	return reply, err
}

type coCreateLogEntry struct {
	Time         time.Time         `json:"time"`
	DurationMS   int64             `json:"duration_ms"`
	InputHistory []CoCreateMessage `json:"input_history"`
	RawResponse  string            `json:"raw_response"`
	RawLen       int               `json:"raw_len"`
	Thinking     string            `json:"thinking,omitempty"`
	ParsedReply  string            `json:"parsed_reply"`
	ParsedDraft  string            `json:"parsed_draft"`
	ParsedReady  bool              `json:"parsed_ready"`
	ParsedSugs   []string          `json:"parsed_sugs,omitempty"`
	Error        string            `json:"error,omitempty"`
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func assistantMsg(text string) agentcore.Message {
	return agentcore.Message{
		Role:      agentcore.RoleAssistant,
		Content:   []agentcore.ContentBlock{agentcore.TextBlock(text)},
		Timestamp: time.Now(),
	}
}

func parseCoCreateResponse(raw string) (CoCreateReply, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return CoCreateReply{}, fmt.Errorf("cocreate empty response")
	}

	reply, draft, ready, suggestions := splitCoCreateMarkers(raw)
	if reply == "" {
		return CoCreateReply{Message: raw, Prompt: "", Ready: false, Raw: raw}, nil
	}
	return CoCreateReply{
		Message:     reply,
		Prompt:      draft,
		Ready:       ready,
		Suggestions: suggestions,
		Raw:         raw,
	}, nil
}

func splitCoCreateMarkers(s string) (reply, draft string, ready bool, suggestions []string) {
	reply = extractTagContent(s, tagReply)
	draft = extractTagContent(s, tagDraft)
	readyStr := strings.ToLower(extractTagContent(s, tagReady))
	ready = readyStr == "true" || readyStr == "yes"
	suggestions = parseSuggestions(extractTagContent(s, tagSuggestions))
	return
}

func extractTagContent(s, tag string) string {
	open := "<" + tag + ">"
	closeTag := "</" + tag + ">"
	oIdx := strings.Index(s, open)
	if oIdx >= 0 {
		rest := s[oIdx+len(open):]
		if cIdx := strings.Index(rest, closeTag); cIdx >= 0 {
			return strings.TrimSpace(rest[:cIdx])
		}
		for _, other := range []string{"<reply>", "<draft>", "<ready>", "<suggestions>"} {
			if other == open {
				continue
			}
			if idx := strings.Index(rest, other); idx >= 0 {
				rest = rest[:idx]
			}
		}
		return strings.TrimSpace(rest)
	}

	if cIdx := strings.Index(s, closeTag); cIdx >= 0 {
		prefix := s[:cIdx]
		start := 0
		for _, t := range []string{"</reply>", "</draft>", "</ready>", "</suggestions>"} {
			if t == closeTag {
				continue
			}
			if i := strings.LastIndex(prefix, t); i >= 0 {
				if end := i + len(t); end > start {
					start = end
				}
			}
		}
		return strings.TrimSpace(prefix[start:])
	}
	return ""
}

func parseSuggestions(text string) []string {
	if text == "" {
		return nil
	}
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "- "):
			line = strings.TrimSpace(line[2:])
		case strings.HasPrefix(line, "* "):
			line = strings.TrimSpace(line[2:])
		case isOrderedSuggestion(line):
			line = stripOrderedPrefix(line)
		}
		if len([]rune(line)) < 2 {
			continue
		}
		out = append(out, line)
		if len(out) >= 3 {
			break
		}
	}
	return out
}

func isOrderedSuggestion(line string) bool {
	i := 0
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	return i > 0 && i+1 < len(line) && line[i] == '.' && line[i+1] == ' '
}

func stripOrderedPrefix(line string) string {
	i := 0
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i == 0 || i+1 >= len(line) {
		return line
	}
	return strings.TrimSpace(line[i+2:])
}

func extractReplyPreview(raw string) string {
	trimmed := strings.TrimSpace(raw)
	open := "<" + tagReply + ">"
	closeTag := "</" + tagReply + ">"
	draftOpen := "<" + tagDraft + ">"

	rest := trimmed
	if rIdx := strings.Index(trimmed, open); rIdx >= 0 {
		rest = trimmed[rIdx+len(open):]
	}
	if cIdx := strings.Index(rest, closeTag); cIdx >= 0 {
		return strings.TrimSpace(rest[:cIdx])
	}
	if dIdx := strings.Index(rest, draftOpen); dIdx >= 0 {
		rest = rest[:dIdx]
	}
	return strings.TrimSpace(rest)
}
