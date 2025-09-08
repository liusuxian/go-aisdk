/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-02 04:49:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-17 18:56:58
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package aisdk

import "github.com/Mrzhouyl/go-aisdk/httpclient"

// WithMiddleware 添加中间件
func WithMiddleware(m httpclient.Middleware) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, m)
	}
}

// WithLogging 添加日志中间件
func WithLogging(config httpclient.LoggingMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, httpclient.NewLoggingMiddleware(config))
	}
}

// WithMetrics 添加监控中间件
func WithMetrics(config httpclient.MetricsMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, httpclient.NewMetricsMiddleware(config))
	}
}

// WithRetry 添加重试中间件
func WithRetry(config httpclient.RetryMiddlewareConfig) (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, httpclient.NewRetryMiddleware(config))
	}
}

// WithDefaultMiddlewares 添加默认中间件（日志、监控、重试）
func WithDefaultMiddlewares() (opt SDKClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares,
			httpclient.NewMetricsMiddleware(httpclient.DefaultMetricsConfig()),
			httpclient.NewRetryMiddleware(httpclient.DefaultRetryConfig()),
			httpclient.NewLoggingMiddleware(httpclient.DefaultLoggingConfig()),
		)
	}
}

// GetMetrics 获取指标数据（如果启用了监控中间件）
func (c *SDKClient) GetMetrics() (metrics map[string]any) {
	for _, mw := range c.middlewareChain.GetMiddlewares() {
		if metricsMiddleware, ok := mw.(*httpclient.MetricsMiddleware); ok {
			return metricsMiddleware.GetMetrics()
		}
	}
	return
}
