/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-16 17:58:22
 * @Description: DeepSeek服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"github.com/liusuxian/aisdk/conf"
	"github.com/liusuxian/aisdk/consts"
	"github.com/liusuxian/aisdk/core"
	"github.com/liusuxian/aisdk/models"
)

// deepseekProvider DeepSeek提供商
type deepseekProvider struct {
	supportedModels   map[consts.ModelType][]string // 支持的模型
	providerConfig    *conf.ProviderConfig          // 提供商配置
	connectionOptions *conf.ConnectionOptions       // 连接选项
}

var (
	deepseekService *deepseekProvider // deepseek提供商实例
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	deepseekService = &deepseekProvider{
		supportedModels: map[consts.ModelType][]string{
			consts.ChatModel: {
				consts.DeepSeekChat,
				consts.DeepSeekReasoner,
			},
		},
	}
	core.RegisterProvider(consts.DeepSeek, deepseekService)
}

// GetSupportedModels 获取支持的模型
func (s *deepseekProvider) GetSupportedModels() (supportedModels map[consts.ModelType][]string) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *deepseekProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
}

// InitializeConnectionOptions 初始化连接选项
func (s *deepseekProvider) InitializeConnectionOptions(options *conf.ConnectionOptions) {
	s.connectionOptions = options
}

// TODO: CreateChatCompletion 创建聊天
func (s *deepseekProvider) CreateChatCompletion(ctx context.Context, request models.ChatRequest) (response models.ChatResponse, err error) {
	return
}
