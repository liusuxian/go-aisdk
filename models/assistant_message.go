/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:14:56
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/Mrzhouyl/go-aisdk/internal/utils"

// ChatAssistantMsgAudio 音频
type ChatAssistantMsgAudio struct {
	// 音频ID
	//
	// 提供商支持: OpenAI
	ID string `json:"id,omitempty" providers:"openai"`
}

// ChatAssistantMsgPartType 多模态内容类型
type ChatAssistantMsgPartType string

const (
	// 提供商支持: OpenAI
	ChatAssistantMsgPartTypeText ChatAssistantMsgPartType = "text"
	// 提供商支持: OpenAI
	ChatAssistantMsgPartTypeRefusal ChatAssistantMsgPartType = "refusal"
)

// ChatAssistantMsgPart 多模态内容
type ChatAssistantMsgPart struct {
	// 内容类型
	//
	// 提供商支持: OpenAI
	Type ChatAssistantMsgPartType `json:"type,omitempty" providers:"openai"`
	// 文本内容
	//
	// 提供商支持: OpenAI
	Text string `json:"text,omitempty" providers:"openai"`
	// 拒绝消息
	//
	// 提供商支持: OpenAI
	Refusal string `json:"refusal,omitempty" providers:"openai"`
}

// AssistantMessage 助手消息
type AssistantMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 音频
	//
	// 提供商支持: OpenAI
	Audio *ChatAssistantMsgAudio `json:"audio,omitempty" providers:"openai"`
	// 文本内容
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Content string `json:"content,omitempty" providers:"openai,deepseek,alibl"`
	// 多模态内容
	//
	// 提供商支持: OpenAI
	MultimodalContent []ChatAssistantMsgPart `json:"multimodal_content,omitempty" providers:"openai" copyto:"Content"`
	// 消息角色
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Role string `json:"role,omitempty" providers:"openai,deepseek,alibl" default:"assistant"`
	// 参与者名称
	//
	// 提供商支持: OpenAI | DeepSeek
	Name string `json:"name,omitempty" providers:"openai,deepseek"`
	// 拒绝消息
	//
	// 提供商支持: OpenAI
	Refusal string `json:"refusal,omitempty" providers:"openai"`
	// 工具调用
	//
	// 提供商支持: OpenAI | DeepSeek
	ToolCalls []ToolCalls `json:"tool_calls,omitempty" providers:"openai,deepseek"`
	// 设置此参数为 true，来强制模型在其回答中以此 assistant 消息中提供的前缀内容开始
	//
	// 提供商支持: DeepSeek | AliBL
	Prefix *bool `json:"prefix,omitempty" providers:"deepseek,alibl" mapping:"alibl:partial"`
	// 用于模型在对话前缀续写功能下，作为最后一条 assistant 思维链内容的输入。使用此功能时，prefix 参数必须设置为 true
	//
	// 提供商支持: DeepSeek
	ReasoningContent string `json:"reasoning_content,omitempty" providers:"deepseek"`
}

// SetProvider 设置提供商
func (m *AssistantMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m AssistantMessage) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(m.provider).Serialize(m)
}
