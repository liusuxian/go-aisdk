/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 22:36:41
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/internal/utils"
)

// ChatMessage 聊天消息的通用接口
type ChatMessage interface {
	SetProvider(provider string)        // 设置提供商（序列化参数时，处理差异化参数）
	MarshalJSON() (b []byte, err error) // 序列化JSON
}

// ChatAudioFormatType 输出音频的格式
type ChatAudioFormatType string

const (
	// 提供商支持: OpenAI
	ChatAudioFormatTypeWAV ChatAudioFormatType = "wav"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeMP3 ChatAudioFormatType = "mp3"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeFLAC ChatAudioFormatType = "flac"
	// 提供商支持: OpenAI
	ChatAudioFormatTypeOPUS ChatAudioFormatType = "opus"
	// 提供商支持: OpenAI
	ChatAudioFormatTypePCM16 ChatAudioFormatType = "pcm16"
)

// ChatAudioVoiceType 输出音频的音色
type ChatAudioVoiceType string

const (
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeAlloy ChatAudioVoiceType = "alloy"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeAsh ChatAudioVoiceType = "ash"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeBallad ChatAudioVoiceType = "ballad"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeCoral ChatAudioVoiceType = "coral"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeEcho ChatAudioVoiceType = "echo"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeFable ChatAudioVoiceType = "fable"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeNova ChatAudioVoiceType = "nova"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeOnyx ChatAudioVoiceType = "onyx"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeSage ChatAudioVoiceType = "sage"
	// 提供商支持: OpenAI
	ChatAudioVoiceTypeShimmer ChatAudioVoiceType = "shimmer"
)

// ChatAudioOutputArgs 音频输出参数
type ChatAudioOutputArgs struct {
	// 输出音频的格式
	//
	// 提供商支持: OpenAI
	Format ChatAudioFormatType `json:"format,omitempty" providers:"openai"`
	// 输出音频的音色
	//
	// 提供商支持: OpenAI
	Voice ChatAudioVoiceType `json:"voice,omitempty" providers:"openai"`
}

// ChatModalitiesType 输出数据的模态
type ChatModalitiesType string

const (
	// 提供商支持: OpenAI
	ChatModalitiesTypeText ChatModalitiesType = "text"
	// 提供商支持: OpenAI
	ChatModalitiesTypeAudio ChatModalitiesType = "audio"
)

// ChatPredictionType 预测内容的类型
type ChatPredictionType string

const (
	// 提供商支持: OpenAI
	ChatPredictionTypeContent ChatPredictionType = "content"
)

// ChatPredictionContentPart 预测内容
type ChatPredictionContentPart struct {
	// 内容的类型
	//
	// 提供商支持: OpenAI
	Type string `json:"type,omitempty" providers:"openai"`
	// 文本内容
	//
	// 提供商支持: OpenAI
	Text string `json:"text,omitempty" providers:"openai"`
}

// ChatPrediction 预测输出配置
type ChatPrediction struct {
	// 预测内容的类型
	//
	// 提供商支持: OpenAI
	Type ChatPredictionType `json:"type,omitempty" providers:"openai"`
	// 预测内容
	//
	// 提供商支持: OpenAI
	Content []ChatPredictionContentPart `json:"content,omitempty" providers:"openai"`
}

// ChatReasoningEffortType 推理努力程度
type ChatReasoningEffortType string

const (
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeLow ChatReasoningEffortType = "low"
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeMedium ChatReasoningEffortType = "medium"
	// 提供商支持: OpenAI
	ChatReasoningEffortTypeHigh ChatReasoningEffortType = "high"
)

// ChatResponseFormatType 响应格式的类型
type ChatResponseFormatType string

const (
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatResponseFormatTypeText ChatResponseFormatType = "text"
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatResponseFormatTypeJSONObject ChatResponseFormatType = "json_object"
	// 提供商支持: OpenAI
	ChatResponseFormatTypeJSONSchema ChatResponseFormatType = "json_schema"
)

// ChatResponseFormatJSONSchema JSON Schema 配置
type ChatResponseFormatJSONSchema struct {
	// 响应格式名称，必须是 a-z、A-Z、0-9 或包含下划线和破折号，最大长度为 64
	//
	// 提供商支持: OpenAI
	Name string `json:"name,omitempty" providers:"openai"`
	// 响应格式的描述，用于指导模型如何响应
	//
	// 提供商支持: OpenAI
	Description string `json:"description,omitempty" providers:"openai"`
	// 响应格式的 JSON Schema
	//
	// 提供商支持: OpenAI
	Schema map[string]any `json:"schema,omitempty" providers:"openai"`
	// 是否启用严格模式，默认为 false
	//
	// 提供商支持: OpenAI
	Strict *bool `json:"strict,omitempty" providers:"openai"`
}

// ChatResponseFormat 响应格式
type ChatResponseFormat struct {
	// 响应格式的类型
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Type ChatResponseFormatType `json:"type,omitempty" providers:"openai,deepseek,alibl"`
	// JSON Schema 配置，仅当 Type 为 "json_schema" 时使用
	//
	// 提供商支持: OpenAI
	JSONSchema *ChatResponseFormatJSONSchema `json:"json_schema,omitempty" providers:"openai"`
}

// ChatStreamOptions 流式传输选项
type ChatStreamOptions struct {
	// 是否包含令牌使用统计信息
	//
	// 提供商支持: OpenAI | DeepSeek
	IncludeUsage *bool `json:"include_usage,omitempty" providers:"openai,deepseek"`
}

// ChatToolChoiceType 工具调用类型
type ChatToolChoiceType string

