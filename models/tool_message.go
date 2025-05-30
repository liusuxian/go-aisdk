/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 18:54:06
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "encoding/json"

// ToolMessage 工具消息
type ToolMessage struct {
	Content    string `json:"content"`      // 文本内容
	ToolCallID string `json:"tool_call_id"` // 工具调用ID
}

// GetRole 获取消息角色
func (m ToolMessage) GetRole() (role string) { return "tool" }

// MarshalJSON 序列化JSON
func (m ToolMessage) MarshalJSON() (b []byte, err error) {
	type Alias ToolMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		Alias
	}{
		Role:  "tool",
		Alias: Alias(m),
	})
}
