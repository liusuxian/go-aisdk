/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-19 17:11:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-09 14:49:55
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package models

import (
	"io"

	"github.com/Mrzhouyl/go-aisdk/consts"
	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/internal/utils"
)

// ImageBackground 背景透明度类型
type ImageBackground string

const (
	// 透明
	//
	// 提供商支持: OpenAI
	ImageBackgroundTransparent ImageBackground = "transparent"
	// 不透明
	//
	// 提供商支持: OpenAI
	ImageBackgroundOpaque ImageBackground = "opaque"
	// 自动
	//
	// 提供商支持: OpenAI
	ImageBackgroundAuto ImageBackground = "auto"
)

// ImageModeration 内容审核级别
type ImageModeration string

const (
	// 低
	//
	// 提供商支持: OpenAI
	ImageModerationLow ImageModeration = "low"
	// 自动
	//
	// 提供商支持: OpenAI
	ImageModerationAuto ImageModeration = "auto"
)

// ImageOutputFormat 返回图像格式
type ImageOutputFormat string

const (
	// PNG
	//
	// 提供商支持: OpenAI
	ImageOutputFormatPNG ImageOutputFormat = "png"
	// JPEG
	//
	// 提供商支持: OpenAI
	ImageOutputFormatJPEG ImageOutputFormat = "jpeg"
	// WebP
	//
	// 提供商支持: OpenAI
	ImageOutputFormatWebP ImageOutputFormat = "webp"
)

// ImageQuality 图像质量
type ImageQuality string

const (
	// 标准
	//
	// 提供商支持: OpenAI
	ImageQualityStandard ImageQuality = "standard"
	// 高清
	//
	// 提供商支持: OpenAI
	ImageQualityHD ImageQuality = "hd"
	// 高
	//
	// 提供商支持: OpenAI
	ImageQualityHigh ImageQuality = "high"
	// 中
	//
	// 提供商支持: OpenAI
	ImageQualityMedium ImageQuality = "medium"
	// 低
	//
	// 提供商支持: OpenAI
	ImageQualityLow ImageQuality = "low"
	// 自动
	//
	// 提供商支持: OpenAI
	ImageQualityAuto ImageQuality = "auto"
)

// ImageResponseFormat 响应格式
type ImageResponseFormat string

const (
	// URL
	//
	// 提供商支持: OpenAI
	ImageResponseFormatURL ImageResponseFormat = "url"
	// Base64
	//
	// 提供商支持: OpenAI
	ImageResponseFormatB64JSON ImageResponseFormat = "b64_json"
)

// ImageSize 图像尺寸
type ImageSize string

const (
	// 自动
	//
	// 提供商支持: OpenAI
	ImageSizeAuto ImageSize = "auto"
	// 256x256
	//
	// 提供商支持: OpenAI
	ImageSize256x256 ImageSize = "256x256"
	// 512x512
	//
	// 提供商支持: OpenAI
	ImageSize512x512 ImageSize = "512x512"
	// 1024x1024
	//
	// 提供商支持: OpenAI
	ImageSize1024x1024 ImageSize = "1024x1024"
	// 1792x1024
	//
	// 提供商支持: OpenAI
	ImageSize1792x1024 ImageSize = "1792x1024"
	// 1024x1792
	//
	// 提供商支持: OpenAI
	ImageSize1024x1792 ImageSize = "1024x1792"
	// 1536x1024
	//
	// 提供商支持: OpenAI
	ImageSize1536x1024 ImageSize = "1536x1024"
	// 1024x1536
	//
	// 提供商支持: OpenAI
	ImageSize1024x1536 ImageSize = "1024x1536"
)

// ImageStyle 图像风格
type ImageStyle string

const (
	// 超现实
	//
	// 提供商支持: OpenAI
	ImageStyleVivid ImageStyle = "vivid"
	// 自然
	//
	// 提供商支持: OpenAI
	ImageStyleNatural ImageStyle = "natural"
)

// ImageRequest 创建图像请求
type ImageRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 提示词
	//
	// 提供商支持: OpenAI
	Prompt string `json:"prompt,omitempty" providers:"openai"`
	// 设置生成图像的背景透明度
	//
	// 提供商支持: OpenAI
	Background ImageBackground `json:"background,omitempty" providers:"openai"`
	// 模型名称
	//
	// 提供商支持: OpenAI
	Model string `json:"model,omitempty" providers:"openai"`
	// 内容审核级别
	//
	// 提供商支持: OpenAI
	Moderation ImageModeration `json:"moderation,omitempty" providers:"openai"`
	// 生成的图像数量
	//
	// 提供商支持: OpenAI
	N int `json:"n,omitempty" providers:"openai"`
	// 图像压缩级别(0-100%)
	//
	// 提供商支持: OpenAI
	OutputCompression int `json:"output_compression,omitempty" providers:"openai"`
	// 返回图像格式
	//
	// 提供商支持: OpenAI
	OutputFormat ImageOutputFormat `json:"output_format,omitempty" providers:"openai"`
	// 图像质量
	//
	// 提供商支持: OpenAI
	Quality ImageQuality `json:"quality,omitempty" providers:"openai"`
	// 响应格式
	//
	// 提供商支持: OpenAI
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty" providers:"openai"`
	// 图像尺寸
	//
	// 提供商支持: OpenAI
	Size ImageSize `json:"size,omitempty" providers:"openai"`
	// 图像风格
	//
	// 提供商支持: OpenAI
	Style ImageStyle `json:"style,omitempty" providers:"openai"`
}

