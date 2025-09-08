/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:15:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 23:39:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mrzhouyl/go-aisdk"
	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/errors"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/models"
)

func getApiKeys(envKey string) (apiKeys string) {
	raw := strings.TrimSpace(os.Getenv(envKey))
	if raw == "" {
		log.Fatal("DEEPSEEK_API_KEYS is empty. Set it first, e.g. PowerShell: $env:DEEPSEEK_API_KEYS=\"sk-xxxxxxxx\"")
	}
	parts := strings.Split(raw, ",")
	quoted := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			quoted = append(quoted, fmt.Sprintf(`"%s"`, v))
		}
	}
	if len(quoted) == 0 {
		log.Fatal("DEEPSEEK_API_KEYS contains no valid keys after filtering empties. Example: sk-xxxxxxxx or sk-xxx1,sk-xxx2")
	}
	return strings.Join(quoted, ",")
}

func isError(err error) {
	if err != nil {
		originalErr := errors.Unwrap(err)
		fmt.Println("originalErr =", originalErr)
		fmt.Println("Cause Error =", errors.Cause(err))
		switch {
		case errors.IsFailedToCreateConfigManagerError(originalErr):
			fmt.Println("IsFailedToCreateConfigManagerError =", true)
		case errors.IsFailedToCreateFlakeInstanceError(originalErr):
			fmt.Println("IsFailedToCreateFlakeInstanceError =", true)
		case errors.IsProviderNotSupportedError(originalErr):
			fmt.Println("IsProviderNotSupportedError =", true)
		case errors.IsModelTypeNotSupportedError(originalErr):
			fmt.Println("IsModelTypeNotSupportedError =", true)
		case errors.IsModelNotSupportedError(originalErr):
			fmt.Println("IsModelNotSupportedError =", true)
		case errors.IsMethodNotSupportedError(originalErr):
			fmt.Println("IsMethodNotSupportedError =", true)
		case errors.IsCompletionStreamNotSupportedError(originalErr):
			fmt.Println("IsCompletionStreamNotSupportedError =", true)
		case errors.IsTooManyEmptyStreamMessagesError(originalErr):
			fmt.Println("IsTooManyEmptyStreamMessagesError =", true)
		case errors.IsStreamReturnIntervalTimeoutError(originalErr):
			fmt.Println("IsStreamReturnIntervalTimeoutError =", true)
		case errors.IsCanceledError(originalErr):
			fmt.Println("IsCanceledError =", true)
		case errors.IsDeadlineExceededError(originalErr):
			fmt.Println("IsDeadlineExceededError =", true)
		case errors.IsNetError(originalErr):
			fmt.Println("IsNetError =", true)
		default:
			fmt.Println("unknown error =", err)
		}
	}
}

func listModels(ctx context.Context, client *aisdk.SDKClient) (response models.ListModelsResponse, err error) {
	return client.ListModels(ctx, models.ListModelsRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.DeepSeek,
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	// 定义工具函数
	params := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，例如：北京",
			},
			"unit": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"celsius", "fahrenheit"},
				"description": "温度单位，celsius 表示摄氏度，fahrenheit 表示华氏度",
			},
		},
		"required": []string{"location"},
	}

	tools := []models.ChatTool{
		{
			Type: "function",
			Function: &models.ChatToolFunction{
				Name:        "get_weather",
				Description: "获取指定城市的天气信息",
				Parameters:  params,
			},
		},
	}

	// 创建聊天请求
	return client.CreateChatCompletion(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.DeepSeek,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "北京今天的天气怎么样？",
			},
		},
		Model: consts.DeepSeekChat,
		Tools: tools,
		ToolChoice: &models.ChatToolChoice{
			ToolChoiceType: models.ChatToolChoiceTypeAuto,
		},
		Temperature:         models.Float32(0.7),
		MaxCompletionTokens: models.Int(1024),
	}, httpclient.WithTimeout(time.Minute*2))
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponseStream, err error) {
	return client.CreateChatCompletionStream(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.DeepSeek,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好",
				MultimodalContent: []models.ChatUserMsgPart{
					{
						Type: models.ChatUserMsgPartTypeImageURL,
						ImageURL: &models.ChatUserMsgImageURL{
							URL:    "https://www.gstatic.com/webp/gallery/1.webp",
							Detail: models.ChatUserMsgImageURLDetailHigh,
						},
					},
					{
						Type: models.ChatUserMsgPartTypeText,
						Text: "这些是什么?",
					},
				}, // 不会被序列化，也不会复制到Content字段，因为deepseek不支持多模态
			},
		},
		Model:               consts.DeepSeekReasoner,
		MaxCompletionTokens: models.Int(4096),
		Stream:              models.Bool(true),
		StreamOptions: &models.ChatStreamOptions{
			IncludeUsage: models.Bool(true),
		},
	}, httpclient.WithTimeout(time.Minute*5), httpclient.WithStreamReturnIntervalTimeout(time.Second*5))
}

