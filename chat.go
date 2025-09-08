/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-16 14:31:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 22:29:12
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import (
	"context"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/core"
	"github.com/Mrzhouyl/go-aisdk/errors"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/models"
)

// CreateChatCompletion 创建聊天
func (c *SDKClient) CreateChatCompletion(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		// 判断是否流式传输
		if models.BoolValue(chatReq.Stream) {
			return nil, errors.ErrCompletionStreamNotSupported
		}
		// 创建聊天
		return ps.CreateChatCompletion(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider:  request.Provider,
		ModelType: consts.ChatModel,
		Model:     request.Model,
	}, request.UserInfo, "CreateChatCompletion", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ChatResponse)
	return
}

// CreateChatCompletionStream 创建流式聊天
func (c *SDKClient) CreateChatCompletionStream(ctx context.Context, request models.ChatRequest, opts ...httpclient.HTTPClientOption) (response models.ChatResponseStream, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, ps core.ProviderService, req any) (resp any, err error) {
		chatReq := req.(models.ChatRequest)
		chatReq.Stream = models.Bool(true)
		// 创建流式聊天
		return ps.CreateChatCompletionStream(ctx, chatReq, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, models.ModelInfo{
		Provider:  request.Provider,
		ModelType: consts.ChatModel,
		Model:     request.Model,
	}, request.UserInfo, "CreateChatCompletionStream", request, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ChatResponseStream)
	return
}