// MarshalJSON 序列化JSON
func (r ImageRequest) MarshalJSON() (b []byte, err error) {
	provider := r.Provider.String()
	// 序列化JSON
	r.Provider = ""
	return utils.NewSerializer(provider).Serialize(r)
}

// ImageEditRequest 编辑图像请求
//
//	提供商支持: OpenAI
type ImageEditRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 要编辑的图像源数组
	//
	// 提供商支持: OpenAI
	Image []io.Reader `json:"image,omitempty" providers:"openai"`
	// 提示词
	//
	// 提供商支持: OpenAI
	Prompt string `json:"prompt,omitempty" providers:"openai"`
	// 设置生成图像的背景透明度
	//
	// 提供商支持: OpenAI
	Background ImageBackground `json:"background,omitempty" providers:"openai"`
	// mask图像源，其中完全透明的区域指示应该编辑的位置
	//
	// 提供商支持: OpenAI
	Mask io.Reader `json:"mask,omitempty" providers:"openai"`
	// 模型名称
	//
	// 提供商支持: OpenAI
	Model string `json:"model,omitempty" providers:"openai"`
	// 生成的图像数量。必须在1到10之间
	//
	// 提供商支持: OpenAI
	N int `json:"n,omitempty" providers:"openai"`
	// 图像压缩级别(0-100%)
	//
	// 提供商支持: OpenAI
	OutputCompression int `json:"output_compression,omitempty" providers:"openai"`
	// 返回图像格式
	//
	// 提供商支持: OpenAI
	OutputFormat ImageOutputFormat `json:"output_format,omitempty" providers:"openai"`
	// 图像质量
	//
	// 提供商支持: OpenAI
	Quality ImageQuality `json:"quality,omitempty" providers:"openai"`
	// 响应格式
	//
	// 提供商支持: OpenAI
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty" providers:"openai"`
	// 图像尺寸
	//
	// 提供商支持: OpenAI
	Size ImageSize `json:"size,omitempty" providers:"openai"`
}

// MarshalJSON 序列化JSON
func (r ImageEditRequest) MarshalJSON() (b []byte, err error) {
	provider := r.Provider.String()
	// 序列化JSON
	r.Provider = ""
	return utils.NewSerializer(provider).Serialize(r)
}

// ImageVariationRequest 变换图像请求
//
//	提供商支持: OpenAI
type ImageVariationRequest struct {
	UserInfo
	Provider consts.Provider `json:"provider,omitempty"` // 提供商
	// 用作变换的基础图像。必须是有效的PNG文件，小于4MB，且为正方形
	//
	// 提供商支持: OpenAI
	Image io.Reader `json:"image,omitempty" providers:"openai"`
	// 模型名称
	//
	// 提供商支持: OpenAI
	Model string `json:"model,omitempty" providers:"openai"`
	// 生成的图像数量
	//
	// 提供商支持: OpenAI
	N int `json:"n,omitempty" providers:"openai"`
	// 响应格式
	//
	// 提供商支持: OpenAI
	ResponseFormat ImageResponseFormat `json:"response_format,omitempty" providers:"openai"`
	// 图像尺寸
	//
	// 提供商支持: OpenAI
	Size ImageSize `json:"size,omitempty" providers:"openai"`
}

// MarshalJSON 序列化JSON
func (r ImageVariationRequest) MarshalJSON() (b []byte, err error) {
	provider := r.Provider.String()
	// 序列化JSON
	r.Provider = ""
	return utils.NewSerializer(provider).Serialize(r)
}

// ImageResponseData 生成图像的列表
type ImageResponseData struct {
	URL           string `json:"url,omitempty"`            // 生成图像的URL地址
	B64JSON       string `json:"b64_json,omitempty"`       // 生成图像的base64编码JSON
	RevisedPrompt string `json:"revised_prompt,omitempty"` // 用于生成图像的修订提示词
}

// ImageUsageInputTokensDetails 图像生成输入token的详细信息
type ImageUsageInputTokensDetails struct {
	ImageTokens int `json:"image_tokens,omitempty"` // 输入提示中的图像token数量
	TextTokens  int `json:"text_tokens,omitempty"`  // 输入提示中的文本token数量
}

// ImageUsage 图像生成的token使用信息
type ImageUsage struct {
	InputTokens        int                           `json:"input_tokens,omitempty"`         // 输入提示中的token总数（包括图像和文本）
	InputTokensDetails *ImageUsageInputTokensDetails `json:"input_tokens_details,omitempty"` // 图像生成输入token的详细信息
	OutputTokens       int                           `json:"output_tokens,omitempty"`        // 输出图像中的图像token数量
	TotalTokens        int                           `json:"total_tokens,omitempty"`         // 图像生成使用的token总数（图像和文本）
}

// ImageResponse 创建图像响应
type ImageResponse struct {
	Created int64               `json:"created,omitempty"` // 创建图像完成时的 Unix 时间戳（以秒为单位）
	Data    []ImageResponseData `json:"data,omitempty"`    // 生成图像的列表
	Usage   *ImageUsage         `json:"usage,omitempty"`   // 图像生成的token使用信息
	httpclient.HttpHeader
}
