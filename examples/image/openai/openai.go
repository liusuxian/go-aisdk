/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-11 14:53:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-07 23:48:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"io"
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
	"github.com/Mrzhouyl/go-aisdk/utils"
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

func createImage(ctx context.Context, client *aisdk.SDKClient) (response models.ImageResponse, err error) {
	return client.CreateImage(ctx, models.ImageRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider:     consts.OpenAI,
		Prompt:       "一间有着精致雕花窗户的花店，漂亮的深色木质门上挂着铜制把手。店内摆放着各式各样的鲜花，包括玫瑰、百合和向日葵，色彩鲜艳，生机勃勃。背景是温馨的室内场景，透过窗户可以看到街道。高清写实摄影，中景构图。",
		Model:        consts.OpenAIGPTImage1,
		N:            2,
		OutputFormat: models.ImageOutputFormatJPEG,
		Quality:      models.ImageQualityHigh,
		Size:         models.ImageSize1024x1024,
	}, httpclient.WithTimeout(time.Minute*5))
}

func createImageEdit(ctx context.Context, client *aisdk.SDKClient, filenames []string) (response models.ImageResponse, err error) {
	imageReaders := make([]io.Reader, 0, len(filenames)*2)
	for i, filename := range filenames {
		var reader *utils.ImageReader
		if reader, err = utils.FileToReader(filename); err != nil {
			return
		}
		imageReaders = append(imageReaders, reader)
		if reader, err = utils.URLToReader(fmt.Sprintf("https://www.gstatic.com/webp/gallery/%d.webp", i+1), fmt.Sprintf("%d", i+1), time.Second*10); err != nil {
			return
		}
		imageReaders = append(imageReaders, reader)
	}
	return client.CreateImageEdit(ctx, models.ImageEditRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider:     consts.OpenAI,
		Image:        imageReaders,
		Prompt:       "合并这些图片，并添加文字：这是一张美丽的图片。",
		Model:        consts.OpenAIGPTImage1,
		N:            2,
		OutputFormat: models.ImageOutputFormatJPEG,
		Quality:      models.ImageQualityHigh,
		Size:         models.ImageSize1024x1024,
	}, httpclient.WithTimeout(time.Minute*5))
}

func createImageVariation(ctx context.Context, client *aisdk.SDKClient, filename string) (response models.ImageResponse, err error) {
	var imageReader *utils.ImageReader
	if imageReader, err = utils.FileToReader(filename); err != nil {
		return
	}
	return client.CreateImageVariation(ctx, models.ImageVariationRequest{
		UserInfo: models.UserInfo{
			User: "123456",
		},
		Provider:       consts.OpenAI,
		Image:          imageReader,
		Model:          consts.OpenAIDallE2,
		N:              2,
		ResponseFormat: models.ImageResponseFormatURL,
		Size:           models.ImageSize1024x1024,
	}, httpclient.WithTimeout(time.Minute*5))
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
    "openai": {
      "base_url": "https://niubi.zeabur.app/v1",
			"api_keys": [%v]
    }
  }
}`
	configData = fmt.Sprintf(configData, getApiKeys("OPENAI_API_KEYS"))
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
		log.Printf("metrics = %+v\n", metrics)
	}()

	ctx := context.Background()
	// 创建图像
	response1, err := createImage(ctx, client)
	isError(err)
	if err != nil {
		log.Printf("createImage error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	// 保存每张生成的图片
	filenames := make([]string, 0, len(response1.Data))
	for i, v := range response1.Data {
		if v.B64JSON != "" {
			if filename, err := utils.SaveBase64Image(v.B64JSON, "generated_images", fmt.Sprintf("image_%d", i+1)); err != nil {
				log.Printf("save image %d error: %v", i+1, err)
			} else {
				filenames = append(filenames, filename)
			}
		} else {
			log.Printf("image %d base64 data is empty", i+1)
		}
	}
	// 编辑图像
	response2, err := createImageEdit(ctx, client, filenames)
	isError(err)
	if err != nil {
		log.Printf("createImageEdit error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	// 保存每张编辑的图片
	for i, v := range response2.Data {
		if v.B64JSON != "" {
			if _, err := utils.SaveBase64Image(v.B64JSON, "generated_images", fmt.Sprintf("image_edit_%d", i+1)); err != nil {
				log.Printf("save image edit %d error: %v", i+1, err)
			}
		} else {
			log.Printf("image edit %d base64 data is empty", i+1)
		}
	}
	// 变换图像
	response3, err := createImageVariation(ctx, client, filenames[0])
	isError(err)
	if err != nil {
		log.Printf("createImageVariation error = %v, request_id = %s", err, errors.RequestID(err))
		return
	}
	// 保存每张变换的图片
	for i, v := range response3.Data {
		if v.URL != "" {
			if _, err := utils.SaveURLImage(v.URL, "generated_images", fmt.Sprintf("image_variation_%d", i+1), time.Second*10); err != nil {
				log.Printf("save image variation %d error: %v", i+1, err)
			}
		} else {
			log.Printf("image variation %d base64 data is empty", i+1)
		}
	}
}
