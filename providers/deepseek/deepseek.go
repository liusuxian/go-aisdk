/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-09 17:13:03
 * @Description: DeepSeek服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"net/http"

	"github.com/Mrzhouyl/go-aisdk/conf"
	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/core"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/loadbalancer"
	"github.com/Mrzhouyl/go-aisdk/models"
	"github.com/Mrzhouyl/go-aisdk/providers/common"
)

// deepseekProvider DeepSeek提供商
type deepseekProvider struct {
	core.DefaultProviderService
	supportedModels map[consts.ModelType]map[string]consts.ModelFeature // 支持的模型
	providerConfig  *conf.ProviderConfig                                // 提供商配置
	lb              *loadbalancer.LoadBalancer                          // 负载均衡器
}

var (
	deepseekService *deepseekProvider // deepseek提供商实例
)

const (
	apiModels = "/models"
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	deepseekService = &deepseekProvider{
		supportedModels: map[consts.ModelType]map[string]consts.ModelFeature{
			consts.ChatModel: {
				// chat
				consts.DeepSeekChat:     consts.ModelFeatureNone,
				consts.DeepSeekReasoner: consts.ModelFeatureReasoning,
			},
		},
	}
	core.RegisterProvider(consts.DeepSeek, deepseekService)
}

// GetSupportedModels 获取支持的模型
func (s *deepseekProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]consts.ModelFeature) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *deepseekProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *deepseekProvider) ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Provider: consts.DeepSeek,
		Method:   http.MethodGet,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiModels,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
	})
	return
}
