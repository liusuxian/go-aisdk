/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:42:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-01 00:27:25
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import "github.com/Mrzhouyl/go-aisdk/consts"

// ModelInfo 模型信息结构
type ModelInfo struct {
	Provider  consts.Provider  `json:"provider,omitempty"`   // 提供商
	ModelType consts.ModelType `json:"model_type,omitempty"` // 模型类型
	Model     string           `json:"model,omitempty"`      // 模型名称
}

// UserInfo 用户信息结构
type UserInfo struct {
	User string `json:"user,omitempty" providers:"openai"` // 代表你的终端用户的唯一标识符
}
