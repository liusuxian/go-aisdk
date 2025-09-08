/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:15:57
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/Mrzhouyl/go-aisdk/internal/utils"

// DeveloperMessage 开发者消息
type DeveloperMessage struct {
	// 用于序列化参数时，处理差异化参数
	provider string
	// 文本内容
	//
	// 提供商支持: OpenAI
	Content string `json:"content,omitempty" providers:"openai"`
	// 消息角色
	//
	// 提供商支持: OpenAI
	Role string `json:"role,omitempty" providers:"openai" default:"developer"`
	// 参与者名称
	//
	// 提供商支持: OpenAI
	Name string `json:"name,omitempty" providers:"openai"`
}

// SetProvider 设置提供商
func (m *DeveloperMessage) SetProvider(provider string) {
	m.provider = provider
}

// MarshalJSON 序列化JSON
func (m DeveloperMessage) MarshalJSON() (b []byte, err error) {
	return utils.NewSerializer(m.provider).Serialize(m)
}
