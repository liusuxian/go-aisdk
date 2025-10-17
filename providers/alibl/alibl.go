/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:31:10
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-10-17 19:04:41
 * @Description: AliBL服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package alibl

import (
	"github.com/liusuxian/go-aisdk/conf"
	"github.com/liusuxian/go-aisdk/consts"
	"github.com/liusuxian/go-aisdk/core"
	"github.com/liusuxian/go-aisdk/loadbalancer"
)

// aliblProvider AliBL提供商
type aliblProvider struct {
	core.DefaultProviderService
	supportedModels map[consts.ModelType]map[string]consts.ModelFeature // 支持的模型
	providerConfig  *conf.ProviderConfig                                // 提供商配置
	lb              *loadbalancer.LoadBalancer                          // 负载均衡器
}

var (
	aliblService *aliblProvider // alibl提供商实例
)

// init 包初始化时创建 deepseekProvider 实例并注册到工厂
func init() {
	aliblService = &aliblProvider{
		supportedModels: map[consts.ModelType]map[string]consts.ModelFeature{
			consts.ChatModel: {
				// chat
				consts.AliBLQwqPlus:                        consts.ModelFeatureReasoningStream,
				consts.AliBLQwqPlusLatest:                  consts.ModelFeatureReasoningStream,
				consts.AliBLQwqPlus20250305:                consts.ModelFeatureReasoningStream,
				consts.AliBLQwen3Max:                       consts.ModelFeatureNone,
				consts.AliBLQwen3Max20250923:               consts.ModelFeatureNone,
				consts.AliBLQwen3MaxPreview:                consts.ModelFeatureNone,
				consts.AliBLQwenMax:                        consts.ModelFeatureNone,
				consts.AliBLQwenMaxLatest:                  consts.ModelFeatureNone,
				consts.AliBLQwenMax20250125:                consts.ModelFeatureNone,
				consts.AliBLQwenMax20240919:                consts.ModelFeatureNone,
				consts.AliBLQwenMax20240428:                consts.ModelFeatureNone,
				consts.AliBLQwenMax20240403:                consts.ModelFeatureNone,
				consts.AliBLQwenPlus:                       consts.ModelFeatureReasoning,
				consts.AliBLQwenPlusLatest:                 consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250911:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250728:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250714:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250428:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250125:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20250112:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20241220:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20241127:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20241125:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20240919:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20240806:               consts.ModelFeatureReasoning,
				consts.AliBLQwenPlus20240723:               consts.ModelFeatureReasoning,
				consts.AliBLQwenFlash:                      consts.ModelFeatureReasoning,
				consts.AliBLQwenFlash20250923:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo:                      consts.ModelFeatureReasoning,
				consts.AliBLQwenTurboLatest:                consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20250715:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20250428:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20250211:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20241101:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20240919:              consts.ModelFeatureReasoning,
				consts.AliBLQwenTurbo20240624:              consts.ModelFeatureReasoning,
				consts.AliBLQwenLong:                       consts.ModelFeatureNone,
				consts.AliBLQwenLongLatest:                 consts.ModelFeatureNone,
				consts.AliBLQwenLong20250125:               consts.ModelFeatureNone,
				consts.AliBLQwen3OmniFlash:                 consts.ModelFeatureAdvancedStream,
				consts.AliBLQwen3OmniFlash20250915:         consts.ModelFeatureAdvancedStream,
				consts.AliBLQwenOmniTurbo:                  consts.ModelFeatureMultimodal,
				consts.AliBLQwenOmniTurboLatest:            consts.ModelFeatureMultimodal,
				consts.AliBLQwenOmniTurbo20250326:          consts.ModelFeatureMultimodal,
				consts.AliBLQwenOmniTurbo20250119:          consts.ModelFeatureMultimodal,
				consts.AliBLQwen3OmniFlashRealtime:         consts.ModelFeatureMultimodal, // TODO: 需要 WebSocket
				consts.AliBLQwen3OmniFlashRealtime20250915: consts.ModelFeatureMultimodal, // TODO: 需要 WebSocket
				consts.AliBLQwenOmniTurboRealtime:          consts.ModelFeatureMultimodal, // TODO: 需要 WebSocket
				consts.AliBLQwenOmniTurboRealtimeLatest:    consts.ModelFeatureMultimodal, // TODO: 需要 WebSocket
				consts.AliBLQwenOmniTurboRealtime20250508:  consts.ModelFeatureMultimodal, // TODO: 需要 WebSocket
				consts.AliBLQvqMax:                         consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqMaxLatest:                   consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqMax20250515:                 consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqMax20250325:                 consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqPlus:                        consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqPlusLatest:                  consts.ModelFeatureAdvancedStream,
				consts.AliBLQvqPlus20250515:                consts.ModelFeatureAdvancedStream,
				consts.AliBLQwenVlMax:                      consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMaxLatest:                consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20250813:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20250408:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20250402:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20250125:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20241230:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20241119:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20241030:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlMax20240809:              consts.ModelFeatureMultimodalStream,
				consts.AliBLQwen3VlPlus:                    consts.ModelFeatureAdvancedStream,
				consts.AliBLQwen3VlPlus20250923:            consts.ModelFeatureAdvancedStream,
				consts.AliBLQwenVlPlus:                     consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlusLatest:               consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20250815:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20250710:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20250507:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20250125:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20250102:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlPlus20240809:             consts.ModelFeatureMultimodalStream,
				consts.AliBLQwenVlOcr:                      consts.ModelFeatureMultimodal,
				consts.AliBLQwenVlOcrLatest:                consts.ModelFeatureMultimodal,
				consts.AliBLQwenVlOcr20250828:              consts.ModelFeatureMultimodal,
				consts.AliBLQwenVlOcr20250413:              consts.ModelFeatureMultimodal,
				consts.AliBLQwenVlOcr20241028:              consts.ModelFeatureMultimodal,
				consts.AliBLQwenAudioTurbo:                 consts.ModelFeatureMultimodal,
				consts.AliBLQwenAudioTurboLatest:           consts.ModelFeatureMultimodal,
				consts.AliBLQwenAudioTurbo20241204:         consts.ModelFeatureMultimodal,
				consts.AliBLQwenAudioTurbo20240807:         consts.ModelFeatureMultimodal,
				consts.AliBLQwen3AsrFlash:                  consts.ModelFeatureMultimodal,
				consts.AliBLQwen3AsrFlash20250908:          consts.ModelFeatureMultimodal,
				consts.AliBLQwenAudioAsr:                   consts.ModelFeatureNone,
				consts.AliBLQwenAudioAsrLatest:             consts.ModelFeatureNone,
				consts.AliBLQwenAudioAsr20241204:           consts.ModelFeatureNone,
				consts.AliBLQwenMathPlus:                   consts.ModelFeatureNone,
				consts.AliBLQwenMathTurbo:                  consts.ModelFeatureNone,
				consts.AliBLQwen3CoderPlus:                 consts.ModelFeatureNone,
				consts.AliBLQwen3CoderPlus20250923:         consts.ModelFeatureNone,
				consts.AliBLQwen3CoderPlus20250722:         consts.ModelFeatureNone,
				consts.AliBLQwen3CoderFlash:                consts.ModelFeatureNone,
				consts.AliBLQwen3CoderFlash20250728:        consts.ModelFeatureNone,
				consts.AliBLQwenCoderPlus:                  consts.ModelFeatureNone,
				consts.AliBLQwenCoderPlusLatest:            consts.ModelFeatureNone,
				consts.AliBLQwenCoderPlus20241106:          consts.ModelFeatureNone,
				consts.AliBLQwenCoderTurbo:                 consts.ModelFeatureNone,
				consts.AliBLQwenCoderTurboLatest:           consts.ModelFeatureNone,
				consts.AliBLQwenCoderTurbo20240919:         consts.ModelFeatureNone,
				consts.AliBLQwenMtPlus:                     consts.ModelFeatureNone,
				consts.AliBLQwenMtTurbo:                    consts.ModelFeatureNone,
				consts.AliBLQwenDocTurbo:                   consts.ModelFeatureNone,
				consts.AliBLQwenDeepResearch:               consts.ModelFeatureNone,
				consts.AliBLQwen3Next80ba3bThinking:        consts.ModelFeatureNone,
				consts.AliBLQwen3Next80ba3bInstruct:        consts.ModelFeatureNone,
				consts.AliBLQwen3_235bA22bThinking2507:     consts.ModelFeatureNone,
				consts.AliBLQwen3_235bA22bInstruct2507:     consts.ModelFeatureNone,
				consts.AliBLQwen3_30bA3bThinking2507:       consts.ModelFeatureNone,
				consts.AliBLQwen3_30bA3bInstruct2507:       consts.ModelFeatureNone,
				consts.AliBLQwen3_235bA22b:                 consts.ModelFeatureNone,
				consts.AliBLQwen3_32b:                      consts.ModelFeatureNone,
				consts.AliBLQwen3_30bA3b:                   consts.ModelFeatureNone,
				consts.AliBLQwen3_14b:                      consts.ModelFeatureNone,
				consts.AliBLQwen3_8b:                       consts.ModelFeatureNone,
				consts.AliBLQwen3_4b:                       consts.ModelFeatureNone,
				consts.AliBLQwen3_17b:                      consts.ModelFeatureNone,
				consts.AliBLQwen3_06b:                      consts.ModelFeatureNone,
				consts.AliBLQwq32b:                         consts.ModelFeatureNone,
				consts.AliBLQwq32bPreview:                  consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_14bInstruct1m:        consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_7bInstruct1m:         consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_72bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_32bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_14bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_7bInstruct:           consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_3bInstruct:           consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_15bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5_05bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2_72bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwen2_57bA14bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2_7bInstruct:               consts.ModelFeatureNone,
				consts.AliBLQwen2_15bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwen2_05bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_110bChat:             consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_72bChat:              consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_32bChat:              consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_14bChat:              consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_7bChat:               consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_18bChat:              consts.ModelFeatureNone,
				consts.AliBLQwen1Dot5_05bChat:              consts.ModelFeatureNone,
				consts.AliBLQvq72bPreview:                  consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Omni7b:                consts.ModelFeatureNone,
				consts.AliBLQwen3Omni30bA3bCaptioner:       consts.ModelFeatureNone,
				consts.AliBLQwen3Vl30bA3bThinking:          consts.ModelFeatureNone,
				consts.AliBLQwen3Vl30bA3bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen3Vl8bThinking:              consts.ModelFeatureNone,
				consts.AliBLQwen3Vl8bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwen3Vl235bA22bThinking:        consts.ModelFeatureNone,
				consts.AliBLQwen3Vl235bA22bInstruct:        consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Vl72bInstruct:         consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Vl32bInstruct:         consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Vl7bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Vl3bInstruct:          consts.ModelFeatureNone,
				consts.AliBLQwen2Vl72bInstruct:             consts.ModelFeatureNone,
				consts.AliBLQwen2Vl7bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwen2Vl2bInstruct:              consts.ModelFeatureNone,
				consts.AliBLQwenVlV1:                       consts.ModelFeatureNone,
				consts.AliBLQwenVlChatV1:                   consts.ModelFeatureNone,
				consts.AliBLQwen2AudioInstruct:             consts.ModelFeatureNone,
				consts.AliBLQwenAudioChat:                  consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Math72bInstruct:       consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Math7bInstruct:        consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Math15bInstruct:       consts.ModelFeatureNone,
				consts.AliBLQwen3Coder480bA35bInstruct:     consts.ModelFeatureNone,
				consts.AliBLQwen3Coder30bA3bInstruct:       consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder32bInstruct:      consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder14bInstruct:      consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder7bInstruct:       consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder3bInstruct:       consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder15bInstruct:      consts.ModelFeatureNone,
				consts.AliBLQwen2Dot5Coder05bInstruct:      consts.ModelFeatureNone,
				consts.AliBLDeepSeekV3Dot2Exp:              consts.ModelFeatureNone,
				consts.AliBLDeepSeekV3Dot1:                 consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1:                     consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1_0528:                consts.ModelFeatureNone,
				consts.AliBLDeepSeekV3:                     consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillQwen15b:       consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillQwen7b:        consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillQwen14b:       consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillQwen32b:       consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillLlama8b:       consts.ModelFeatureNone,
				consts.AliBLDeepSeekR1DistillLlama70b:      consts.ModelFeatureNone,
				consts.AliBLMoonshotKimiK2Instruct:         consts.ModelFeatureNone,
				consts.AliBLGlm4Dot5:                       consts.ModelFeatureNone,
				consts.AliBLGlm4Dot5Air:                    consts.ModelFeatureNone,
				consts.AliBLLlama4Scout17b16eInstruct:      consts.ModelFeatureNone,
				consts.AliBLLlama4Maverick17b128eInstruct:  consts.ModelFeatureNone,
				consts.AliBLBaichuan2Turbo:                 consts.ModelFeatureNone,
				consts.AliBLAbab6Dot5gChat:                 consts.ModelFeatureNone,
				consts.AliBLAbab6Dot5tChat:                 consts.ModelFeatureNone,
				consts.AliBLAbab6Dot5sChat:                 consts.ModelFeatureNone,
			},
		},
	}
	core.RegisterProvider(consts.AliBL, aliblService)
}

// GetSupportedModels 获取支持的模型
func (s *aliblProvider) GetSupportedModels() (supportedModels map[consts.ModelType]map[string]consts.ModelFeature) {
	return s.supportedModels
}

// InitializeProviderConfig 初始化提供商配置
func (s *aliblProvider) InitializeProviderConfig(config *conf.ProviderConfig) {
	s.providerConfig = config
	s.lb = loadbalancer.NewLoadBalancer(s.providerConfig.APIKeys)
}
