package bootstrap

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/voocel/ainovel-cli/internal/localization"
	"github.com/voocel/ainovel-cli/internal/rules"
	"github.com/voocel/ainovel-cli/internal/utils"
)

// exampleConfig is the annotated template written to ~/.ainovel/config.example.jsonc.
// The embedded file must stay in sync with the repository root config.example.jsonc.
//
//go:embed config.example.jsonc
var exampleConfig string

// NeedsSetup reports whether neither global nor project-level configuration exists.
func NeedsSetup() bool {
	if p := DefaultConfigPath(); p != "" {
		if _, err := os.Stat(p); err == nil {
			return false
		}
	}
	if _, err := os.Stat(projectConfigPath()); err == nil {
		return false
	}
	return true
}

type setupProvider struct {
	name           string
	label          string
	baseURL        string
	needType       bool
	apiKeyOptional bool
}

// ProviderPreset is a provider catalog entry shared by first-run setup and /config.
type ProviderPreset struct {
	Name           string
	Label          string
	BaseURL        string
	NeedType       bool
	APIKeyOptional bool
}

var setupProviders = []setupProvider{
	{name: "openrouter", label: "OpenRouter", baseURL: "https://openrouter.ai/api/v1"},
	{name: "anthropic", label: "Anthropic"},
	{name: "gemini", label: "Gemini"},
	{name: "openai", label: "OpenAI"},
	{name: "deepseek", label: "DeepSeek"},
	{name: "qwen", label: "Qwen"},
	{name: "glm", label: "GLM"},
	{name: "grok", label: "Grok"},
	{name: "ollama", label: "Ollama", baseURL: "http://localhost:11434/v1", apiKeyOptional: true},
	{name: "bedrock", label: "Bedrock", apiKeyOptional: true},
	{name: "custom", label: "Custom Proxy", needType: true, apiKeyOptional: true},
}

// ProviderPresets returns a safe-to-modify copy of the provider presets.
func ProviderPresets() []ProviderPreset {
	out := make([]ProviderPreset, 0, len(setupProviders))
	for _, preset := range setupProviders {
		out = append(out, ProviderPreset{
			Name: preset.name, Label: preset.label, BaseURL: preset.baseURL,
			NeedType: preset.needType, APIKeyOptional: preset.apiKeyOptional,
		})
	}
	return out
}

type setupCopy struct {
	noConfig           string
	configPath         string
	advancedSettings   string
	providerTitle      string
	providerName       string
	apiKeyOptional     string
	apiKeyOptionalHint string
	apiKeyRequired     string
	notSet             string
	baseURLTitle       string
	baseURLDefaultHint string
	defaultValue       string
	modelTitle         string
	modelExample       string
	savedTo            string
	defaultModel       string
	roleModelsHint     string
	globalRulesHint    string
	apiProtocolTitle   string
	compatibleSuffix   string
	selectHelp         string
	inputHelp          string
	cancelled          string
}

var viSetupCopy = setupCopy{
	noConfig:           "Không tìm thấy tệp cấu hình, bắt đầu thiết lập...",
	configPath:         "Đường dẫn cấu hình",
	advancedSettings:   "Sau khi hoàn tất, bạn có thể chỉnh sửa tệp này bất cứ lúc nào để thay đổi thiết lập nâng cao.",
	providerTitle:      "[1/4] Chọn Provider",
	providerName:       "Tên Provider",
	apiKeyOptional:     "[2/4] API Key (có thể để trống)",
	apiKeyOptionalHint: "Để trống nếu không dùng API Key",
	apiKeyRequired:     "[2/4] API Key",
	notSet:             "chưa đặt",
	baseURLTitle:       "[3/4] Base URL (Enter để dùng mặc định; nếu dùng proxy hãy nhập địa chỉ proxy)",
	baseURLDefaultHint: "Để trống để dùng địa chỉ chính thức",
	defaultValue:       "mặc định",
	modelTitle:         "[4/4] Tên model",
	modelExample:       "Ví dụ: gpt-4o / claude-sonnet-4 / gemini-2.5-pro",
	savedTo:            "Đã lưu cấu hình vào",
	defaultModel:       "Model mặc định",
	roleModelsHint:     "Muốn cấu hình model riêng cho từng vai trò, hãy chỉnh sửa tệp cấu hình.",
	globalRulesHint:    "Có thể đặt các tệp .md chứa tùy chọn viết toàn cục trong %s (xem README.txt trong thư mục đó).",
	apiProtocolTitle:   "Loại giao thức API",
	compatibleSuffix:   "tương thích",
	selectHelp:         "↑↓ chọn  Enter xác nhận  Esc hủy",
	inputHelp:          "Enter xác nhận, Esc hủy",
	cancelled:          "đã hủy thiết lập",
}

