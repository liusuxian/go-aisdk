/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-20 01:15:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:36:16
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package common

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/internal/utils"
	"github.com/Mrzhouyl/go-aisdk/loadbalancer"
)

const (
	defaultHTTPClientTimeout           time.Duration = 10 * time.Second // 默认HTTP客户端请求超时时间
	defaultStreamReturnIntervalTimeout time.Duration = 15 * time.Second // 默认流式传输返回的间隔超时时间
	defaultEmptyMessagesLimit          uint          = 300              // 默认空消息限制
)

// ExecuteRequestContext 执行请求上下文
type ExecuteRequestContext struct {
	Provider    consts.Provider               // 提供商
	Method      string                        // http方法
	BaseURL     string                        // 基础URL
	ApiPath     string                        // 请求路径
	Opts        []httpclient.HTTPClientOption // 客户端选项
	LB          *loadbalancer.LoadBalancer    // 负载均衡器
	FormHandler httpclient.FormBuilderHandler // 构建表单请求体处理函数
	Response    httpclient.Response           // 响应数据
	ReqSetters  []httpclient.RequestOption    // 请求选项
}

// ExecuteRequest 执行请求
func ExecuteRequest(ctx context.Context, erc *ExecuteRequestContext) (err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClientWithConfig(httpclient.HTTPClientConfig{
		BaseURL:                     erc.BaseURL,
		HTTPClient:                  httpclient.NewDefaultHTTPDoer(defaultHTTPClientTimeout),
		ResponseDecoder:             utils.NewDeserializer(erc.Provider.String(), false),
		EmptyMessagesLimit:          defaultEmptyMessagesLimit,
		StreamReturnIntervalTimeout: defaultStreamReturnIntervalTimeout,
	})
	// 设置客户端选项
	for _, opt := range erc.Opts {
		opt(hc)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = erc.LB.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(erc.ReqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	// 构建表单请求体
	if erc.FormHandler != nil {
		var (
			formBody = &bytes.Buffer{}
			builder  = hc.GetFormBuilder(formBody)
		)
		if err = erc.FormHandler(builder); err != nil {
			return
		}
		setters = append(setters, httpclient.WithBody(formBody), httpclient.WithContentType(builder.FormDataContentType()))
	}
	if req, err = hc.NewRequest(ctx, erc.Method, hc.FullURL(erc.ApiPath), setters...); err != nil {
		return
	}
	// 发送请求
	err = hc.SendRequest(req, erc.Response)
	return
}

// ExecuteStreamRequest 执行流式传输请求
func ExecuteStreamRequest[T httpclient.Streamable](ctx context.Context, erc *ExecuteRequestContext) (stream *httpclient.StreamReader[T], err error) {
	// 新建 HTTP 客户端
	hc := httpclient.NewHTTPClientWithConfig(httpclient.HTTPClientConfig{
		BaseURL:                     erc.BaseURL,
		HTTPClient:                  httpclient.NewDefaultHTTPDoer(defaultHTTPClientTimeout),
		ResponseDecoder:             utils.NewDeserializer(erc.Provider.String(), true),
		EmptyMessagesLimit:          defaultEmptyMessagesLimit,
		StreamReturnIntervalTimeout: defaultStreamReturnIntervalTimeout,
	})
	// 设置客户端选项
	for _, opt := range erc.Opts {
		opt(hc)
	}
	// 获取一个APIKey
	var apiKey *loadbalancer.APIKey
	if apiKey, err = erc.LB.GetAPIKey(); err != nil {
		return
	}
	// 创建请求
	var (
		setters = append(erc.ReqSetters, httpclient.WithKeyValue("Authorization", fmt.Sprintf("Bearer %s", apiKey.Key)))
		req     *http.Request
	)
	if req, err = hc.NewRequest(ctx, erc.Method, hc.FullURL(erc.ApiPath), setters...); err != nil {
		return
	}
	// 发送流式请求
	return httpclient.SendRequestStream[T](hc, req)
}