const (
	// 模型不会调用任何 tool，而是生成一条消息
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatToolChoiceTypeNone ChatToolChoiceType = "none"
	// 模型可以选择生成一条消息或调用一个或多个 tool
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ChatToolChoiceTypeAuto ChatToolChoiceType = "auto"
	// 模型必须调用一个或多个 tool
	//
	// 提供商支持: OpenAI | DeepSeek
	ChatToolChoiceTypeRequired ChatToolChoiceType = "required"
)

// ChatToolChoiceFunction 工具调用函数
type ChatToolChoiceFunction struct {
	// 工具调用函数名称
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Name string `json:"name,omitempty" providers:"openai,deepseek,alibl"`
}

// ToolType 工具类型
type ToolType string

const (
	// 函数
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolTypeFunction ToolType = "function"
)

// ChatToolChoice 指定工具调用的策略
type ChatToolChoice struct {
	// 工具调用类型
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolChoiceType ChatToolChoiceType
	// 工具调用函数
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Function *ChatToolChoiceFunction
	// 工具类型
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Type ToolType
}

// MarshalJSON 序列化JSON
func (c ChatToolChoice) MarshalJSON() (b []byte, err error) {
	if c.ToolChoiceType != "" {
		return json.Marshal(c.ToolChoiceType)
	}
	return json.Marshal(struct {
		Function *ChatToolChoiceFunction `json:"function,omitempty"`
		Type     ToolType                `json:"type,omitempty"`
	}{
		Function: c.Function,
		Type:     c.Type,
	})
}

// ChatToolFunction 工具函数
type ChatToolFunction struct {
	// 函数名称，必须是 a-z, A-Z, 0-9 或者包含下划线和破折号，最大长度为 64
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Name string `json:"name,omitempty" providers:"openai,deepseek,alibl"`
	// 函数描述，用于帮助模型决定何时以及如何调用函数
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Description string `json:"description,omitempty" providers:"openai,deepseek,alibl"`
	// 函数接受的参数，描述为一个 JSON Schema 对象
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Parameters map[string]any `json:"parameters,omitempty" providers:"openai,deepseek,alibl"`
	// 是否启用严格模式，默认为 false
	//
	// 提供商支持: OpenAI
	Strict *bool `json:"strict,omitempty" providers:"openai"`
}

// ChatTool 工具
type ChatTool struct {
	// 工具类型
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Type ToolType `json:"type,omitempty" providers:"openai,deepseek,alibl"`
	// 工具函数
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Function *ChatToolFunction `json:"function,omitempty" providers:"openai,deepseek,alibl"`
}

// ChatSearchContextSize 搜索上下文大小
type ChatSearchContextSize string

const (
	// 提供商支持: OpenAI
	ChatSearchContextSizeLow ChatSearchContextSize = "low"
	// 提供商支持: OpenAI
	ChatSearchContextSizeMedium ChatSearchContextSize = "medium"
	// 提供商支持: OpenAI
	ChatSearchContextSizeHigh ChatSearchContextSize = "high"
)

// ChatApproximateLocation 用户的大致位置参数
type ChatApproximateLocation struct {
	// 用户所在城市
	//
	// 提供商支持: OpenAI
	City string `json:"city,omitempty" providers:"openai"`
	// 用户所在国家的两字母 ISO 代码
	//
	// 提供商支持: OpenAI
	Country string `json:"country,omitempty" providers:"openai"`
	// 用户所在地区
	//
	// 提供商支持: OpenAI
	Region string `json:"region,omitempty" providers:"openai"`
	// 用户的 IANA 时区
	//
	// 提供商支持: OpenAI
	Timezone string `json:"timezone,omitempty" providers:"openai"`
}

// ChatApproximateLocationType 位置近似类型
type ChatApproximateLocationType string

const (
	// 提供商支持: OpenAI
	ChatApproximateLocationTypeApproximate ChatApproximateLocationType = "approximate"
)

// ChatUserLocation 用户位置信息
type ChatUserLocation struct {
	// 大致位置信息
	//
	// 提供商支持: OpenAI
	Approximate *ChatApproximateLocation `json:"approximate,omitempty" providers:"openai"`
	// 位置近似类型
	//
	// 提供商支持: OpenAI
	Type ChatApproximateLocationType `json:"type,omitempty" providers:"openai"`
}

// ChatCitationFormat 角标样式
type ChatCitationFormat string

const (
	// 角标形式为[1]
	//
	// 提供商支持: AliBL
	ChatCitationFormatNumber ChatCitationFormat = "[<number>]"
	// 角标形式为[ref_1]
	//
	// 提供商支持: AliBL
	ChatCitationFormatRefNumber ChatCitationFormat = "[ref_<number>]"
)

// ChatSearchStrategy 搜索互联网信息的数量
type ChatSearchStrategy string

const (
	// 在请求时搜索5条互联网信息
	//
	// 提供商支持: AliBL
	ChatSearchStrategyStandard ChatSearchStrategy = "standard"
	// 在请求时搜索10条互联网信息
	//
	// 提供商支持: AliBL
	ChatSearchStrategyPro ChatSearchStrategy = "pro"
)

