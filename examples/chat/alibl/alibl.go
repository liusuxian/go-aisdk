/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-25 13:01:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-11 18:54:28
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
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
	list := strings.Split(os.Getenv(envKey), ",")
	for i, v := range list {
		if i == 0 {
			apiKeys = fmt.Sprintf(`"%s"`, v)
		} else {
			apiKeys = fmt.Sprintf(`%s,"%s"`, apiKeys, v)
		}
	}
	return
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

func createChatCompletion(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponse, err error) {
	return client.CreateChatCompletion(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.AliBL,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
				MultimodalContent: []models.ChatUserMsgPart{
					{
						Type: models.ChatUserMsgPartTypeImageURL, // 不会被序列化
						ImageURL: &models.ChatUserMsgImageURL{
							URL:    "https://www.gstatic.com/webp/gallery/1.webp",
							Detail: models.ChatUserMsgImageURLDetailHigh, // 不会被序列化
						},
					},
					{
						Type: models.ChatUserMsgPartTypeText, // 不会被序列化
						Text: "这些是什么?",
					},
					{
						InputVideo: &models.ChatUserMsgInputVideo{
							VideoImgList: []string{
								"https://www.gstatic.com/webp/gallery/1.webp",
								"https://www.gstatic.com/webp/gallery/2.webp",
								"https://www.gstatic.com/webp/gallery/3.webp",
								"https://www.gstatic.com/webp/gallery/4.webp",
								"https://www.gstatic.com/webp/gallery/5.webp",
							},
						},
					},
				},
			},
		},
		Model:               consts.AliBLQwenLong,
		FrequencyPenalty:    models.Float32(1.0),
		MaxCompletionTokens: models.Int(4096),
		// Metadata:            map[string]string{"X-DashScope-DataInspection": "{\"input\": \"cip\", \"output\": \"cip\"}"},
		WebSearchOptions: &models.ChatWebSearchOptions{
			EnableSource:   models.Bool(true),
			EnableCitation: models.Bool(true),
			CitationFormat: models.ChatCitationFormatRefNumber,
			ForcedSearch:   models.Bool(true),
			SearchStrategy: models.ChatSearchStrategyPro,
		},
		EnableThinking: models.Bool(true),
	}, httpclient.WithTimeout(time.Minute*5))
}

func createChatCompletionStream(ctx context.Context, client *aisdk.SDKClient) (response models.ChatResponseStream, err error) {
	return client.CreateChatCompletionStream(ctx, models.ChatRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider: consts.AliBL,
		Messages: []models.ChatMessage{
			&models.UserMessage{
				Content: "你好，我是小明，请帮我写一个关于人工智能的论文",
				MultimodalContent: []models.ChatUserMsgPart{
					{
						Type: models.ChatUserMsgPartTypeImageURL, // 不会被序列化
						ImageURL: &models.ChatUserMsgImageURL{
							URL:    "https://www.gstatic.com/webp/gallery/1.webp",
							Detail: models.ChatUserMsgImageURLDetailHigh, // 不会被序列化
						},
					},
					{
						Type: models.ChatUserMsgPartTypeText, // 不会被序列化
						Text: "这些是什么?",
					},
				},
			},
		},
		Model:               consts.AliBLQwenLong,
		FrequencyPenalty:    models.Float32(1.0),
		MaxCompletionTokens: models.Int(4096),
		// Metadata:            map[string]string{"X-DashScope-DataInspection": "{\"input\": \"cip\", \"output\": \"cip\"}"},
		WebSearchOptions: &models.ChatWebSearchOptions{
			EnableSource:   models.Bool(true),
			EnableCitation: models.Bool(true),
			CitationFormat: models.ChatCitationFormatRefNumber,
			ForcedSearch:   models.Bool(true),
			SearchStrategy: models.ChatSearchStrategyPro,
		},
		EnableThinking: models.Bool(true),
		Stream:         models.Bool(true), // 不会被序列化，会放到请求头中
		StreamOptions: &models.ChatStreamOptions{
			IncludeUsage: models.Bool(true),
		}, // 不会被序列化
	}, httpclient.WithTimeout(time.Minute*5), httpclient.WithStreamReturnIntervalTimeout(time.Second*5))
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
    "alibl": {
			"base_url": "https://dashscope.aliyuncs.com/api/v1",
      "api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("ALIBL_API_KEYS"))
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
	// 创建聊天
	response1, err := createChatCompletion(ctx, client)
	isError(err)
	if err != nil {
		log.Printf("createChatCompletion error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	log.Printf("createChatCompletion response = %s, request_id = %s", httpclient.MustString(response1), response1.RequestID())
	// 创建流式聊天
	response2, err := createChatCompletionStream(ctx, client)
	isError(err)
	if err != nil {
		log.Printf("createChatCompletionStream error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	// 读取流式聊天
	log.Printf("createChatCompletionStream request_id = %s", response2.RequestID())
	if err = response2.ForEach(func(item models.ChatBaseResponse, isFinished bool) (err error) {
		if isFinished {
			return nil
		}
		log.Printf("createChatCompletionStream item = %s", httpclient.MustString(item))
		if item.Usage != nil && item.StreamStats != nil {
			log.Printf("createChatCompletionStream usage = %s", httpclient.MustString(item.Usage))
			log.Printf("createChatCompletionStream stream_stats = %s", httpclient.MustString(item.StreamStats))
		}
		return nil
	}); err != nil {
		isError(err)
		log.Printf("createChatCompletionStream item error = %v", err)
		return
	}
}
