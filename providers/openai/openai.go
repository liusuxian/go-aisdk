/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-10-16 16:40:33
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

import (
	"context"
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/httpclient"
	"github.com/liusuxian/go-aisdk/loadbalancer"
	"github.com/liusuxian/go-aisdk/models"
	"github.com/liusuxian/go-aisdk/providers/common"
	"net/http"
)

// openAIProvider OpenAI提供商
type openAIProvider struct {
	core.DefaultProviderService
	supportedModels map[consts.ModelType]map[string]consts.ModelFeature // 支持的模型
	providerConfig  *conf.ProviderConfig                                // 提供商配置
	lb              *loadbalancer.LoadBalancer                          // 负载均衡器
}

var (
	openaiService *openAIProvider // OpenAI提供商实例
)

const (
	apiModels = "/models"
)

// init 包初始化时创建 openAIProvider 实例并注册到工厂
func init() {
	openaiService = &openAIProvider{
		supportedModels: map[consts.ModelType]map[string]consts.ModelFeature{
			consts.ChatModel: {
				// chat
				consts.OpenAIO1Mini:                         consts.ModelFeatureNone,
				consts.OpenAIO1Mini20240912:                 consts.ModelFeatureNone,
				consts.OpenAIO1Preview:                      consts.ModelFeatureNone,
				consts.OpenAIO1Preview20240912:              consts.ModelFeatureNone,
				consts.OpenAIO1:                             consts.ModelFeatureNone,
				consts.OpenAIO1_20241217:                    consts.ModelFeatureNone,
				consts.OpenAIO1Pro:                          consts.ModelFeatureNone,
				consts.OpenAIO1Pro20250319:                  consts.ModelFeatureNone,
				consts.OpenAIO3:                             consts.ModelFeatureNone,
				consts.OpenAIO3_20250416:                    consts.ModelFeatureNone,
				consts.OpenAIO3Mini:                         consts.ModelFeatureNone,
				consts.OpenAIO3Mini20250131:                 consts.ModelFeatureNone,
				consts.OpenAIO4Mini:                         consts.ModelFeatureNone,
				consts.OpenAIO4Mini20250416:                 consts.ModelFeatureNone,
				consts.OpenAIGPT4_32K0613:                   consts.ModelFeatureNone,
				consts.OpenAIGPT4_32K0314:                   consts.ModelFeatureNone,
				consts.OpenAIGPT4_32K:                       consts.ModelFeatureNone,
				consts.OpenAIGPT4_0613:                      consts.ModelFeatureNone,
				consts.OpenAIGPT4_0314:                      consts.ModelFeatureNone,
				consts.OpenAIGPT4o:                          consts.ModelFeatureNone,
				consts.OpenAIGPT4o20240513:                  consts.ModelFeatureNone,
				consts.OpenAIGPT4o20240806:                  consts.ModelFeatureNone,
				consts.OpenAIGPT4o20241120:                  consts.ModelFeatureNone,
				consts.OpenAIChatGPT4oLatest:                consts.ModelFeatureNone,
				consts.OpenAIGPT4oMini:                      consts.ModelFeatureNone,
				consts.OpenAIGPT4oMini20240718:              consts.ModelFeatureNone,
				consts.OpenAIGPT4oSearchPreview:             consts.ModelFeatureNone,
				consts.OpenAIGPT4oSearchPreview20250311:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniSearchPreview:         consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniSearchPreview20250311: consts.ModelFeatureNone,
				consts.OpenAIGPT4Turbo:                      consts.ModelFeatureNone,
				consts.OpenAIGPT4TurboPreview:               consts.ModelFeatureNone,
				consts.OpenAIGPT4Turbo20240409:              consts.ModelFeatureNone,
				consts.OpenAIGPT4_0125Preview:               consts.ModelFeatureNone,
				consts.OpenAIGPT4_1106Preview:               consts.ModelFeatureNone,
				consts.OpenAIGPT4VisionPreview:              consts.ModelFeatureNone,
				consts.OpenAIGPT4:                           consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1:                       consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1_20250414:              consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1Mini:                   consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1Mini20250414:           consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1Nano:                   consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot1Nano20250414:           consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot5Preview:                consts.ModelFeatureNone,
				consts.OpenAIGPT4Dot5Preview20250227:        consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo0125:              consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo1106:              consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo0613:              consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo0301:              consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo16k:               consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo16K0613:           consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5Turbo:                  consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5TurboInstruct:          consts.ModelFeatureNone,
				consts.OpenAIGPT3Dot5TurboInstruct0914:      consts.ModelFeatureNone,
				consts.OpenAIDavinci002:                     consts.ModelFeatureNone,
				consts.OpenAIBabbage002:                     consts.ModelFeatureNone,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20241001:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20241217:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20250603:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview:             consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20241001:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20241217:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20250603:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniAudioPreview:            consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniRealtimePreview:         consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: consts.ModelFeatureNone,
			},
			consts.ImageModel: {
				// image
				consts.OpenAIDallE2:    consts.ModelFeatureNone,
				consts.OpenAIDallE3:    consts.ModelFeatureNone,
				consts.OpenAIGPTImage1: consts.ModelFeatureNone,
			},
			consts.AudioModel: {
				// audio
				consts.OpenAITTS1:                consts.ModelFeatureNone,
				consts.OpenAITTS1_1106:           consts.ModelFeatureNone,
				consts.OpenAITTS1HD:              consts.ModelFeatureNone,
				consts.OpenAITTS1HD1106:          consts.ModelFeatureNone,
				consts.OpenAIWhisper1:            consts.ModelFeatureNone,
				consts.OpenAIGPT4oTranscribe:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniTranscribe: consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniTTS:        consts.ModelFeatureNone,
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20241001:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20241217:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oAudioPreview20250603:        consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview:             consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20241001:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20241217:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oRealtimePreview20250603:     consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniAudioPreview:            consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniAudioPreview20241217:    consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniRealtimePreview:         consts.ModelFeatureNone,
				consts.OpenAIGPT4oMiniRealtimePreview20241217: consts.ModelFeatureNone,
			},
			// moderation
			consts.ModerationModel: {
				consts.OpenAIOmniModerationLatest:   consts.ModelFeatureNone,
				consts.OpenAIOmniModeration20240926: consts.ModelFeatureNone,
			},
			// embed
			consts.EmbedModel: {
				consts.OpenAITextEmbedding3Small: consts.ModelFeatureNone,
				consts.OpenAITextEmbedding3Large: consts.ModelFeatureNone,
				consts.OpenAITextEmbeddingAda002: consts.ModelFeatureNone,
			},
		},
	}
	core.RegisterProvider(consts.OpenAI, openaiService)
}

// GetSupportedModels 获取支持的模型
func (s *openAIProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]consts.ModelFeature) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *openAIProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}

// ListModels 列出模型
func (s *openAIProvider) ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Provider: consts.OpenAI,
		Method:   http.MethodGet,
		BaseURL:  s.providerConfig.BaseURL,
		ApiPath:  apiModels,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
	})
	return
}
