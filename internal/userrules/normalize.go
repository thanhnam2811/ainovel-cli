// Package userrules 是用户规则归一化的服务层：把各来源的自然语言规则经 LLM 结构化调用
// 归一化成候选结构化字段，再由 rules.BuildSnapshot 确定性合并成本书快照。
//
// 分层职责：
//   - rules 包：纯数据 + 确定性合并（Snapshot / Candidate / BuildSnapshot / SystemDefaults）
//   - 本包：LLM 归一化 + 编排 + 落盘（依赖 agentcore + store + rules）
//
// 归一化是增强路径，不是主创作的前置条件：任何来源失败都降级为 raw preferences，主创作必须继续。
package userrules

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/voocel/agentcore"
	"github.com/voocel/agentcore/schema"
	"github.com/voocel/ainovel-cli/internal/llmcontract"
	"github.com/voocel/ainovel-cli/internal/localization"
	"github.com/voocel/ainovel-cli/internal/rules"
)

const normalizeMaxTokens = 8192

var normalizeContract = llmcontract.Contract{
	Name:        "userrules_normalize",
	Description: "把用户自然语言写作规则归一化为结构化字段",
	Schema: schema.Object(
		schema.Property("structured", schema.Object(
			schema.Property("genre", schema.String("题材;无则空字符串")).Required(),
			schema.Property("forbidden_chars", schema.Array("禁止出现的字符", schema.String("字符"))).Required(),
			schema.Property("forbidden_phrases", schema.Array("禁止出现的短语(字面精确匹配)", schema.String("短语"))).Required(),
			schema.Property("fatigue_words", schema.Array("疲劳词及每章出现上限", schema.Object(
				schema.Property("word", schema.String("疲劳词")).Required(),
				schema.Property("max_per_chapter", schema.Int("每章出现次数上限(正整数)")).Required(),
			))).Required(),
		)).Required(),
		schema.Property("preferences", schema.String("自然语言风格/人物/审美偏好;无则空字符串")).Required(),
		schema.Property("uncertain", schema.Array("故意未提升到 structured 的项+原因", schema.String("条目"))).Required(),
	),
}

type Normalizer struct {
	model agentcore.ChatModel
}

func NewNormalizer(model agentcore.ChatModel) *Normalizer {
	return &Normalizer{model: model}
}

func (n *Normalizer) Normalize(ctx context.Context, source, text string) (rules.Candidate, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return rules.Candidate{Source: source}, nil
	}
	if n == nil || n.model == nil {
		return rules.Candidate{}, fmt.Errorf("归一化模型未配置")
	}

	out, err := llmcontract.Execute(ctx, n.model, llmcontract.Request[normalizerOutput]{
		Contract:     normalizeContract,
		SystemPrompt: activeNormalizerSystemPrompt(),
		Payload:      text,
		Options:      []agentcore.CallOption{agentcore.WithMaxTokens(normalizeMaxTokens)},
		Validate: func(out *normalizerOutput) error {
			_, err := out.toCandidate(source)
			return err
		},
		Agent: "rules",
		Hooks: llmcontract.Hooks{
			Resolved: func(res llmcontract.Resolution) {
				slog.Debug("规则归一化协议选择", "module", "rules", "source", source,
					"contract", normalizeContract.Name, "structured_mode", res.Mode,
					"capability_source", res.Source, "provider", res.Provider, "model", res.Model,
					"schema_fingerprint", normalizeContract.Fingerprint())
			},
			Correction: func(ev llmcontract.Correction) {
				slog.Warn("规则归一化输出自愈", "module", "rules", "source", source,
					"attempt", ev.Attempt, "layer", ev.Layer, "structured_mode", ev.Mode, "err", ev.Err)
			},
		},
	})
	if err != nil {
		return rules.Candidate{}, fmt.Errorf("归一化失败: %w", err)
	}
	return out.toCandidate(source)
}

func degraded(source, text string) rules.Candidate {
	return rules.Candidate{
		Source:      source,
		Preferences: text,
		Uncertain:   []string{source + "：归一化失败，已按原文作为风格偏好处理（未提炼机械规则）"},
		Degraded:    true,
	}
}

type normalizerOutput struct {
	Structured  normalizerStructured `json:"structured"`
	Preferences string               `json:"preferences"`
	Uncertain   []string             `json:"uncertain"`
}

type normalizerStructured struct {
	Genre            string             `json:"genre"`
	ForbiddenChars   []string           `json:"forbidden_chars"`
	ForbiddenPhrases []string           `json:"forbidden_phrases"`
	FatigueWords     []fatigueWordEntry `json:"fatigue_words"`
}