// ChatWebSearchOptions 网络搜索选项
type ChatWebSearchOptions struct {
	// 搜索上下文大小
	//
	// 提供商支持: OpenAI
	SearchContextSize ChatSearchContextSize `json:"search_context_size,omitempty" providers:"openai"`
	// 用户位置信息
	//
	// 提供商支持: OpenAI
	UserLocation *ChatUserLocation `json:"user_location,omitempty" providers:"openai"`
	// 在返回结果中是否展示搜索到的信息
	//
	// 提供商支持: AliBL
	EnableSource *bool `json:"enable_source,omitempty" providers:"alibl"`
	// 是否开启[1]或[ref_1]样式的角标标注功能
	//
	// 提供商支持: AliBL
	EnableCitation *bool `json:"enable_citation,omitempty" providers:"alibl"`
	// 角标样式。在enable_citation为true时生效
	//
	// 提供商支持: AliBL
	CitationFormat ChatCitationFormat `json:"citation_format,omitempty" providers:"alibl"`
	// 是否强制开启搜索
	//
	// 提供商支持: AliBL
	ForcedSearch *bool `json:"forced_search,omitempty" providers:"alibl"`
	// 搜索互联网信息的数量
	//
	// 提供商支持: AliBL
	SearchStrategy ChatSearchStrategy `json:"search_strategy,omitempty" providers:"alibl"`
}

// ChatOCRTask OCR模型执行内置任务的名称
type ChatOCRTask string

const (
	// 通用文字识别
	//
	// 提供商支持: AliBL
	ChatOCRTaskTextRecognition ChatOCRTask = "text_recognition"
	// 信息抽取
	//
	// 提供商支持: AliBL
	ChatOCRTaskKeyInformationExtraction ChatOCRTask = "key_information_extraction"
	// 文档解析
	//
	// 提供商支持: AliBL
	ChatOCRTaskDocumentParsing ChatOCRTask = "document_parsing"
	// 表格解析
	//
	// 提供商支持: AliBL
	ChatOCRTaskTableParsing ChatOCRTask = "table_parsing"
	// 公式识别
	//
	// 提供商支持: AliBL
	ChatOCRTaskFormulaRecognition ChatOCRTask = "formula_recognition"
	// 多语言识别
	//
	// 提供商支持: AliBL
	ChatOCRTaskMultiLan ChatOCRTask = "multi_lan"
)

// ChatOCRTaskConfig OCR模型执行内置任务的配置项
type ChatOCRTaskConfig struct {
	// 需要模型抽取的字段，可以是任意形式的JSON结构，最多可嵌套3层JSON 对象。您只需要填写JSON对象的key，value保持为空即可
	//
	// 提供商支持: AliBL
	ResultSchema map[string]any `json:"result_schema,omitempty" providers:"alibl"`
}

// ChatOCROptions OCR模型执行内置任务时需要配置的参数
type ChatOCROptions struct {
	// 内置任务的名称
	//
	// 提供商支持: AliBL
	Task ChatOCRTask `json:"task,omitempty" providers:"alibl"`
	// 当内置任务task为"key_information_extraction"（信息抽取）时使用
	//
	// 提供商支持: AliBL
	TaskConfig []ChatOCRTaskConfig `json:"task_config,omitempty" providers:"alibl"`
}

// ChatTranslationLanguageType 翻译支持的语言类型
type ChatTranslationLanguageType string

const (
	// 自动检测语言
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeAuto ChatTranslationLanguageType = "auto"
	// 中文
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeChinese ChatTranslationLanguageType = "Chinese"
	// 英语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeEnglish ChatTranslationLanguageType = "English"
	// 日语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeJapanese ChatTranslationLanguageType = "Japanese"
	// 韩语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeKorean ChatTranslationLanguageType = "Korean"
	// 泰语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeThai ChatTranslationLanguageType = "Thai"
	// 法语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFrench ChatTranslationLanguageType = "French"
	// 德语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeGerman ChatTranslationLanguageType = "German"
	// 西班牙语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeSpanish ChatTranslationLanguageType = "Spanish"
	// 阿拉伯语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeArabic ChatTranslationLanguageType = "Arabic"
	// 印尼语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeIndonesian ChatTranslationLanguageType = "Indonesian"
	// 越南语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeVietnamese ChatTranslationLanguageType = "Vietnamese"
	// 巴西葡萄牙语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePortuguese ChatTranslationLanguageType = "Portuguese"
	// 意大利语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeItalian ChatTranslationLanguageType = "Italian"
	// 荷兰语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeDutch ChatTranslationLanguageType = "Dutch"
	// 俄语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeRussian ChatTranslationLanguageType = "Russian"
	// 高棉语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeKhmer ChatTranslationLanguageType = "Khmer"
	// 宿务语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCebuano ChatTranslationLanguageType = "Cebuano"
	// 菲律宾语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFilipino ChatTranslationLanguageType = "Filipino"
	// 捷克语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCzech ChatTranslationLanguageType = "Czech"
	// 波兰语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePolish ChatTranslationLanguageType = "Polish"
	// 波斯语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypePersian ChatTranslationLanguageType = "Persian"
	// 希伯来语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHebrew ChatTranslationLanguageType = "Hebrew"
	// 土耳其语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeTurkish ChatTranslationLanguageType = "Turkish"
	// 印地语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHindi ChatTranslationLanguageType = "Hindi"
	// 孟加拉语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeBengali ChatTranslationLanguageType = "Bengali"
	// 乌尔都语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeUrdu ChatTranslationLanguageType = "Urdu"
	// 马来语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeMalay ChatTranslationLanguageType = "Malay"
	// 匈牙利语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeHungarian ChatTranslationLanguageType = "Hungarian"
	// 瑞典语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeSwedish ChatTranslationLanguageType = "Swedish"
	// 芬兰语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeFinnish ChatTranslationLanguageType = "Finnish"
	// 缅甸语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeBurmese ChatTranslationLanguageType = "Burmese"
	// 老挝语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeLao ChatTranslationLanguageType = "Lao"
	// 粤语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeCantonese ChatTranslationLanguageType = "Cantonese"
	// 罗马尼亚语
	//
	// 提供商支持: AliBL
	ChatTranslationLanguageTypeRomanian ChatTranslationLanguageType = "Romanian"
)

