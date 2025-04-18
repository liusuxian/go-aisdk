/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 19:09:15
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-16 17:43:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package conf

import (
	"encoding/json"
	"fmt"
	"github.com/liusuxian/aisdk/consts"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

// ProviderConfig AI服务提供商的配置
type ProviderConfig struct {
	BaseURL          string            `json:"base_url"`          // 基础URL，用于自定义API服务器的地址
	APIKeys          []string          `json:"api_keys"`          // API密钥列表
	OrgID            string            `json:"org_id"`            // 组织ID，对于某些提供商可能需要
	APIVersion       string            `json:"api_version"`       // API版本，对于某些提供商可能需要
	AssistantVersion string            `json:"assistant_version"` // 助手版本，对于某些提供商可能需要
	Extra            map[string]string `json:"extra"`             // 额外参数，对于某些提供商可能需要
}

// ConnectionOptions 连接选项
type ConnectionOptions struct {
	RequestTimeout              time.Duration   `json:"request_timeout"`                // 请求超时时间
	StreamReturnIntervalTimeout time.Duration   `json:"stream_return_interval_timeout"` // 流式返回间隔超时时间
	MaxRetries                  uint            `json:"max_retries"`                    // 最大重试次数。当配置了`retry_delay_list`时，该参数将失效
	RetryDelay                  time.Duration   `json:"retry_delay"`                    // 重试延迟时间。当配置了`retry_delay_list`时，该参数将失效
	RetryIncreaseDelay          bool            `json:"retry_increase_delay"`           // 是否让延迟时间随着重试次数增加而线性增加。当配置了`retry_delay_list`时，该参数将失效
	RetryDelayList              []time.Duration `json:"retry_delay_list"`               // 自定义延迟列表
}

// MarshalJSON 自定义 JSON 编码
func (o ConnectionOptions) MarshalJSON() (b []byte, err error) {
	type Alias ConnectionOptions
	return json.Marshal(&struct {
		RequestTimeout              string   `json:"request_timeout"`
		StreamReturnIntervalTimeout string   `json:"stream_return_interval_timeout"`
		RetryDelay                  string   `json:"retry_delay"`
		RetryDelayList              []string `json:"retry_delay_list"`
		Alias
	}{
		RequestTimeout:              o.RequestTimeout.String(),
		StreamReturnIntervalTimeout: o.StreamReturnIntervalTimeout.String(),
		RetryDelay:                  o.RetryDelay.String(),
		RetryDelayList:              durationSliceToStringSlice(o.RetryDelayList),
		Alias:                       Alias(o),
	})
}

// UnmarshalJSON 自定义 JSON 解码
func (o *ConnectionOptions) UnmarshalJSON(data []byte) (err error) {
	type Alias ConnectionOptions
	aux := &struct {
		RequestTimeout              string   `json:"request_timeout"`
		StreamReturnIntervalTimeout string   `json:"stream_return_interval_timeout"`
		RetryDelay                  string   `json:"retry_delay"`
		RetryDelayList              []string `json:"retry_delay_list"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	if aux.RequestTimeout != "" {
		if o.RequestTimeout, err = time.ParseDuration(aux.RequestTimeout); err != nil {
			return
		}
	}
	if aux.StreamReturnIntervalTimeout != "" {
		if o.StreamReturnIntervalTimeout, err = time.ParseDuration(aux.StreamReturnIntervalTimeout); err != nil {
			return
		}
	}
	if aux.RetryDelay != "" {
		if o.RetryDelay, err = time.ParseDuration(aux.RetryDelay); err != nil {
			return
		}
	}
	if len(aux.RetryDelayList) > 0 {
		if o.RetryDelayList, err = stringSliceToDurationSlice(aux.RetryDelayList); err != nil {
			return
		}
	}
	return
}

// SDKConfig SDK整体配置
type SDKConfig struct {
	Providers         map[consts.Provider]ProviderConfig `json:"providers"`          // AI服务提供商的配置
	ConnectionOptions ConnectionOptions                  `json:"connection_options"` // 连接选项
}

// SDKConfigManager SDK配置管理器
type SDKConfigManager struct {
	configPath string       // 配置文件路径
	config     SDKConfig    // SDK配置
	mu         sync.RWMutex // 读写锁
}

// NewSDKConfigManager 创建SDK配置管理器
func NewSDKConfigManager(configPath string) (manager *SDKConfigManager, err error) {
	if configPath == "" {
		var homeDir string
		if homeDir, err = os.UserHomeDir(); err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".go-aisdk", "config.json")
	}

	manager = &SDKConfigManager{
		configPath: configPath,
		config: SDKConfig{
			Providers: make(map[consts.Provider]ProviderConfig),
			ConnectionOptions: ConnectionOptions{
				RequestTimeout:              10 * time.Second,
				StreamReturnIntervalTimeout: 20 * time.Second,
				MaxRetries:                  3,
				RetryDelay:                  1 * time.Second,
				RetryIncreaseDelay:          false,
				RetryDelayList:              []time.Duration{},
			},
		},
	}
	// 确保目录存在
	configDir := filepath.Dir(configPath)
	if err = os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	// 尝试加载配置
	if _, err = os.Stat(configPath); err == nil {
		if err = manager.Load(); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}
	return
}

// Load 从文件加载配置
func (m *SDKConfigManager) Load() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var data []byte
	if data, err = os.ReadFile(m.configPath); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err = json.Unmarshal(data, &m.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return
}

// Save 保存配置到文件
func (m *SDKConfigManager) Save() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var data []byte
	if data, err = json.MarshalIndent(m.config, "", "  "); err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err = os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return
}

// GetConfig 获取整个配置
func (m *SDKConfigManager) GetConfig() (configCopy SDKConfig) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// 返回配置的副本，防止外部修改
	configCopy = SDKConfig{
		Providers:         make(map[consts.Provider]ProviderConfig),
		ConnectionOptions: cloneConnectionOptions(m.config.ConnectionOptions),
	}

	for k, v := range m.config.Providers {
		configCopy.Providers[k] = cloneProviderConfig(v)
	}
	return
}

// SetProviderConfig 设置提供商配置
func (m *SDKConfigManager) SetProviderConfig(provider consts.Provider, config ProviderConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Providers[provider] = cloneProviderConfig(config)
}

// GetProviderConfig 获取提供商配置
func (m *SDKConfigManager) GetProviderConfig(provider consts.Provider) (config ProviderConfig, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if cfg, ok := m.config.Providers[provider]; ok {
		config = cloneProviderConfig(cfg)
		return
	}

	return ProviderConfig{}, fmt.Errorf("provider %s not configured", provider)
}

// SetConnectionOptions 设置连接选项
func (m *SDKConfigManager) SetConnectionOptions(options ConnectionOptions) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.ConnectionOptions = cloneConnectionOptions(options)
}

// GetConnectionOptions 获取连接选项
func (m *SDKConfigManager) GetConnectionOptions() (options ConnectionOptions) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return cloneConnectionOptions(m.config.ConnectionOptions)
}

// cloneProviderConfig 深拷贝 ProviderConfig
func cloneProviderConfig(source ProviderConfig) (dest ProviderConfig) {
	extraCopy := make(map[string]string)
	maps.Copy(extraCopy, source.Extra)

	dest = ProviderConfig{
		BaseURL:          source.BaseURL,
		APIKeys:          slices.Clone(source.APIKeys),
		OrgID:            source.OrgID,
		APIVersion:       source.APIVersion,
		AssistantVersion: source.AssistantVersion,
		Extra:            extraCopy,
	}
	return
}

// cloneConnectionOptions 深拷贝 ConnectionOptions
func cloneConnectionOptions(source ConnectionOptions) (dest ConnectionOptions) {
	return ConnectionOptions{
		RequestTimeout:              source.RequestTimeout,
		StreamReturnIntervalTimeout: source.StreamReturnIntervalTimeout,
		MaxRetries:                  source.MaxRetries,
		RetryDelay:                  source.RetryDelay,
		RetryIncreaseDelay:          source.RetryIncreaseDelay,
		RetryDelayList:              slices.Clone(source.RetryDelayList),
	}
}

// durationSliceToStringSlice 将时间间隔列表转换为字符串列表
func durationSliceToStringSlice(durations []time.Duration) (list []string) {
	list = make([]string, len(durations))
	for i, d := range durations {
		list[i] = d.String()
	}
	return
}

// stringSliceToDurationSlice 将字符串列表转换为时间间隔列表
func stringSliceToDurationSlice(strs []string) (list []time.Duration, err error) {
	list = make([]time.Duration, len(strs))
	for i, s := range strs {
		var d time.Duration
		if d, err = time.ParseDuration(s); err != nil {
			return
		}
		list[i] = d
	}
	return
}
