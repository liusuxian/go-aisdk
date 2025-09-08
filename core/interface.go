/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:45:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-09 15:49:29
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"context"

	"github.com/Mrzhouyl/go-aisdk/conf"
	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/models"
)

// ProviderService AI服务提供商的服务接口
type ProviderService interface {
	// 获取支持的模型
	// 返回值: map[模型类型]map[模型名称]模型特性位掩码
	GetSupportedModels() (supportedModels map[consts.ModelType]map[string]consts.ModelFeature)
	// 初始化提供商配置
	InitializeProviderConfig(config *conf.ProviderConfig)

	// 模型相关
	ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) // 列出模型
	// 聊天相关
	CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error)             // 创建聊天
	CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) // 创建流式聊天

	// TODO 图像相关
	CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error)                   // 创建图像
	CreateImageEdit(ctx context.Context, request models.ImageEditRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error)           // 编辑图像
	CreateImageVariation(ctx context.Context, request models.ImageVariationRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) // 变换图像

	// TODO 视频相关

	// TODO 音频相关
}