// ChatTranslationTerm 在使用术语干预翻译功能时需要设置的术语
type ChatTranslationTerm struct {
	// 源语言的术语
	//
	// 提供商支持: AliBL
	Source string `json:"source,omitempty" providers:"alibl"`
	// 目标语言的术语
	//
	// 提供商支持: AliBL
	Target string `json:"target,omitempty" providers:"alibl"`
}

// ChatTranslationMemory 在使用翻译记忆功能时需要设置的翻译记忆
type ChatTranslationMemory struct {
	// 源语言的语句
	//
	// 提供商支持: AliBL
	Source string `json:"source,omitempty" providers:"alibl"`
	// 目标语言的语句
	//
	// 提供商支持: AliBL
	Target string `json:"target,omitempty" providers:"alibl"`
}

// ChatTranslationOptions 翻译选项
type ChatTranslationOptions struct {
	// 源语言的英文全称，可以将source_lang设置为"auto"，模型会自动判断输入文本属于哪种语言
	//
	// 提供商支持: AliBL
	SourceLang ChatTranslationLanguageType `json:"source_lang,omitempty" providers:"alibl"`
	// 目标语言的英文全称
	//
	// 提供商支持: AliBL
	TargetLang ChatTranslationLanguageType `json:"target_lang,omitempty" providers:"alibl"`
	// 在使用术语干预翻译功能时需要设置的术语数组
	//
	// 提供商支持: AliBL
	Terms []ChatTranslationTerm `json:"terms,omitempty" providers:"alibl"`
	// 在使用翻译记忆功能时需要设置的翻译记忆数组
	//
	// 提供商支持: AliBL
	TmList []ChatTranslationMemory `json:"tm_list,omitempty" providers:"alibl"`
	// 在使用领域提示功能时需要设置的领域提示语句
	//
	// 提供商支持: AliBL
	Domains string `json:"domains,omitempty" providers:"alibl"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 消息数组
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Messages []ChatMessage `json:"messages,omitempty" providers:"openai,deepseek,alibl" group:"alibl:input"`
	// 模型名称
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Model string `json:"model,omitempty" providers:"openai,deepseek,alibl"`
	// 输出音频的音色与格式
	//
	// 提供商支持: OpenAI
	Audio *ChatAudioOutputArgs `json:"audio,omitempty" providers:"openai"`
	// 介于 -2.0 和 2.0 之间的数字（AliBL取值范围：只要大于0即可，1.0表示不做惩罚）。如果该值为正，那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚，降低模型重复相同内容的可能性
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty" providers:"openai,deepseek,alibl" mapping:"alibl:repetition_penalty" group:"alibl:parameters"`
	// 修改指定标记在补全中出现的可能性
	//
	// 提供商支持: OpenAI
	LogitBias map[string]int `json:"logit_bias,omitempty" providers:"openai"`
	// 是否返回输出 Token 的对数概率
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	LogProbs *bool `json:"logprobs,omitempty" providers:"openai,deepseek,alibl"`
	// 生成补全内容的最大令牌数上限
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	MaxCompletionTokens *int `json:"max_completion_tokens,omitempty" providers:"openai,deepseek,alibl" mapping:"deepseek|alibl:max_tokens" group:"alibl:parameters"`
	// 元数据（AliBL支持该参数，但不会被序列化，会放到请求头中）
	//
	// AliBL：map[string]string{"X-DashScope-DataInspection": "{\"input\": \"cip\", \"output\": \"cip\"}"} -> 在 API 的内容安全能力基础上，进一步识别输入输出内容的违规信息
	//
	// 提供商支持: OpenAI | AliBL
	Metadata map[string]string `json:"metadata,omitempty" providers:"openai"`
	// 输出数据的模态
	//
	// 提供商支持: OpenAI
	Modalities []ChatModalitiesType `json:"modalities,omitempty" providers:"openai"`
	// 生成响应的个数
	//
	// 提供商支持: OpenAI
	N *int `json:"n,omitempty" providers:"openai"`
	// 是否开启并行工具调用
	//
	// 提供商支持: OpenAI | AliBL
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty" providers:"openai,alibl"`
	// 预测输出配置
	//
	// 提供商支持: OpenAI
	Prediction *ChatPrediction `json:"prediction,omitempty" providers:"openai"`
	// 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	PresencePenalty *float32 `json:"presence_penalty,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 仅适用于 o 系列模型，约束推理模型的推理努力程度
	//
	// 提供商支持: OpenAI
	ReasoningEffort ChatReasoningEffortType `json:"reasoning_effort,omitempty" providers:"openai"`
	// 响应格式
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ResponseFormat *ChatResponseFormat `json:"response_format,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 随机种子
	//
	// 提供商支持: OpenAI | AliBL
	Seed *int `json:"seed,omitempty" providers:"openai,alibl" group:"alibl:parameters"`
	// 指定用于处理请求的延迟层级。此参数与订阅了规模层级服务的客户相关
	//
	// 提供商支持: OpenAI
	ServiceTier string `json:"service_tier,omitempty" providers:"openai"`
	// 当API遇到这些序列时将停止生成更多标记。返回的文本不会包含停止序列
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Stop []string `json:"stop,omitempty" providers:"openai,deepseek,alibl"`
	// 是否存储此聊天完成请求的输出，用于我们的模型蒸馏或评估产品
	//
	// 提供商支持: OpenAI
	Store *bool `json:"store,omitempty" providers:"openai"`
	// 是否流式传输响应（AliBL支持该参数，但不会被序列化，会放到请求头中）
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Stream *bool `json:"stream,omitempty" providers:"openai,deepseek"`
	// 流式传输选项
	//
	// 提供商支持: OpenAI | DeepSeek
	StreamOptions *ChatStreamOptions `json:"stream_options,omitempty" providers:"openai,deepseek"`
	// 采样温度，介于 0 和 2 之间（AliBL取值范围：[0,2)）。更高的值，会使输出更随机，而更低的值，会使其更加集中和确定
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Temperature *float32 `json:"temperature,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 指定工具调用的策略
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolChoice *ChatToolChoice `json:"tool_choice,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 可供模型调用的工具数组（AliBL使用 tools 时需要同时指定result_format参数为"message"。无论是发起 Function Calling，还是向模型提交工具函数的执行结果，均需设置tools参数）
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Tools []ChatTool `json:"tools,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 一个介于0和20之间的整数（AliBL取值范围：[0,5]），指定在每个标记位置返回的最可能标记的数量，每个标记都有相关的对数概率。如果使用此参数，必须将logprobs设置为true
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	TopLogProbs *int `json:"top_logprobs,omitempty" providers:"openai,deepseek,alibl"`
	// 核采样的概率阈值，介于 0 和 1 之间（AliBL取值范围：（0,1.0]）。较高的值，会使输出更随机，而较低的值，会使其更加集中和确定
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	TopP *float32 `json:"top_p,omitempty" providers:"openai,deepseek,alibl" group:"alibl:parameters"`
	// 网络搜索选项
	//
	// 提供商支持: OpenAI | AliBL
	WebSearchOptions *ChatWebSearchOptions `json:"web_search_options,omitempty" providers:"openai,alibl" mapping:"alibl:search_options" group:"alibl:parameters"`
	// 生成过程中采样候选集的大小。取值越大，生成的随机性越高，取值越小，生成的确定性越高。不赋值或当top_k大于100时，表示不启用top_k策略，此时仅有top_p策略生效。取值需要大于或等于0
	//
	// 提供商支持: AliBL
	TopK *int `json:"top_k,omitempty" providers:"alibl" group:"alibl:parameters"`
	// 是否开启思考模式
	//
	// 提供商支持: AliBL
	EnableThinking *bool `json:"enable_thinking,omitempty" providers:"alibl" group:"alibl:parameters"`
	// 思考过程的最大长度，在enable_thinking为true时生效
	//
	// 提供商支持: AliBL
	ThinkingBudget *int `json:"thinking_budget,omitempty" providers:"alibl"`
	// OCR模型执行内置任务时需要配置的参数
	//
	// 提供商支持: AliBL
	OcrOptions *ChatOCROptions `json:"ocr_options,omitempty" providers:"alibl" group:"alibl:parameters"`
	// 翻译选项
	//
	// 提供商支持: AliBL
	TranslationOptions *ChatTranslationOptions `json:"translation_options,omitempty" providers:"alibl" group:"alibl:parameters"`
}

