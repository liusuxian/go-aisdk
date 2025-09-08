/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-07-07 21:01:34
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 23:16:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package errors

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Mrzhouyl/go-aisdk/httpclient"
)

var (
	ErrFailedToCreateConfigManager  = errors.New("failed to create config manager")                                                    // 创建配置管理器失败
	ErrFailedToCreateFlakeInstance  = errors.New("failed to create flake instance")                                                    // 创建分布式唯一ID生成器失败
	ErrProviderNotSupported         = errors.New("provider is not supported")                                                          // 提供商不支持
	ErrModelTypeNotSupported        = errors.New("model type is not supported")                                                        // 模型类型不支持
	ErrModelNotSupported            = errors.New("model is not supported")                                                             // 模型不支持
	ErrMethodNotSupported           = errors.New("method is not supported")                                                            // 方法不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrTooManyEmptyStreamMessages   = httpclient.ErrTooManyEmptyStreamMessages                                                         // 流式传输发送了太多空消息
	ErrStreamReturnIntervalTimeout  = httpclient.ErrStreamReturnIntervalTimeout                                                        // 流式传输返回间隔超时
)

// WrapFailedToCreateConfigManager 包装创建配置管理器失败错误
func WrapFailedToCreateConfigManager(text string) (err error) {
	return fmt.Errorf("%s: %w", text, ErrFailedToCreateConfigManager)
}

// WrapFailedToCreateFlakeInstance 包装创建分布式唯一ID生成器失败错误
func WrapFailedToCreateFlakeInstance(text string) (err error) {
	return fmt.Errorf("%s: %w", text, ErrFailedToCreateFlakeInstance)
}

// WrapProviderNotSupported 包装提供商不支持错误
func WrapProviderNotSupported(provider fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] is not supported: %w", provider.String(), ErrProviderNotSupported)
}

// WrapModelTypeNotSupported 包装模型类型不支持错误
func WrapModelTypeNotSupported(provider, modelType fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] does not support model type [%s]: %w",
		provider.String(), modelType.String(), ErrModelTypeNotSupported)
}

// WrapModelNotSupported 包装模型不支持错误
func WrapModelNotSupported(provider fmt.Stringer, model string, modelType fmt.Stringer) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s]: %w",
		provider.String(), model, modelType.String(), ErrModelNotSupported)
}

// WrapMethodNotSupported 包装方法不支持错误
func WrapMethodNotSupported(provider, modelType fmt.Stringer, model, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] model of type [%s] with method [%s]: %w",
		provider.String(), model, modelType.String(), method, ErrMethodNotSupported)
}

// WrapMethodNotSupportedByProvider 包装方法不支持错误
func WrapMethodNotSupportedByProvider(provider fmt.Stringer, method string) (err error) {
	return fmt.Errorf("provider [%s] does not support the [%s] method: %w", provider.String(), method, ErrMethodNotSupported)
}

// IsFailedToCreateConfigManagerError 判断是否是创建配置管理器失败错误
func IsFailedToCreateConfigManagerError(err error) (is bool) {
	return errors.Is(err, ErrFailedToCreateConfigManager)
}

// IsFailedToCreateFlakeInstanceError 判断是否是创建分布式唯一ID生成器失败错误
func IsFailedToCreateFlakeInstanceError(err error) (is bool) {
	return errors.Is(err, ErrFailedToCreateFlakeInstance)
}

// IsProviderNotSupportedError 判断是否是提供商不支持错误
func IsProviderNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrProviderNotSupported)
}

// IsModelTypeNotSupportedError 判断是否是模型类型不支持错误
func IsModelTypeNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrModelTypeNotSupported)
}

// IsModelNotSupportedError 判断是否是模型不支持错误
func IsModelNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrModelNotSupported)
}

// IsMethodNotSupportedError 判断是否是方法不支持错误
func IsMethodNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrMethodNotSupported)
}

// IsCompletionStreamNotSupportedError 判断是否是流式传输不支持错误
func IsCompletionStreamNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrCompletionStreamNotSupported)
}

// IsTooManyEmptyStreamMessagesError 判断是否是流式传输发送了太多空消息错误
func IsTooManyEmptyStreamMessagesError(err error) (is bool) {
	return errors.Is(err, ErrTooManyEmptyStreamMessages)
}

// IsStreamReturnIntervalTimeoutError 判断是否是流式传输返回间隔超时错误
func IsStreamReturnIntervalTimeoutError(err error) (is bool) {
	return errors.Is(err, ErrStreamReturnIntervalTimeout)
}

// IsCanceledError 判断是否是取消错误
func IsCanceledError(err error) (is bool) {
	return errors.Is(err, context.Canceled)
}

// IsDeadlineExceededError 判断是否是截止时间错误
func IsDeadlineExceededError(err error) (is bool) {
	return errors.Is(err, context.DeadlineExceeded)
}

// IsNetError 判断是否是网络错误
func IsNetError(err error) (is bool) {
	var netErr net.Error
	return errors.As(err, &netErr)
}

// SDKError SDK错误
type SDKError struct {
	RequestID string // 请求ID
	Err       error  // 原始错误
}

// Error 错误信息
func (e *SDKError) Error() (errStr string) {
	return fmt.Sprintf("request_id: %s, error: %v", e.RequestID, e.Err)
}

// RequestID 获取请求ID
func RequestID(err error) (requestId string) {
	if err == nil {
		return ""
	}

	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.RequestID
	}

	return ""
}

// Unwrap 解包错误
func Unwrap(err error) (originalError error) {
	if err == nil {
		return nil
	}
	// 解包 SDKError
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		if sdkErr.Err != nil {
			return sdkErr.Err
		}
		return err // 如果内部错误为 nil，返回 SDKError 本身
	}
	// 解包 RequestError
	var requestError *httpclient.RequestError
	if errors.As(err, &requestError) {
		if requestError.Err != nil {
			return requestError.Err
		}
		return err // 如果内部错误为 nil，返回 RequestError 本身
	}
	// 其他类型的错误
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err // 已经是最底层错误，返回原错误
	}
	return unwrapped
}

// Cause 错误根因
func Cause(err error) (causeError error) {
	return doCause(err)
}

// doCause 递归获取错误根因
func doCause(err error) (causeError error) {
	if err == nil {
		return nil
	}
	// 解包错误
	unwrapped := Unwrap(err)
	if unwrapped == nil {
		return err // 已经到达最底层错误，返回当前错误
	}
	// 防止无限递归：如果解包后的错误与原错误相同，直接返回
	if unwrapped == err {
		return err
	}
	return doCause(unwrapped)
}
