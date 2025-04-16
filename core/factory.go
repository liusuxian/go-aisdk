/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 18:45:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-16 10:47:58
 * @Description: 提供AI服务的核心功能，包括提供商工厂和相关接口，采用单例模式实现，通过包级函数直接访问功能
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package core

import (
	"fmt"
	"github.com/liusuxian/aisdk/consts"
	"sync"
)

// providerFactory 管理所有AI服务提供商的工厂
type providerFactory struct {
	providers map[consts.Provider]ProviderService // 所有提供商
	mu        sync.RWMutex                        // 读写锁
}

var (
	factory *providerFactory // 全局工厂单例实例(非导出)
)

// init 包初始化时创建 providerFactory 单例
// 确保在导入包时自动初始化工厂
func init() {
	factory = &providerFactory{
		providers: make(map[consts.Provider]ProviderService),
	}
}

// RegisterProvider 注册提供商
func RegisterProvider(provider consts.Provider, service ProviderService) {
	factory.mu.Lock()
	defer factory.mu.Unlock()

	factory.providers[provider] = service
}

// GetProvider 获取提供商
func GetProvider(provider consts.Provider) (service ProviderService, err error) {
	factory.mu.RLock()
	defer factory.mu.RUnlock()

	if p, ok := factory.providers[provider]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider %s not registered", provider)
}

// ListProviders 列出所有注册的提供商
func ListProviders() (providers []consts.Provider) {
	factory.mu.RLock()
	defer factory.mu.RUnlock()

	providers = make([]consts.Provider, 0, len(factory.providers))
	for p := range factory.providers {
		providers = append(providers, p)
	}
	return
}