// MarshalJSON 序列化JSON
func (r ChatRequest) MarshalJSON() (b []byte, err error) {
	// 设置提供商
	provider := r.Provider.String()
	for _, message := range r.Messages {
		message.SetProvider(provider)
	}
	// 处理提供商差异化内容
	switch r.Provider {
	case consts.AliBL:
		// 创建一个别名结构体
		type Alias ChatRequest
		temp := struct {
			Alias
			// 模型在生成文本时是否使用互联网搜索结果进行参考
			EnableSearch *bool `json:"enable_search,omitempty" providers:"alibl" group:"alibl:parameters"`
			// 是否提高输入图片的默认Token上限。输入图片的默认Token上限为1280，配置为true时输入图片的Token上限为16384
			VlHighResolutionImages *bool `json:"vl_high_resolution_images,omitempty" providers:"alibl" group:"alibl:parameters"`
			// 在流式输出模式下是否开启增量输出
			IncrementalOutput *bool `json:"incremental_output,omitempty" providers:"alibl" group:"alibl:parameters"`
			// 返回数据的格式
			ResultFormat string `json:"result_format,omitempty" providers:"alibl" group:"alibl:parameters"`
		}{
			Alias:                  Alias(r),
			VlHighResolutionImages: Bool(true),
			ResultFormat:           "message",
		}
		if r.WebSearchOptions != nil {
			temp.EnableSearch = Bool(true)
		}
		if BoolValue(r.Stream) {
			temp.IncrementalOutput = Bool(true)
		}
		// TODO 特殊处理（qwen-plus系列模型在开启思考模式时，incremental_output必须为true）
		if strings.Contains(r.Model, "qwen-plus") && BoolValue(r.EnableThinking) {
			temp.IncrementalOutput = Bool(true)
		}
		// 序列化JSON
		temp.Provider = ""
		temp.Metadata = nil
		temp.Stream = nil
		return utils.NewSerializer(provider).Serialize(temp)
	default:
		// 序列化JSON
		r.Provider = ""
		return utils.NewSerializer(provider).Serialize(r)
	}
}

// ChatFinishReason 模型停止生成 token 的原因
type ChatFinishReason string

const (
	ChatFinishReasonNull                       ChatFinishReason = "null"
	ChatFinishReasonStop                       ChatFinishReason = "stop"                         // 模型自然停止生成，或遇到 stop 序列中列出的字符串
	ChatFinishReasonLength                     ChatFinishReason = "length"                       // 输出长度达到了模型上下文长度限制，或达到了 max_completion_tokens 的限制
	ChatFinishReasonContentFilter              ChatFinishReason = "content_filter"               // 输出内容因触发过滤策略而被过滤
	ChatFinishReasonToolCalls                  ChatFinishReason = "tool_calls"                   // 模型调用了工具
	ChatFinishReasonInsufficientSystemResource ChatFinishReason = "insufficient_system_resource" // 系统推理资源不足，生成被打断
)