var zhSetupCopy = setupCopy{
	noConfig:           "未检测到配置文件，开始初始化设置...",
	configPath:         "配置文件路径",
	advancedSettings:   "完成后可随时编辑该文件调整高级设置。",
	providerTitle:      "[1/4] 选择 Provider",
	providerName:       "Provider 名称",
	apiKeyOptional:     "[2/4] API Key（可留空）",
	apiKeyOptionalHint: "留空表示不使用 API Key",
	apiKeyRequired:     "[2/4] API Key",
	notSet:             "未设置",
	baseURLTitle:       "[3/4] Base URL（直接回车使用默认，代理用户填写代理地址）",
	baseURLDefaultHint: "留空使用官方地址",
	defaultValue:       "默认",
	modelTitle:         "[4/4] 模型名称",
	modelExample:       "例如：gpt-4o / claude-sonnet-4 / gemini-2.5-pro",
	savedTo:            "配置已保存到",
	defaultModel:       "默认模型",
	roleModelsHint:     "如需按角色配置不同模型，编辑配置文件即可。",
	globalRulesHint:    "全局写作偏好可放 %s 下的 .md 文件（见其中 README.txt）",
	apiProtocolTitle:   "API 协议类型",
	compatibleSuffix:   "兼容",
	selectHelp:         "↑↓ 选择  Enter 确认  Esc 取消",
	inputHelp:          "Enter 确认, Esc 取消",
	cancelled:          "setup cancelled",
}

func activeSetupCopy() setupCopy {
	if localization.IsVietnamese() {
		return viSetupCopy
	}
	return zhSetupCopy
}

// RunSetup runs the first-run wizard and returns the generated configuration.
func RunSetup() (Config, error) {
	copy := activeSetupCopy()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Render(copy.noConfig))
	fmt.Fprintf(os.Stderr, "  %s: %s\n", copy.configPath, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(DefaultConfigPath()))
	fmt.Fprintf(os.Stderr, "  %s\n", copy.advancedSettings)
	fmt.Fprintln(os.Stderr)

	sp, err := runProviderSelect()
	if err != nil {
		return Config{}, err
	}

	providerName := sp.name
	var pc ProviderConfig
	printStepDone("Provider", sp.label)

	if sp.needType {
		providerName, err = runTextInput(copy.providerName, "my-proxy")
		if err != nil {
			return Config{}, err
		}
		providerType, err := runTypeSelect()
		if err != nil {
			return Config{}, err
		}
		pc.Type = providerType
	}

	var apiKey string
	if sp.apiKeyOptional {
		apiKey, err = runOptionalTextInput(copy.apiKeyOptional, copy.apiKeyOptionalHint)
	} else {
		apiKey, err = runTextInput(copy.apiKeyRequired, "sk-xxx")
	}
	if err != nil {
		return Config{}, err
	}
	pc.APIKey = apiKey
	if apiKey == "" {
		printStepDone("API Key", copy.notSet)
	} else {
		printStepDone("API Key", maskKey(apiKey))
	}

	baseDefault := sp.baseURL
	baseHint := copy.baseURLDefaultHint
	if baseDefault != "" {
		baseHint = baseDefault
	}
	baseURL, err := runTextInputWithDefault(copy.baseURLTitle, baseHint, baseDefault)
	if err != nil {
		return Config{}, err
	}
	pc.BaseURL = baseURL
	if baseURL != "" {
		printStepDone("Base URL", baseURL)
	} else {
		printStepDone("Base URL", copy.defaultValue)
	}

	modelName, err := runTextInput(copy.modelTitle, copy.modelExample)
	if err != nil {
		return Config{}, err
	}
	printStepDone("Model", modelName)
	pc.Models = []ModelConfig{{Name: modelName}}

	cfg := Config{
		Provider:  providerName,
		ModelName: modelName,
		Providers: map[string]ProviderConfig{providerName: pc},
		Roles:     map[string]RoleConfig{},
		Style:     "default",
	}

	path := DefaultConfigPath()
	if err := SaveConfig(path, cfg); err != nil {
		return cfg, fmt.Errorf("save config: %w", err)
	}

	saveExampleConfig()
	rulesDir := rules.DefaultHomeRulesDir()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%s %s %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓"), copy.savedTo, path)
	fmt.Fprintf(os.Stderr, "  %s: %s\n", copy.defaultModel, modelName)
	fmt.Fprintf(os.Stderr, "  %s\n", copy.roleModelsHint)
	if rulesDir != "" {
		fmt.Fprintf(os.Stderr, "  %s\n", fmt.Sprintf(copy.globalRulesHint, rulesDir))
	}
	fmt.Fprintln(os.Stderr)

	return cfg, nil
}

