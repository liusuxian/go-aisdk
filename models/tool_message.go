/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:16:59
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/Mrzhouyl/go-aisdk/internal/utils"

// ToolMessage 工具消息
type ToolMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 文本内容
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Content string `json:"content,omitempty" providers:"openai,deepseek,alibl"`
	// 消息角色
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	Role string `json:"role,omitempty" providers:"openai,deepseek,alibl" default:"tool"`
	// 工具调用ID
	//
	// 提供商支持: OpenAI | DeepSeek | AliBL
	ToolCallID string `json:"tool_call_id,omitempty" providers:"openai,deepseek,alibl"`
}

// SetProvider 设置提供商
func (m *ToolMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m ToolMessage) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(m.provider).Serialize(m)
}