type fatigueWordEntry struct {
	Word          string `json:"word"`
	MaxPerChapter int    `json:"max_per_chapter"`
}

func (o normalizerOutput) toCandidate(source string) (rules.Candidate, error) {
	var fatigue map[string]int
	for _, e := range o.Structured.FatigueWords {
		word := strings.TrimSpace(e.Word)
		if word == "" {
			return rules.Candidate{}, fmt.Errorf("fatigue_words 含空词条目")
		}
		if e.MaxPerChapter < 1 {
			return rules.Candidate{}, fmt.Errorf("fatigue_words[%q].max_per_chapter 必须是正整数, got %d", word, e.MaxPerChapter)
		}
		if fatigue == nil {
			fatigue = make(map[string]int, len(o.Structured.FatigueWords))
		}
		fatigue[word] = e.MaxPerChapter
	}
	return rules.Candidate{
		Source: source,
		Structured: rules.Structured{
			Genre:            strings.TrimSpace(o.Structured.Genre),
			ForbiddenChars:   nonEmpty(o.Structured.ForbiddenChars),
			ForbiddenPhrases: nonEmpty(o.Structured.ForbiddenPhrases),
			FatigueWords:     fatigue,
		},
		Preferences: strings.TrimSpace(o.Preferences),
		Uncertain:   nonEmpty(o.Uncertain),
	}, nil
}

func nonEmpty(in []string) []string {
	var out []string
	for _, s := range in {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func activeNormalizerSystemPrompt() string {
	if localization.IsVietnamese() {
		return viNormalizerSystemPrompt
	}
	return zhNormalizerSystemPrompt
}

const zhNormalizerSystemPrompt = `你是 AI 小说写作系统的「规则归一化器」。你读取用户某一个来源的长期写作规则（自然语言），把明确且可机械检查的规则提升到 structured，其余内容归入 preferences 或 uncertain。

【保守提升——最重要】
- 只有用户明确、无歧义时才写入 structured。
- forbidden_chars/forbidden_phrases 是 error 级:只有「不要出现X/禁用X/别写X」这类明确禁止才提升。
- fatigue_words:只有同时给出「明确的词」和「明确的次数阈值」才提升;「少用X/别老用X」没给数字的放进 preferences,绝不自己发明阈值。
- 字数/篇幅类意愿(「每章3000字」「短一点」)一律放 preferences:章节长度是叙事节奏问题,由创作时自然把握,不做机械检查。
- 不可机械检查、无明确阈值、依赖语境的,一律放 preferences。
- 原则:宁可漏进 structured,也不要错误提升(那会每章误报)。

preferences 用一段可读的自然语言保留风格、人物与审美偏好。
uncertain 说明你故意没有提升到 structured 的项目及原因。`

const viNormalizerSystemPrompt = `Bạn là bộ chuẩn hóa quy tắc của một hệ thống viết tiểu thuyết bằng AI. Bạn đọc một nguồn quy tắc viết dài hạn do người dùng cung cấp bằng ngôn ngữ tự nhiên, nâng những yêu cầu thật sự rõ ràng và có thể kiểm tra máy móc vào structured; phần còn lại phải nằm trong preferences hoặc uncertain.

Nguyên tắc bảo thủ — quan trọng nhất:
- Chỉ đưa vào structured khi ý của người dùng rõ ràng và không mơ hồ.
- forbidden_chars/forbidden_phrases là ràng buộc mức lỗi: chỉ nâng những lệnh cấm trực tiếp kiểu "không được xuất hiện X", "cấm X", "đừng viết X".
- fatigue_words chỉ được dùng khi người dùng nêu cả một từ/cụm từ cụ thể và một ngưỡng số lần cụ thể. Các yêu cầu kiểu "ít dùng X" hoặc "đừng lạm dụng X" nhưng không có số phải để trong preferences; tuyệt đối không tự bịa threshold.
- Mọi mong muốn về độ dài như "mỗi chương khoảng 3000 từ" hoặc "viết ngắn hơn" phải nằm trong preferences. Độ dài chương là vấn đề nhịp kể, không biến thành kiểm tra cơ học ở đây.
- Yêu cầu phụ thuộc ngữ cảnh, không thể kiểm tra máy móc hoặc không có ngưỡng rõ ràng phải nằm trong preferences.
- Thà bỏ sót một mục khỏi structured còn hơn nâng sai và gây cảnh báo giả ở mọi chương.

preferences phải giữ lại bằng tiếng Việt tự nhiên các sở thích về phong cách, nhân vật và thẩm mỹ của người dùng.
uncertain phải giải thích ngắn gọn bằng tiếng Việt những mục bạn chủ động không nâng vào structured và lý do.`