func saveExampleConfig() {
	dir, err := configDir()
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, "config.example.jsonc"), []byte(exampleConfig), 0o644)
}

func printStepDone(label, value string) {
	fmt.Fprintf(os.Stderr, "  %s %s: %s\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓"),
		label,
		lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(value))
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func runProviderSelect() (setupProvider, error) {
	m := setupSelectModel{title: activeSetupCopy().providerTitle, items: setupProviders}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return setupProvider{}, err
	}
	result := final.(setupSelectModel)
	if result.cancelled {
		return setupProvider{}, fmt.Errorf("%s", activeSetupCopy().cancelled)
	}
	return result.items[result.cursor], nil
}

func apiTypeOptions() []setupProvider {
	suffix := activeSetupCopy().compatibleSuffix
	return []setupProvider{
		{name: "openai", label: "OpenAI " + suffix},
		{name: "anthropic", label: "Anthropic " + suffix},
		{name: "gemini", label: "Gemini " + suffix},
	}
}

func runTypeSelect() (string, error) {
	m := setupSelectModel{title: activeSetupCopy().apiProtocolTitle, items: apiTypeOptions()}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	result := final.(setupSelectModel)
	if result.cancelled {
		return "", fmt.Errorf("%s", activeSetupCopy().cancelled)
	}
	return result.items[result.cursor].name, nil
}

func runTextInput(label, placeholder string) (string, error) {
	return runTextInputWithDefault(label, placeholder, "")
}

func runOptionalTextInput(label, placeholder string) (string, error) {
	m := setupInputModel{label: label, placeholder: placeholder, allowEmpty: true}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	result := final.(setupInputModel)
	if result.cancelled {
		return "", fmt.Errorf("%s", activeSetupCopy().cancelled)
	}
	return utils.CleanInputLine(result.value), nil
}

func runTextInputWithDefault(label, placeholder, defaultValue string) (string, error) {
	m := setupInputModel{label: label, placeholder: placeholder, defaultValue: defaultValue}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	result := final.(setupInputModel)
	if result.cancelled {
		return "", fmt.Errorf("%s", activeSetupCopy().cancelled)
	}
	if result.value == "" && result.defaultValue != "" {
		return result.defaultValue, nil
	}
	return utils.CleanInputLine(result.value), nil
}

var (
	setupCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	setupDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	setupHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	setupInputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

type setupSelectModel struct {
	title     string
	items     []setupProvider
	cursor    int
	cancelled bool
}

func (m setupSelectModel) Init() tea.Cmd { return nil }

func (m setupSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m setupSelectModel) View() string {
	var b strings.Builder
	b.WriteString(setupHeaderStyle.Render(m.title))
	b.WriteString("\n\n")
	for i, item := range m.items {
		cursor := "  "
		label := item.label
		if i == m.cursor {
			cursor = setupCursorStyle.Render("❯ ")
			label = setupCursorStyle.Render(label)
		}
		b.WriteString(cursor + label + "\n")
	}
	b.WriteString(setupDimStyle.Render("\n  " + activeSetupCopy().selectHelp))
	return b.String()
}

type setupInputModel struct {
	label        string
	placeholder  string
	defaultValue string
	allowEmpty   bool
	value        string
	cancelled    bool
}

func (m setupInputModel) Init() tea.Cmd { return nil }

func (m setupInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter":
			if utils.CleanInputLine(m.value) != "" || m.defaultValue != "" || m.allowEmpty {
				return m, tea.Quit
			}
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "backspace":
			if len(m.value) > 0 {
				runes := []rune(m.value)
				m.value = string(runes[:len(runes)-1])
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.value += utils.CleanInputRunes(msg.Runes)
			} else if msg.Type == tea.KeySpace {
				m.value += " "
			}
		}
	}
	return m, nil
}

func (m setupInputModel) View() string {
	var b strings.Builder
	b.WriteString(setupHeaderStyle.Render(m.label))
	b.WriteString("\n\n")
	b.WriteString(setupInputStyle.Render("❯ "))
	if m.value == "" {
		b.WriteString(setupCursorStyle.Render("▌"))
		b.WriteString(setupDimStyle.Render(m.placeholder))
	} else {
		b.WriteString(m.value)
		b.WriteString(setupCursorStyle.Render("▌"))
	}
	b.WriteString(setupDimStyle.Render("  (" + activeSetupCopy().inputHelp + ")"))
	b.WriteString("\n")
	return b.String()
}
