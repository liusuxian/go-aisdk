/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:20:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-24 10:49:18
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/core"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/models"
)

// CreateImage 创建图像
func (c *SDKClient) CreateImage(ctx context.Context, request models.ImageRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		imageReq := req.(models.ImageRequest)
		// 创建图像
		return ps.CreateImage(ctx, imageReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider:  request.Provider,
		ModelType: consts.ImageModel,
		Model:     request.Model,
	}, request.UserInfo, "CreateImage", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ImageResponse)
	return
}

// CreateImageEdit 编辑图像
func (c *SDKClient) CreateImageEdit(ctx context.Context, request models.ImageEditRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		imageEditReq := req.(models.ImageEditRequest)
		// 编辑图像
		return ps.CreateImageEdit(ctx, imageEditReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider:  request.Provider,
		ModelType: consts.ImageModel,
		Model:     request.Model,
	}, request.UserInfo, "CreateImageEdit", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ImageResponse)
	return
}

// CreateImageVariation 变换图像
func (c *SDKClient) CreateImageVariation(ctx context.Context, request models.ImageVariationRequest, opts ...httpclient.HTTPClientOption) (response models.ImageResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		imageVariationReq := req.(models.ImageVariationRequest)
		// 变换图像
		return ps.CreateImageVariation(ctx, imageVariationReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider:  request.Provider,
		ModelType: consts.ImageModel,
		Model:     request.Model,
	}, request.UserInfo, "CreateImageVariation", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ImageResponse)
	return
}
