/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 12:31:10
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 20:57:51
 * @Description: AliBL服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package alibl

import (
	"github.com/Mrzhouyl/go-aisdk/conf"
	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/core"
	"github.com/Mrzhouyl/go-aisdk/loadbalancer"
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
				consts.AliBLQwqPlus:                       consts.ModelFeature(0),
				consts.AliBLQwqPlusLatest:                 consts.ModelFeature(0),
				consts.AliBLQwqPlus20250305:               consts.ModelFeature(0),
				consts.AliBLQwenMax:                       consts.ModelFeature(0),
				consts.AliBLQwenMaxLatest:                 consts.ModelFeature(0),
				consts.AliBLQwenMax20250125:               consts.ModelFeature(0),
				consts.AliBLQwenMax20240919:               consts.ModelFeature(0),
				consts.AliBLQwenMax20240428:               consts.ModelFeature(0),
				consts.AliBLQwenMax20240403:               consts.ModelFeature(0),
				consts.AliBLQwenPlus:                      consts.ModelFeature(0),
				consts.AliBLQwenPlusLatest:                consts.ModelFeature(0),
				consts.AliBLQwenPlus20250428:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20250125:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20250112:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20241220:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20241127:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20241125:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20240919:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20240806:              consts.ModelFeature(0),
				consts.AliBLQwenPlus20240723:              consts.ModelFeature(0),
				consts.AliBLQwenTurbo:                     consts.ModelFeature(0),
				consts.AliBLQwenTurboLatest:               consts.ModelFeature(0),
				consts.AliBLQwenTurbo20250428:             consts.ModelFeature(0),
				consts.AliBLQwenTurbo20250211:             consts.ModelFeature(0),
				consts.AliBLQwenTurbo20241101:             consts.ModelFeature(0),
				consts.AliBLQwenTurbo20240919:             consts.ModelFeature(0),
				consts.AliBLQwenTurbo20240624:             consts.ModelFeature(0),
				consts.AliBLQwenLong:                      consts.ModelFeature(0),
				consts.AliBLQwenLongLatest:                consts.ModelFeature(0),
				consts.AliBLQwenLong20250125:              consts.ModelFeature(0),
				consts.AliBLQwenOmniTurbo:                 consts.ModelFeature(1),
				consts.AliBLQwenOmniTurboLatest:           consts.ModelFeature(1),
				consts.AliBLQwenOmniTurbo20250326:         consts.ModelFeature(1),
				consts.AliBLQwenOmniTurbo20250119:         consts.ModelFeature(1),
				consts.AliBLQwenOmniTurboRealtime:         consts.ModelFeature(1),
				consts.AliBLQwenOmniTurboRealtimeLatest:   consts.ModelFeature(1),
				consts.AliBLQwenOmniTurboRealtime20250508: consts.ModelFeature(1),
				consts.AliBLQvqMax:                        consts.ModelFeature(1),
				consts.AliBLQvqMaxLatest:                  consts.ModelFeature(1),
				consts.AliBLQvqMax20250515:                consts.ModelFeature(1),
				consts.AliBLQvqMax20250325:                consts.ModelFeature(1),
				consts.AliBLQvqPlus:                       consts.ModelFeature(1),
				consts.AliBLQvqPlusLatest:                 consts.ModelFeature(1),
				consts.AliBLQvqPlus20250515:               consts.ModelFeature(1),
				consts.AliBLQwenVlMax:                     consts.ModelFeature(1),
				consts.AliBLQwenVlMaxLatest:               consts.ModelFeature(1),
				consts.AliBLQwenVlMax20250408:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20250402:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20250125:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20241230:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20241119:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20241030:             consts.ModelFeature(1),
				consts.AliBLQwenVlMax20240809:             consts.ModelFeature(1),
				consts.AliBLQwenVlPlus:                    consts.ModelFeature(1),
				consts.AliBLQwenVlPlusLatest:              consts.ModelFeature(1),
				consts.AliBLQwenVlPlus20250507:            consts.ModelFeature(1),
				consts.AliBLQwenVlPlus20250125:            consts.ModelFeature(1),
				consts.AliBLQwenVlPlus20250102:            consts.ModelFeature(1),
				consts.AliBLQwenVlPlus20240809:            consts.ModelFeature(1),
				consts.AliBLQwenVlPlus20231201:            consts.ModelFeature(1),
				consts.AliBLQwenVlOcr:                     consts.ModelFeature(1),
				consts.AliBLQwenVlOcrLatest:               consts.ModelFeature(1),
				consts.AliBLQwenVlOcr20250413:             consts.ModelFeature(1),
				consts.AliBLQwenVlOcr20241028:             consts.ModelFeature(1),
				consts.AliBLQwenAudioTurbo:                consts.ModelFeature(1),
				consts.AliBLQwenAudioTurboLatest:          consts.ModelFeature(1),
				consts.AliBLQwenAudioTurbo20241204:        consts.ModelFeature(1),
				consts.AliBLQwenAudioTurbo20240807:        consts.ModelFeature(1),
				consts.AliBLQwenAudioAsr:                  consts.ModelFeature(1),
				consts.AliBLQwenAudioAsrLatest:            consts.ModelFeature(1),
				consts.AliBLQwenAudioAsr20241204:          consts.ModelFeature(1),
				consts.AliBLQwenMathPlus:                  consts.ModelFeature(0),
				consts.AliBLQwenMathPlusLatest:            consts.ModelFeature(0),
				consts.AliBLQwenMathPlus20240919:          consts.ModelFeature(0),
				consts.AliBLQwenMathPlus20240816:          consts.ModelFeature(0),
				consts.AliBLQwenMathTurbo:                 consts.ModelFeature(0),
				consts.AliBLQwenMathTurboLatest:           consts.ModelFeature(0),
				consts.AliBLQwenMathTurbo20240919:         consts.ModelFeature(0),
				consts.AliBLQwenCoderPlus:                 consts.ModelFeature(0),
				consts.AliBLQwenCoderPlusLatest:           consts.ModelFeature(0),
				consts.AliBLQwenCoderPlus20241106:         consts.ModelFeature(0),
				consts.AliBLQwenCoderTurbo:                consts.ModelFeature(0),
				consts.AliBLQwenCoderTurboLatest:          consts.ModelFeature(0),
				consts.AliBLQwenCoderTurbo20240919:        consts.ModelFeature(0),
				consts.AliBLQwenMtPlus:                    consts.ModelFeature(0),
				consts.AliBLQwenMtTurbo:                   consts.ModelFeature(0),
				consts.AliBLQwen3_235bA22b:                consts.ModelFeature(0),
				consts.AliBLQwen3_32b:                     consts.ModelFeature(0),
				consts.AliBLQwen3_30bA3b:                  consts.ModelFeature(0),
				consts.AliBLQwen3_14b:                     consts.ModelFeature(0),
				consts.AliBLQwen3_8b:                      consts.ModelFeature(0),
				consts.AliBLQwen3_4b:                      consts.ModelFeature(0),
				consts.AliBLQwen3_17b:                     consts.ModelFeature(0),
				consts.AliBLQwen3_06b:                     consts.ModelFeature(0),
				consts.AliBLQwq32b:                        consts.ModelFeature(0),
				consts.AliBLQwq32bPreview:                 consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_14bInstruct1m:       consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_7bInstruct1m:        consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_72bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_32bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_14bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_7bInstruct:          consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_3bInstruct:          consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_15bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2Dot5_05bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2_72bInstruct:             consts.ModelFeature(0),
				consts.AliBLQwen2_57bA14bInstruct:         consts.ModelFeature(0),
				consts.AliBLQwen2_7bInstruct:              consts.ModelFeature(0),
				consts.AliBLQwen2_15bInstruct:             consts.ModelFeature(0),
				consts.AliBLQwen2_05bInstruct:             consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_110bChat:            consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_72bChat:             consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_32bChat:             consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_14bChat:             consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_7bChat:              consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_18bChat:             consts.ModelFeature(0),
				consts.AliBLQwen1Dot5_05bChat:             consts.ModelFeature(0),
				consts.AliBLQvq72bPreview:                 consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Omni7b:               consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Vl72bInstruct:        consts.ModelFeature(1),
				consts.AliBLQwen2Dot5Vl32bInstruct:        consts.ModelFeature(1),
				consts.AliBLQwen2Dot5Vl7bInstruct:         consts.ModelFeature(1),
				consts.AliBLQwen2Dot5Vl3bInstruct:         consts.ModelFeature(1),
				consts.AliBLQwen2Vl72bInstruct:            consts.ModelFeature(1),
				consts.AliBLQwen2Vl7bInstruct:             consts.ModelFeature(1),
				consts.AliBLQwen2Vl2bInstruct:             consts.ModelFeature(1),
				consts.AliBLQwenVlV1:                      consts.ModelFeature(1),
				consts.AliBLQwenVlChatV1:                  consts.ModelFeature(1),
				consts.AliBLQwen2AudioInstruct:            consts.ModelFeature(1),
				consts.AliBLQwenAudioChat:                 consts.ModelFeature(1),
				consts.AliBLQwen2Dot5Math72bInstruct:      consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Math7bInstruct:       consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Math15bInstruct:      consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder32bInstruct:     consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder14bInstruct:     consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder7bInstruct:      consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder3bInstruct:      consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder15bInstruct:     consts.ModelFeature(0),
				consts.AliBLQwen2Dot5Coder05bInstruct:     consts.ModelFeature(0),
				consts.AliBLDeepSeekR1:                    consts.ModelFeature(0),
				consts.AliBLDeepSeekR1_0528:               consts.ModelFeature(0),
				consts.AliBLDeepSeekV3:                    consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillQwen15b:      consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillQwen7b:       consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillQwen14b:      consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillQwen32b:      consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillLlama8b:      consts.ModelFeature(0),
				consts.AliBLDeepSeekR1DistillLlama70b:     consts.ModelFeature(0),
				consts.AliBLLlama3Dot3_70bInstruct:        consts.ModelFeature(0),
				consts.AliBLLlama3Dot2_3bInstruct:         consts.ModelFeature(0),
				consts.AliBLLlama3Dot2_1bInstruct:         consts.ModelFeature(0),
				consts.AliBLLlama3Dot1_405bInstruct:       consts.ModelFeature(0),
				consts.AliBLLlama3Dot1_70bInstruct:        consts.ModelFeature(0),
				consts.AliBLLlama3Dot1_8bInstruct:         consts.ModelFeature(0),
				consts.AliBLLlama3_70bInstruct:            consts.ModelFeature(0),
				consts.AliBLLlama3_8bInstruct:             consts.ModelFeature(0),
				consts.AliBLLlama2_13bChatV2:              consts.ModelFeature(0),
				consts.AliBLLlama2_7bChatV2:               consts.ModelFeature(0),
				consts.AliBLLlama4Scout17b16eInstruct:     consts.ModelFeature(0),
				consts.AliBLLlama4Maverick17b128eInstruct: consts.ModelFeature(0),
				consts.AliBLLlama3Dot2_90bVisionInstruct:  consts.ModelFeature(0),
				consts.AliBLLlama3Dot2_11bVision:          consts.ModelFeature(0),
				consts.AliBLBaichuan2Turbo:                consts.ModelFeature(0),
				consts.AliBLBaichuan2_13bChatV1:           consts.ModelFeature(0),
				consts.AliBLBaichuan2_7bChatV1:            consts.ModelFeature(0),
				consts.AliBLBaichuan7bV1:                  consts.ModelFeature(0),
				consts.AliBLChatglm3_6b:                   consts.ModelFeature(0),
				consts.AliBLChatglm6bV2:                   consts.ModelFeature(0),
				consts.AliBLYiLarge:                       consts.ModelFeature(0),
				consts.AliBLYiMedium:                      consts.ModelFeature(0),
				consts.AliBLYiLargeRag:                    consts.ModelFeature(0),
				consts.AliBLYiLargeTurbo:                  consts.ModelFeature(0),
				consts.AliBLAbab6Dot5gChat:                consts.ModelFeature(0),
				consts.AliBLAbab6Dot5tChat:                consts.ModelFeature(0),
				consts.AliBLAbab6Dot5sChat:                consts.ModelFeature(0),
				consts.AliBLZiyaLlama13bV1:                consts.ModelFeature(0),
				consts.AliBLBelleLlama13b2mV1:             consts.ModelFeature(0),
				consts.AliBLChatyuanLargeV2:               consts.ModelFeature(0),
				consts.AliBLBilla7bSftV1:                  consts.ModelFeature(0),
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