// MarshalJSON 序列化JSON
func (r ChatFinishReason) MarshalJSON() (b []byte, err error) {
	if r == ChatFinishReasonNull || r == "" {
		return []byte("null"), nil
	}
	return []byte(`"` + string(r) + `"`), nil
}

// ChatTopLogProbs 输出概率 top N 的 token 的列表，以及它们的对数概率
type ChatTopLogProbs struct {
	Bytes   []byte  `json:"bytes,omitempty"`   // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
	LogProb float64 `json:"logprob,omitempty"` // 该 token 的对数概率
	Token   string  `json:"token,omitempty"`   // 输出的 token
}

// ChatLogProb 对数概率信息
type ChatLogProb struct {
	Bytes       []byte            `json:"bytes,omitempty"`        // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
	LogProb     float64           `json:"logprob,omitempty"`      // 该 token 的对数概率
	Token       string            `json:"token,omitempty"`        // 输出的 token
	TopLogProbs []ChatTopLogProbs `json:"top_logprobs,omitempty"` // 一个包含在该输出位置上，输出概率 top N 的 token 的列表，以及它们的对数概率
}

// ChatLogProbs 对数概率信息
type ChatLogProbs struct {
	Content []ChatLogProb `json:"content,omitempty"` // 一个包含输出 token 对数概率信息的列表
	Refusal []ChatLogProb `json:"refusal,omitempty"` // 一个带有对数概率信息的消息拒绝 token 列表
}

// ChatAnnotationType 消息的注释类型
type ChatAnnotationType string

const (
	ChatAnnotationTypeURLCitation ChatAnnotationType = "url_citation" // 网络搜索工具的 URL 引用
)

// ChatAnnotationURLCitation 网络搜索工具的 URL 引用
type ChatAnnotationURLCitation struct {
	EndIndex   int    `json:"end_index,omitempty"`   // URL引用在消息中的结束位置索引
	StartIndex int    `json:"start_index,omitempty"` // URL引用在消息中的起始位置索引
	SiteName   string `json:"site_name,omitempty"`   // 搜索结果来源的网站名称
	Icon       string `json:"icon,omitempty"`        // 来源网站的图标URL，如果没有图标则为空字符串
	Title      string `json:"title,omitempty"`       // 被引用网页的标题
	URL        string `json:"url,omitempty"`         // 被引用网页的URL地址
}

// ChatAnnotation 消息的注释
type ChatAnnotation struct {
	Index       int                        `json:"index,omitempty"`        // 搜索结果的序号，表示该搜索结果在 []ChatAnnotation 中的索引
	Type        ChatAnnotationType         `json:"type,omitempty"`         // 注释的类型
	URLCitation *ChatAnnotationURLCitation `json:"url_citation,omitempty"` // 网络搜索工具的 URL 引用
}

// ChatAudioOutput 音频响应数据
type ChatAudioOutput struct {
	Data       string `json:"data,omitempty"`       // Base64 编码的音频字节，格式在请求中指定
	ExpiresAt  int64  `json:"expires_at,omitempty"` // 音频响应在服务器上不再可访问的 Unix 时间戳（秒）
	ID         string `json:"id,omitempty"`         // 音频响应的唯一标识符
	Transcript string `json:"transcript,omitempty"` // 模型生成的音频文本转录
}

// ToolCallsFunction 函数
type ToolCallsFunction struct {
	Arguments string `json:"arguments,omitempty"` // 函数参数
	Name      string `json:"name,omitempty"`      // 函数名
}

// ToolCalls 工具调用
type ToolCalls struct {
	Index    int                `json:"index,omitempty"`    // 索引
	Function *ToolCallsFunction `json:"function,omitempty"` // 函数调用
	ID       string             `json:"id,omitempty"`       // 工具ID
	Type     ToolType           `json:"type,omitempty"`     // 工具类型
}

// ChatCompletionMessage 模型生成的 completion 消息
type ChatCompletionMessage struct {
	Content          string           `json:"content,omitempty"`           // 文本内容
	ReasoningContent string           `json:"reasoning_content,omitempty"` // 推理内容
	Refusal          string           `json:"refusal,omitempty"`           // 拒绝消息
	Role             string           `json:"role,omitempty"`              // 角色
	Annotations      []ChatAnnotation `json:"annotations,omitempty"`       // 消息的注释，在适用情况下提供，例如使用网络搜索工具时
	Audio            *ChatAudioOutput `json:"audio,omitempty"`             // 音频响应数据
	ToolCalls        []ToolCalls      `json:"tool_calls,omitempty"`        // 工具调用
}

// ChatChoice 模型生成的 completion
type ChatChoice struct {
	FinishReason ChatFinishReason       `json:"finish_reason,omitempty"` // 模型停止生成 token 的原因
	Index        int                    `json:"index,omitempty"`         // 该 completion 在模型生成的 completion 的选择列表中的索引
	LogProbs     *ChatLogProbs          `json:"logprobs,omitempty"`      // 该 choice 的对数概率信息
	Message      *ChatCompletionMessage `json:"message,omitempty"`       // 模型生成的 completion 消息
	Delta        *ChatCompletionMessage `json:"delta,omitempty"`         // 流式传输的增量信息
}

// CompletionTokensDetails completion tokens 的详细信息
type CompletionTokensDetails struct {
	TextTokens               int `json:"text_tokens,omitempty"`                // 模型生成的文本 token 数量
	AudioTokens              int `json:"audio_tokens,omitempty"`               // 模型生成的音频 token 数量
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`           // 模型用于推理而生成的 token 数量
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"` // 使用预测输出时，预测中出现在完成结果中的 token 数量
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"` // 使用预测输出时，预测中未出现在完成结果中的 token 数量
}

