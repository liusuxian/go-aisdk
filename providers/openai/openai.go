/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:56:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 20:59:08
 * @Description: OpenAI服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package openai

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
				consts.OpenAIO1Mini:                         consts.ModelFeature(0),
				consts.OpenAIO1Mini20240912:                 consts.ModelFeature(0),
				consts.OpenAIO1Preview:                      consts.ModelFeature(0),
				consts.OpenAIO1Preview20240912:              consts.ModelFeature(0),
				consts.OpenAIO1:                             consts.ModelFeature(0),
				consts.OpenAIO1_20241217:                    consts.ModelFeature(0),
				consts.OpenAIO1Pro:                          consts.ModelFeature(0),
				consts.OpenAIO1Pro20250319:                  consts.ModelFeature(0),
				consts.OpenAIO3:                             consts.ModelFeature(1),
				consts.OpenAIO3_20250416:                    consts.ModelFeature(1),
				consts.OpenAIO3Mini:                         consts.ModelFeature(1),
				consts.OpenAIO3Mini20250131:                 consts.ModelFeature(1),
				consts.OpenAIO4Mini:                         consts.ModelFeature(1),
				consts.OpenAIO4Mini20250416:                 consts.ModelFeature(1),
				consts.OpenAIGPT4_32K0613:                   consts.ModelFeature(0),
				consts.OpenAIGPT4_32K0314:                   consts.ModelFeature(0),
				consts.OpenAIGPT4_32K:                       consts.ModelFeature(0),
				consts.OpenAIGPT4_0613:                      consts.ModelFeature(0),
				consts.OpenAIGPT4_0314:                      consts.ModelFeature(0),
				consts.OpenAIGPT4o:                          consts.ModelFeature(1),
				consts.OpenAIGPT4o20240513:                  consts.ModelFeature(1),
				consts.OpenAIGPT4o20240806:                  consts.ModelFeature(1),
				consts.OpenAIGPT4o20241120:                  consts.ModelFeature(1),
				consts.OpenAIChatGPT4oLatest:                consts.ModelFeature(1),
				consts.OpenAIGPT4oMini:                      consts.ModelFeature(1),
				consts.OpenAIGPT4oMini20240718:              consts.ModelFeature(1),
				consts.OpenAIGPT4oSearchPreview:             consts.ModelFeature(1),
				consts.OpenAIGPT4oSearchPreview20250311:     consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniSearchPreview:         consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniSearchPreview20250311: consts.ModelFeature(1),
				consts.OpenAIGPT4Turbo:                      consts.ModelFeature(1),
				consts.OpenAIGPT4TurboPreview:               consts.ModelFeature(1),
				consts.OpenAIGPT4Turbo20240409:              consts.ModelFeature(1),
				consts.OpenAIGPT4_0125Preview:               consts.ModelFeature(1),
				consts.OpenAIGPT4_1106Preview:               consts.ModelFeature(1),
				consts.OpenAIGPT4VisionPreview:              consts.ModelFeature(1),
				consts.OpenAIGPT4:                           consts.ModelFeature(0),
				consts.OpenAIGPT4Dot1:                       consts.ModelFeature(1),
				consts.OpenAIGPT4Dot1_20250414:              consts.ModelFeature(1),
				consts.OpenAIGPT4Dot1Mini:                   consts.ModelFeature(1),
				consts.OpenAIGPT4Dot1Mini20250414:           consts.ModelFeature(1),
				consts.OpenAIGPT4Dot1Nano:                   consts.ModelFeature(1),
				consts.OpenAIGPT4Dot1Nano20250414:           consts.ModelFeature(1),
				consts.OpenAIGPT4Dot5Preview:                consts.ModelFeature(1),
				consts.OpenAIGPT4Dot5Preview20250227:        consts.ModelFeature(1),
				consts.OpenAIGPT3Dot5Turbo0125:              consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo1106:              consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo0613:              consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo0301:              consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo16k:               consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo16K0613:           consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5Turbo:                  consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5TurboInstruct:          consts.ModelFeature(0),
				consts.OpenAIGPT3Dot5TurboInstruct0914:      consts.ModelFeature(0),
				consts.OpenAIDavinci002:                     consts.ModelFeature(0),
				consts.OpenAIBabbage002:                     consts.ModelFeature(0),
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20241001:        consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20241217:        consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20250603:        consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview:             consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20241001:     consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20241217:     consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20250603:     consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniAudioPreview:            consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniAudioPreview20241217:    consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniRealtimePreview:         consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniRealtimePreview20241217: consts.ModelFeature(1),
			},
			consts.ImageModel: {
				// image
				consts.OpenAIDallE2:    consts.ModelFeature(1),
				consts.OpenAIDallE3:    consts.ModelFeature(1),
				consts.OpenAIGPTImage1: consts.ModelFeature(1),
			},
			consts.AudioModel: {
				// audio
				consts.OpenAITTS1:                consts.ModelFeature(1),
				consts.OpenAITTS1_1106:           consts.ModelFeature(1),
				consts.OpenAITTS1HD:              consts.ModelFeature(1),
				consts.OpenAITTS1HD1106:          consts.ModelFeature(1),
				consts.OpenAIWhisper1:            consts.ModelFeature(1),
				consts.OpenAIGPT4oTranscribe:     consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniTranscribe: consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniTTS:        consts.ModelFeature(1),
				// chat, audio
				consts.OpenAIGPT4oAudioPreview:                consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20241001:        consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20241217:        consts.ModelFeature(1),
				consts.OpenAIGPT4oAudioPreview20250603:        consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview:             consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20241001:     consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20241217:     consts.ModelFeature(1),
				consts.OpenAIGPT4oRealtimePreview20250603:     consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniAudioPreview:            consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniAudioPreview20241217:    consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniRealtimePreview:         consts.ModelFeature(1),
				consts.OpenAIGPT4oMiniRealtimePreview20241217: consts.ModelFeature(1),
			},
			// moderation
			consts.ModerationModel: {
				consts.OpenAIOmniModerationLatest:   consts.ModelFeature(0),
				consts.OpenAIOmniModeration20240926: consts.ModelFeature(0),
			},
			// embed
			consts.EmbedModel: {
				consts.OpenAITextEmbedding3Small: consts.ModelFeature(0),
				consts.OpenAITextEmbedding3Large: consts.ModelFeature(0),
				consts.OpenAITextEmbeddingAda002: consts.ModelFeature(0),
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
