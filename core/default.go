/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:59:35
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 22:23:18
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"context"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/errors"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/models"
)

// DefaultProviderService 默认提供商服务实现
type DefaultProviderService struct {
}

// ListModels 列出模型
func (s *DefaultProviderService) ListModels(ctx context.Context, provider consts.Provider, opts ...httpclient.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = errors.WrapMethodNotSupportedByProvider(provider, "ListModels")
	return
}

// CreateChatCompletion 创建聊天
func (s *DefaultProviderService) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	err = errors.WrapMethodNotSupported(request.Provider, consts.ChatModel, request.Model, "CreateChatCompletion")
	return
}

// CreateChatCompletionStream 创建流式聊天
func (s *DefaultProviderService) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	err = errors.WrapMethodNotSupported(request.Provider, consts.ChatModel, request.Model, "CreateChatCompletionStream")
	return
}

// CreateImage 创建图像
func (s *DefaultProviderService) CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = errors.WrapMethodNotSupported(request.Provider, consts.ImageModel, request.Model, "CreateImage")
	return
}

// CreateImageEdit 编辑图像
func (s *DefaultProviderService) CreateImageEdit(ctx context.Context, request models.ImageEditRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = errors.WrapMethodNotSupported(request.Provider, consts.ImageModel, request.Model, "CreateImageEdit")
	return
}

// CreateImageVariation 变换图像
func (s *DefaultProviderService) CreateImageVariation(ctx context.Context, request models.ImageVariationRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	err = errors.WrapMethodNotSupported(request.Provider, consts.ImageModel, request.Model, "CreateImageVariation")
	return
}