// PromptTokensDetails prompt tokens 的详细信息
type PromptTokensDetails struct {
	AudioTokens  int `json:"audio_tokens,omitempty"`  // prompt 中存在的音频输入 token 数
	CachedTokens int `json:"cached_tokens,omitempty"` // prompt 中存在的缓存 token 数
	TextTokens   int `json:"text_tokens,omitempty"`   // prompt 中存在的文本 token 数
	ImageTokens  int `json:"image_tokens,omitempty"`  // prompt 中存在的图像 token 数
	VideoTokens  int `json:"video_tokens,omitempty"`  // prompt 中存在的视频 token 数
}

// ChatUsage 该对话补全请求的用量信息
type ChatUsage struct {
	CompletionTokens        int                      `json:"completion_tokens,omitempty"`         // 模型 completion 产生的 token 数
	PromptTokens            int                      `json:"prompt_tokens,omitempty"`             // 用户 prompt 所包含的 token 数
	PromptCacheHitTokens    int                      `json:"prompt_cache_hit_tokens,omitempty"`   // 用户 prompt 中，命中上下文缓存的 token 数
	PromptCacheMissTokens   int                      `json:"prompt_cache_miss_tokens,omitempty"`  // 用户 prompt 中，未命中上下文缓存的 token 数
	TotalTokens             int                      `json:"total_tokens,omitempty"`              // 该请求中，所有 token 的数量（prompt + completion）
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // completion tokens 的详细信息
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // prompt tokens 的详细信息
}

// ChatBaseResponse 聊天响应基础信息
type ChatBaseResponse struct {
	provider   string // 用于反序列化数据时，处理差异化数据
	streamable bool   // 用于反序列化数据时，处理流式传输数据
	// 以下字段为响应数据
	Choices           []ChatChoice            `json:"choices,omitempty"`            // 模型生成的 completion 的选择列表
	Created           int64                   `json:"created,omitempty"`            // 创建聊天完成时的 Unix 时间戳（以秒为单位）
	ID                string                  `json:"id,omitempty"`                 // 该对话的唯一标识符
	Model             string                  `json:"model,omitempty"`              // 生成该 completion 的模型名
	Object            string                  `json:"object,omitempty"`             // 对象的类型
	ServiceTier       string                  `json:"service_tier,omitempty"`       // 用于处理请求的服务层级
	SystemFingerprint string                  `json:"system_fingerprint,omitempty"` // 此指纹表示模型运行的后端配置。可以与 seed 请求参数一起使用，以了解何时进行了可能影响确定性的后端更改
	Usage             *ChatUsage              `json:"usage,omitempty"`              // 该对话补全请求的用量信息
	StreamStats       *httpclient.StreamStats `json:"stream_stats,omitempty"`       // 流式传输统计信息
}

// SetProvider 设置提供商
func (c *ChatBaseResponse) SetProvider(provider string) {
	c.provider = provider
}

// SetStreamable 设置流式传输
func (c *ChatBaseResponse) SetStreamable(streamable bool) {
	c.streamable = streamable
}

// SetStreamStats 设置流式传输统计信息
func (c *ChatBaseResponse) SetStreamStats(stats httpclient.StreamStats) {
	c.StreamStats = &stats
}

// UnmarshalJSON 反序列化JSON
func (c *ChatBaseResponse) UnmarshalJSON(data []byte) (err error) {
	switch consts.Provider(c.provider) {
	case consts.AliBL:
		return c.unmarshalAliBL(data)
	default:
		// 默认反序列化
		type Alias ChatBaseResponse
		temp := (*Alias)(c)
		return json.Unmarshal(data, temp)
	}
}