// 模拟天气查询函数
func getWeather(location string, unit string) string {
	// 实际项目中这里会调用天气API
	return fmt.Sprintf(`{"weather": "sunny", "temperature": 25, "unit": "%s"}`, unit)
}

// processToolCalls 处理工具调用
func processToolCalls(toolCalls []models.ToolCalls) ([]*models.ToolMessage, error) {
	var responses []*models.ToolMessage

	for _, toolCall := range toolCalls {
		switch toolCall.Function.Name {
		case "get_weather":
			// 解析参数
			var args struct {
				Location string `json:"location"`
				Unit     string `json:"unit"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				return nil, fmt.Errorf("failed to parse weather args: %v", err)
			}

			// 获取天气
			weather := getWeather(args.Location, args.Unit)

			// 构建响应
			response := &models.ToolMessage{
				Role:       "tool",
				Content:    weather,
				ToolCallID: toolCall.ID,
			}
			responses = append(responses, response)

		default:
			return nil, fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
		}
	}

	return responses, nil
}

func main() {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		log.Printf("Failed to create temporary test directory: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v", err)
		return
	}

	configData := `{
		"providers": {
			"deepseek": {
				"base_url": "https://api.deepseek.com",
				"api_keys": [%v]
			}
		}
	}`

	configData = fmt.Sprintf(configData, getApiKeys("DEEPSEEK_API_KEYS"))
	log.Printf("configData: %s", configData)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		log.Printf("Failed to create empty config file: %v", err)
		return
	}

	client, err := aisdk.NewSDKClient(configPath, aisdk.WithDefaultMiddlewares())
	if err != nil {
		log.Printf("NewSDKClient() error = %v", err)
		return
	}
	defer func() {
		metrics := client.GetMetrics()
		log.Printf("metrics = %s\n", httpclient.MustString(metrics))
	}()

	ctx := context.Background()
	// 列出模型
	response1, err := listModels(ctx, client)
	isError(err)
	if err != nil {
		log.Printf("listModels error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	log.Printf("listModels response = %s, request_id = %s", httpclient.MustString(response1), response1.RequestID())

	// 创建聊天（启用工具调用）
	response2, err := createChatCompletion(ctx, client)
	isError(err)
	if err != nil {
		log.Printf("createChatCompletion error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}

	// 处理工具调用
	if len(response2.Choices) > 0 && len(response2.Choices[0].Message.ToolCalls) > 0 {
		log.Printf("Tool calls detected: %s", httpclient.MustString(response2.Choices[0].Message.ToolCalls))

		// 处理工具调用
		toolResponses, err := processToolCalls(response2.Choices[0].Message.ToolCalls)
		if err != nil {
			log.Printf("Failed to process tool calls: %v", err)
			return
		}

		// 将工具调用结果发送回模型
		followUpResponse, err := client.CreateChatCompletion(ctx, models.ChatRequest{
			UserInfo: models.UserInfo{User: "123456"},
			Provider: consts.DeepSeek,
			Model:    consts.DeepSeekChat,
			Messages: []models.ChatMessage{
				&models.UserMessage{Content: "北京今天的天气怎么样？"},
				&models.AssistantMessage{
					Role:      "assistant",
					Content:   "",
					ToolCalls: response2.Choices[0].Message.ToolCalls,
				},
				toolResponses[0],
			},
			ToolChoice: &models.ChatToolChoice{
				ToolChoiceType: models.ChatToolChoiceTypeNone, // 不需要再调用工具
			},
		}, httpclient.WithTimeout(time.Minute*2))

		if err != nil {
			log.Printf("Follow-up completion error: %v", err)
			return
		}

		log.Printf("Final response: %s", httpclient.MustString(followUpResponse.Choices[0].Message.Content))
	} else {
		log.Printf("Response: %s", httpclient.MustString(response2.Choices[0].Message.Content))
	}

	// // 创建流式聊天
	// response3, err := createChatCompletionStream(ctx, client)
	// isError(err)
	// if err != nil {
	// 	log.Printf("createChatCompletionStream error = %v, request_id = %s", err, errors.RequestID(err))
	// 	return
	// }
	// // 读取流式聊天
	// log.Printf("createChatCompletionStream request_id = %s", response3.RequestID())
	// if err = response3.ForEach(func(item models.ChatBaseResponse, isFinished bool) (err error) {
	// 	if isFinished {
	// 		return nil
	// 	}
	// 	log.Printf("createChatCompletionStream item = %s", httpclient.MustString(item))
	// 	if item.Usage != nil && item.StreamStats != nil {
	// 		log.Printf("createChatCompletionStream usage = %s", httpclient.MustString(item.Usage))
	// 		log.Printf("createChatCompletionStream stream_stats = %s", httpclient.MustString(item.StreamStats))
	// 	}
	// 	return nil
	// }); err != nil {
	// 	isError(err)
	// 	log.Printf("createChatCompletionStream item error = %v", err)
	// 	return
	// }
}
