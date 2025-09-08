/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 19:11:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-25 10:55:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package conf_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Mrzhouyl/go-aisdk/conf"
	"github.com/Mrzhouyl/go-aisdk/consts"
)

func TestSDKConfigManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temporary test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")
	// 1. Test NewSDKConfigManager with non-existent config file
	_, err = conf.NewSDKConfigManager(configPath)
	if err == nil {
		t.Error("NewSDKConfigManager should return an error for non-existent config file")
	}
	// 2. Test NewSDKConfigManager with empty config path
	_, err = conf.NewSDKConfigManager("")
	if err == nil {
		t.Error("NewSDKConfigManager should return an error for empty config path")
	}
	// 3. Test with existing config file
	// Create a test config file with sample data
	testConfig := conf.SDKConfig{
		Providers: map[string]conf.ProviderConfig{
			"openai": {
				BaseURL:          "https://api.openai.com/v1",
				APIKeys:          []string{"test-key-1", "test-key-2"},
				OrgID:            "test-org-id",
				APIVersion:       "v1",
				AssistantVersion: "v2",
				Extra: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			"deepseek": {
				BaseURL: "https://api.deepseek.com",
				APIKeys: []string{"deepseek-key"},
			},
		},
	}
	// Write test config to file
	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	// 4. Test NewSDKConfigManager with existing config file
	manager2, err := conf.NewSDKConfigManager(configPath)
	if err != nil {
		t.Fatalf("NewSDKConfigManager with existing config failed: %v", err)
	}
	// 5. Test GetConfig
	loadedConfig := manager2.GetConfig()
	if len(loadedConfig.Providers) != 2 {
		t.Errorf("Expected 2 providers, got: %d", len(loadedConfig.Providers))
	}
	// 6. Test GetProviderConfig for OpenAI
	openaiConfig := manager2.GetProviderConfig(consts.OpenAI)
	// Verify OpenAI configuration
	expectedOpenAI := testConfig.Providers[consts.OpenAI.String()]
	if openaiConfig.BaseURL != expectedOpenAI.BaseURL {
		t.Errorf("OpenAI BaseURL mismatch, got: %v, want: %v", openaiConfig.BaseURL, expectedOpenAI.BaseURL)
	}
	if !reflect.DeepEqual(openaiConfig.APIKeys, expectedOpenAI.APIKeys) {
		t.Errorf("OpenAI APIKeys mismatch, got: %v, want: %v", openaiConfig.APIKeys, expectedOpenAI.APIKeys)
	}
	if openaiConfig.OrgID != expectedOpenAI.OrgID {
		t.Errorf("OpenAI OrgID mismatch, got: %v, want: %v", openaiConfig.OrgID, expectedOpenAI.OrgID)
	}
	if openaiConfig.APIVersion != expectedOpenAI.APIVersion {
		t.Errorf("OpenAI APIVersion mismatch, got: %v, want: %v", openaiConfig.APIVersion, expectedOpenAI.APIVersion)
	}
	if openaiConfig.AssistantVersion != expectedOpenAI.AssistantVersion {
		t.Errorf("OpenAI AssistantVersion mismatch, got: %v, want: %v", openaiConfig.AssistantVersion, expectedOpenAI.AssistantVersion)
	}
	if !reflect.DeepEqual(openaiConfig.Extra, expectedOpenAI.Extra) {
		t.Errorf("OpenAI Extra mismatch, got: %v, want: %v", openaiConfig.Extra, expectedOpenAI.Extra)
	}
	// 7. Test GetProviderConfig for DeepSeek
	deepseekConfig := manager2.GetProviderConfig(consts.DeepSeek)
	deepseekConfig2 := testConfig.Providers[consts.DeepSeek.String()]
	if deepseekConfig.BaseURL != deepseekConfig2.BaseURL {
		t.Errorf("Anthropic BaseURL mismatch, got: %v, want: %v", deepseekConfig.BaseURL, deepseekConfig2.BaseURL)
	}
	if !reflect.DeepEqual(deepseekConfig.APIKeys, deepseekConfig2.APIKeys) {
		t.Errorf("Anthropic APIKeys mismatch, got: %v, want: %v", deepseekConfig.APIKeys, deepseekConfig2.APIKeys)
	}
	// 8. Test deep copy functionality
	// Modify the returned config and verify it doesn't affect the internal config
	openaiConfig.BaseURL = "modified-url"
	openaiConfig.APIKeys[0] = "modified-key"
	openaiConfig.Extra["key1"] = "modified-value"
	// Get the config again and verify it's unchanged
	openaiConfig2 := manager2.GetProviderConfig(consts.OpenAI)
	if openaiConfig2.BaseURL != expectedOpenAI.BaseURL {
		t.Error("GetProviderConfig should return a deep copy, BaseURL should not be affected by modifications")
	}
	if openaiConfig2.APIKeys[0] != expectedOpenAI.APIKeys[0] {
		t.Error("GetProviderConfig should return a deep copy, APIKeys should not be affected by modifications")
	}
	if openaiConfig2.Extra["key1"] != expectedOpenAI.Extra["key1"] {
		t.Error("GetProviderConfig should return a deep copy, Extra should not be affected by modifications")
	}
	// 9. Test GetConfig deep copy functionality
	fullConfig := manager2.GetConfig()
	// Modify the returned config
	if providerConfig, ok := fullConfig.Providers[consts.OpenAI.String()]; ok {
		providerConfig.BaseURL = "modified-full-config-url"
		fullConfig.Providers[consts.OpenAI.String()] = providerConfig
	}
	// Verify internal config is unchanged
	openaiConfig3 := manager2.GetProviderConfig(consts.OpenAI)
	if openaiConfig3.BaseURL != expectedOpenAI.BaseURL {
		t.Error("GetConfig should return a deep copy, modifications should not affect internal config")
	}
	// 10. Test Load method with invalid JSON
	invalidConfigPath := filepath.Join(tempDir, "invalid-config.json")
	if err := os.WriteFile(invalidConfigPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}
	_, err = conf.NewSDKConfigManager(invalidConfigPath)
	if err == nil {
		t.Error("NewSDKConfigManager should return an error for invalid JSON config")
	}
	// 11. Test GetProviderConfig for AliBL provider
	providerConfig := manager2.GetProviderConfig(consts.AliBL)
	if providerConfig.BaseURL != "" {
		t.Error("GetProviderConfig should return an empty config for AliBL provider")
	}
}