// unmarshalAliBL 反序列化阿里百炼响应
func (c *ChatBaseResponse) unmarshalAliBL(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	if err = json.Unmarshal(data, &rawMap); err != nil {
		return
	}
	// 解析 output 数据
	if _, ok := rawMap["output"]; ok {
		var outputMap map[string]json.RawMessage
		if err = json.Unmarshal(rawMap["output"], &outputMap); err != nil {
			return
		}
		// 解析 search_info 数据
		var annotations []ChatAnnotation
		if _, ok := outputMap["search_info"]; ok {
			var tmpSearchInfo struct {
				SearchResults []struct {
					SiteName string `json:"site_name,omitempty"` // 搜索结果来源的网站名称
					Icon     string `json:"icon,omitempty"`      // 来源网站的图标URL，如果没有图标则为空字符串
					Index    int    `json:"index,omitempty"`     // 搜索结果的序号，表示该搜索结果在search_results中的索引
					Title    string `json:"title,omitempty"`     // 搜索结果的标题
					URL      string `json:"url,omitempty"`       // 搜索结果的链接地址
				} `json:"search_results,omitempty"` // 联网搜索到的结果
			}
			if err = json.Unmarshal(outputMap["search_info"], &tmpSearchInfo); err != nil {
				return
			}
			// 创建 annotations 数组
			annotations = make([]ChatAnnotation, len(tmpSearchInfo.SearchResults))
			for i, searchResult := range tmpSearchInfo.SearchResults {
				annotations[i] = ChatAnnotation{
					Index: searchResult.Index,
					URLCitation: &ChatAnnotationURLCitation{
						SiteName: searchResult.SiteName,
						Icon:     searchResult.Icon,
						Title:    searchResult.Title,
						URL:      searchResult.URL,
					},
				}
			}
		}
		// 解析 choices 数据
		if _, ok := outputMap["choices"]; ok {
			var tmpChoices []struct {
				FinishReason ChatFinishReason `json:"finish_reason,omitempty"` // 模型停止生成 token 的原因
				Message      *struct {
					Role             string          `json:"role,omitempty"`              // 输出消息的角色
					Content          json.RawMessage `json:"content,omitempty"`           // 输出消息的内容
					ReasoningContent string          `json:"reasoning_content,omitempty"` // 模型的深度思考内容
					ToolCalls        []ToolCalls     `json:"tool_calls,omitempty"`        // 如果模型需要调用工具，则会生成tool_calls参数
				} `json:"message,omitempty"` // 模型输出的消息对象
				LogProbs *ChatLogProbs `json:"logprobs,omitempty"` // 当前 choices 对象的概率信息
			}
			if err = json.Unmarshal(outputMap["choices"], &tmpChoices); err != nil {
				return
			}
			// 创建 choices 数组
			c.Choices = make([]ChatChoice, len(tmpChoices))
			for i, tmpChoice := range tmpChoices {
				// 创建 message 对象
				message := &ChatCompletionMessage{
					Role: tmpChoice.Message.Role,
					Content: func() string {
						// 尝试将 Content 转换为 []map[string]any 类型
						var contentMapList []map[string]any
						if e := json.Unmarshal(tmpChoice.Message.Content, &contentMapList); e == nil {
							var content string
							for _, contentMap := range contentMapList {
								if text, ok := contentMap["text"].(string); ok {
									if content == "" {
										content = text
									} else {
										content = fmt.Sprintf("%s\n%s", content, text)
									}
								}
							}
							return content
						}
						// 尝试将 Content 转换为 string 类型
						var content string
						if e := json.Unmarshal(tmpChoice.Message.Content, &content); e == nil {
							return content
						}
						// 如果转换失败，返回空字符串
						return ""
					}(),
					ReasoningContent: tmpChoice.Message.ReasoningContent,
					Annotations:      annotations,
					ToolCalls:        tmpChoice.Message.ToolCalls,
				}
				c.Choices[i] = ChatChoice{
					FinishReason: tmpChoice.FinishReason,
					LogProbs:     tmpChoice.LogProbs,
				}
				// 如果流式传输，则设置 Delta 字段，否则设置 Message 字段
				if c.streamable {
					c.Choices[i].Delta = message
				} else {
					c.Choices[i].Message = message
				}
			}
		}
	}
	// 解析 usage 数据
	if _, ok := rawMap["usage"]; ok {
		var tmpUsage struct {
			InputTokens        int `json:"input_tokens,omitempty"`  // 用户输入内容转换成Token后的长度
			OutputTokens       int `json:"output_tokens,omitempty"` // 模型输出内容转换成Token后的长度
			InputTokensDetails *struct {
				TextTokens  int `json:"text_tokens,omitempty"`  // 输入的文本转换为Token后的长度
				ImageTokens int `json:"image_tokens,omitempty"` // 输入的图像转换为Token后的长度
				VideoTokens int `json:"video_tokens,omitempty"` // 输入的视频文件或图像列表转换为Token后的长度
			} `json:"input_tokens_details,omitempty"` // 输入内容转换成Token后的长度详情
			TotalTokens         int `json:"total_tokens,omitempty"` // 当输入为纯文本时返回该字段，为input_tokens与output_tokens之和
			ImageTokens         int `json:"image_tokens,omitempty"` // 输入内容包含image时返回该字段。为用户输入图片内容转换成Token后的长度
			VideoTokens         int `json:"video_tokens,omitempty"` // 输入内容包含video时返回该字段。为用户输入视频内容转换成Token后的长度
			AudioTokens         int `json:"audio_tokens,omitempty"` // 输入内容包含audio时返回该字段。为用户输入音频内容转换成Token后的长度
			OutputTokensDetails *struct {
				TextTokens      int `json:"text_tokens,omitempty"`      // 输出的文本转换为Token后的长度
				ReasoningTokens int `json:"reasoning_tokens,omitempty"` // 模型思考过程转换为Token后的长度
			} `json:"output_tokens_details,omitempty"` // 输出内容转换成 Token后的长度详情
			PromptTokensDetails *PromptTokensDetails `json:"prompt_tokens_details,omitempty"` // 输入 Token 的细粒度分类
		}
		if err = json.Unmarshal(rawMap["usage"], &tmpUsage); err != nil {
			return
		}
		// 创建 Usage 对象
		c.Usage = &ChatUsage{
			CompletionTokens: tmpUsage.OutputTokens,
			PromptTokens:     tmpUsage.InputTokens,
			TotalTokens:      tmpUsage.TotalTokens,
		}
		if tmpUsage.OutputTokensDetails != nil {
			c.Usage.CompletionTokensDetails = &CompletionTokensDetails{
				TextTokens:      tmpUsage.OutputTokensDetails.TextTokens,
				ReasoningTokens: tmpUsage.OutputTokensDetails.ReasoningTokens,
			}
		}
		if tmpUsage.PromptTokensDetails != nil {
			c.Usage.PromptTokensDetails = tmpUsage.PromptTokensDetails
		}
		if tmpUsage.InputTokensDetails != nil {
			if c.Usage.PromptTokensDetails == nil {
				c.Usage.PromptTokensDetails = &PromptTokensDetails{}
			}
			c.Usage.PromptTokensDetails.TextTokens = tmpUsage.InputTokensDetails.TextTokens
			c.Usage.PromptTokensDetails.ImageTokens = tmpUsage.InputTokensDetails.ImageTokens
			c.Usage.PromptTokensDetails.VideoTokens = tmpUsage.InputTokensDetails.VideoTokens
		}
	}
	return
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ChatBaseResponse
	httpclient.HttpHeader
}

// ChatResponseStream 流式传输的聊天响应
type ChatResponseStream struct {
	*httpclient.StreamReader[ChatBaseResponse]
}
