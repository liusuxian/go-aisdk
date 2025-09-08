/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-07-03 15:11:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-07-03 15:12:43
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// specialFieldContext 特殊字段上下文
type specialFieldContext struct {
	structValue   reflect.Value
	field         reflect.StructField
	fieldValue    reflect.Value
	jsonFieldName string
	visited       *map[uintptr]bool
	result        *map[string]any
	copyToResult  *map[string]any
	depth         int
}

// specialFieldHandler 特殊字段处理函数
type specialFieldHandler func(ctx *specialFieldContext) (isSkip bool, err error)

// Serializer 序列化器
type Serializer struct {
	provider string // 提供商
	maxDepth int    // 最大深度
}

// NewSerializer 创建序列化器
func NewSerializer(provider string, maxDepth ...int) (s *Serializer) {
	var depth int
	if len(maxDepth) > 0 && maxDepth[0] > 0 {
		depth = maxDepth[0]
	} else {
		depth = 100
	}
	return &Serializer{
		provider: strings.ToLower(provider),
		maxDepth: depth,
	}
}

// Serialize 序列化对象为JSON
func (s *Serializer) Serialize(obj any) (b []byte, err error) {
	// 创建局部的循环引用检测map
	visited := make(map[uintptr]bool)

	v := reflect.ValueOf(obj)
	// 处理指针
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("expected struct, got %s", v.Kind())
		return
	}

	var (
		result  map[string]any
		isValid bool
	)
	if result, isValid, err = s.serializeStruct(v, &visited, 0); err != nil {
		return
	}
	if !isValid {
		return
	}
	return json.Marshal(result)
}

// serializeValue 递归序列化值
func (s *Serializer) serializeValue(v reflect.Value, visited *map[uintptr]bool, depth int) (result any, isValid bool, err error) {
	// 检查递归深度
	if depth > s.maxDepth {
		err = fmt.Errorf("serialization depth exceeded maximum limit of %d at depth %d", s.maxDepth, depth)
		return
	}

	if !v.IsValid() {
		return
	}
	// 处理接口类型
	if v.Kind() == reflect.Interface && !v.IsNil() {
		return s.serializeValue(v.Elem(), visited, depth)
	}
	// 处理指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		// 循环引用检测
		ptr := v.Pointer()
		if (*visited)[ptr] {
			return
		}
		(*visited)[ptr] = true
		defer delete((*visited), ptr) // 序列化完成后清理，允许同一对象在不同路径中出现
		return s.serializeValue(v.Elem(), visited, depth+1)
	}
	// 优先处理实现了 json.Marshaler 的类型（尊重自定义序列化逻辑）
	if v.CanInterface() {
		if marshaler, ok := v.Interface().(json.Marshaler); ok {
			var b []byte
			if b, err = marshaler.MarshalJSON(); err != nil {
				return
			}
			// 将自定义 JSON 反序列化为任意类型以合并到结果中
			var anyVal any
			if len(b) == 0 {
				return
			}
			if err = json.Unmarshal(b, &anyVal); err != nil {
				// 如果不是标准 JSON（例如返回的是裸字符串而不是带引号），尝试按字符串处理
				anyVal = string(b)
			}
			result = anyVal
			isValid = true
			return
		}
	}
	// 处理不同类型
	switch v.Kind() {
	case reflect.Struct:
		return s.serializeStruct(v, visited, depth+1)
	case reflect.Slice, reflect.Array:
		return s.serializeSlice(v, visited, depth+1)
	case reflect.Map:
		return s.serializeMap(v, visited, depth+1)
	case reflect.Chan, reflect.Func:
		// 通道和函数类型无法序列化，返回nil
		return
	case reflect.UnsafePointer, reflect.Uintptr:
		// 不安全指针类型，为安全起见返回nil
		return
	default:
		// 基本类型处理 - 添加安全检查
		if !v.CanInterface() {
			err = fmt.Errorf("cannot serialize unexported field of type %s", v.Type())
			return
		}
		result = v.Interface()
		isValid = true
		return
	}
}

// serializeStruct 序列化结构体
func (s *Serializer) serializeStruct(v reflect.Value, visited *map[uintptr]bool, depth int) (result map[string]any, isValid bool, err error) {
	result = make(map[string]any)
	var (
		copyToResult = make(map[string]any)
		vType        = v.Type()
	)

	for i := 0; i < v.NumField(); i++ {
		var (
			field      = vType.Field(i)
			fieldValue = v.Field(i)
		)
		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}
		// 判断是否应该跳过该字段
		if s.shouldSkipField(field, fieldValue) {
			continue
		}
		// 处理特殊字段，每个特殊字段是互斥的，不能同时存在
		var (
			jsonFieldName = s.getJsonFieldName(field)
			isSkip        bool
		)
		for _, handler := range []specialFieldHandler{
			s.handleAnonymousEmbeddedField, // 处理匿名嵌入字段
			s.handleCopyToField,            // 处理copyto标签字段
			s.handleFlattenField,           // 处理flatten标签字段
			s.handleGroupField,             // 处理group标签字段
			s.handleDefaultField,           // 处理default标签字段
		} {
			if isSkip, err = handler(&specialFieldContext{
				structValue:   v,
				field:         field,
				fieldValue:    fieldValue,
				jsonFieldName: jsonFieldName,
				visited:       visited,
				result:        &result,
				copyToResult:  &copyToResult,
				depth:         depth,
			}); err != nil {
				return
			}
			if isSkip {
				break
			}
		}
		if isSkip {
			continue
		}
		// 检查该字段是否已经被copyto特殊字段处理过了，如果处理过，则不处理该字段
		if _, exists := copyToResult[jsonFieldName]; exists {
			continue
		}
		// 递归序列化字段值
		var (
			serializedValue any
			serializedValid bool
		)
		if serializedValue, serializedValid, err = s.serializeValue(fieldValue, visited, depth); err != nil {
			err = fmt.Errorf("error serializing field %s: %w", field.Name, err)
			return
		}
		if serializedValid {
			result[jsonFieldName] = serializedValue
		}
	}

	if len(result) == 0 {
		return
	}
	isValid = true
	return
}

// serializeSlice 序列化切片和数组
func (s *Serializer) serializeSlice(v reflect.Value, visited *map[uintptr]bool, depth int) (result []any, isValid bool, err error) {
	// 只有切片可能为nil，数组不会
	if v.Kind() == reflect.Slice && v.IsNil() {
		return
	}

	result = make([]any, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		var (
			item      any
			itemValid bool
		)
		if item, itemValid, err = s.serializeValue(v.Index(i), visited, depth); err != nil {
			err = fmt.Errorf("error serializing %s item %d: %w", v.Kind(), i, err)
			return
		}
		if itemValid {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return
	}
	isValid = true
	return
}

// serializeMap 序列化映射
func (s *Serializer) serializeMap(v reflect.Value, visited *map[uintptr]bool, depth int) (result map[string]any, isValid bool, err error) {
	if v.IsNil() {
		return
	}

	result = make(map[string]any)
	for _, key := range v.MapKeys() {
		var (
			keyStr     = fmt.Sprintf("%v", key.Interface())
			value      any
			valueValid bool
		)
		if value, valueValid, err = s.serializeValue(v.MapIndex(key), visited, depth); err != nil {
			err = fmt.Errorf("error serializing map value for key %s: %w", keyStr, err)
			return
		}
		if valueValid {
			result[keyStr] = value
		}
	}

	if len(result) == 0 {
		return
	}
	isValid = true
	return
}

// 处理匿名嵌入字段
func (s *Serializer) handleAnonymousEmbeddedField(ctx *specialFieldContext) (isSkip bool, err error) {
	if ctx.field.Anonymous {
		// 展开结构体字段
		var (
			flattenedFields map[string]any
			shouldFlatten   bool
		)
		if flattenedFields, shouldFlatten, err = s.flattenField(ctx.fieldValue, ctx.visited, ctx.depth); err != nil {
			return
		}
		if shouldFlatten {
			// 检查字段冲突，避免覆盖已存在的普通字段
			for key, value := range flattenedFields {
				if _, exists := (*ctx.result)[key]; !exists {
					(*ctx.result)[key] = value
				}
			}
		}
		isSkip = true
		return
	}
	return
}

// 处理copyto标签字段
func (s *Serializer) handleCopyToField(ctx *specialFieldContext) (isSkip bool, err error) {
	if copyToTag := ctx.field.Tag.Get("copyto"); copyToTag != "" {
		// 获取copyto标签字段对应的JSON字段名
		var targetJsonFieldName string
		if targetJsonFieldName, err = s.getCopyToTargetJsonFieldName(copyToTag, ctx.structValue); err != nil {
			return
		}
		if targetJsonFieldName == "" {
			// copyto tag 没有当前提供商的配置，之后会当成普通字段进行处理
			return
		}
		// 递归序列化值
		var (
			sourceValue any
			sourceValid bool
		)
		if sourceValue, sourceValid, err = s.serializeValue(ctx.fieldValue, ctx.visited, ctx.depth); err != nil {
			err = fmt.Errorf("error serializing copyto field %s: %w", ctx.field.Name, err)
			return
		}
		if sourceValid {
			(*ctx.result)[targetJsonFieldName] = sourceValue
			(*ctx.copyToResult)[targetJsonFieldName] = true // 记录目标字段被处理过
		}
		isSkip = true
		return
	}
	return
}

// 处理flatten标签字段
func (s *Serializer) handleFlattenField(ctx *specialFieldContext) (isSkip bool, err error) {
	if flattenTag := ctx.field.Tag.Get("flatten"); flattenTag != "" && s.isSupported(flattenTag) {
		// 展开结构体字段
		var (
			flattenedFields map[string]any
			shouldFlatten   bool
		)
		if flattenedFields, shouldFlatten, err = s.flattenField(ctx.fieldValue, ctx.visited, ctx.depth); err != nil {
			return
		}
		if shouldFlatten {
			// 检查字段冲突，避免覆盖已存在的普通字段
			for key, value := range flattenedFields {
				if _, exists := (*ctx.result)[key]; !exists {
					(*ctx.result)[key] = value
				}
			}
		}
		isSkip = true
		return
	}
	return
}

// 处理group标签字段
func (s *Serializer) handleGroupField(ctx *specialFieldContext) (isSkip bool, err error) {
	if groupTag := ctx.field.Tag.Get("group"); groupTag != "" {
		// 获取当前提供商的分组名
		var groupName string
		if groupName, err = s.getTagValue(groupTag); err != nil {
			return
		}
		if groupName == "" {
			// group tag 没有当前提供商的配置，之后会当成普通字段进行处理
			return
		}
		// 递归序列化字段值
		var (
			serializedValue any
			serializedValid bool
		)
		if serializedValue, serializedValid, err = s.serializeValue(ctx.fieldValue, ctx.visited, ctx.depth); err != nil {
			err = fmt.Errorf("error serializing group field %s: %w", ctx.field.Name, err)
			return
		}
		if serializedValid {
			// 直接在 result 中处理分组，检查是否已存在
			if existingGroup, exists := (*ctx.result)[groupName]; exists {
				// 如果分组已存在，合并字段
				if existingGroupMap, ok := existingGroup.(map[string]any); ok {
					// 只添加不存在的字段，避免覆盖
					if _, fieldExists := existingGroupMap[ctx.jsonFieldName]; !fieldExists {
						existingGroupMap[ctx.jsonFieldName] = serializedValue
					}
				} else {
					// 如果已存在的不是 map 类型，创建新的分组并添加字段
					(*ctx.result)[groupName] = map[string]any{
						ctx.jsonFieldName: serializedValue,
					}
				}
			} else {
				// 分组不存在，创建新的分组
				(*ctx.result)[groupName] = map[string]any{
					ctx.jsonFieldName: serializedValue,
				}
			}
		}
		isSkip = true
		return
	}
	return
}

// 处理default标签字段
func (s *Serializer) handleDefaultField(ctx *specialFieldContext) (isSkip bool, err error) {
	if defaultTag := ctx.field.Tag.Get("default"); defaultTag != "" {
		// 获取当前提供商的默认值
		var defaultValue any
		if defaultValue, err = s.getDefaultValue(defaultTag); err != nil {
			return
		}
		if defaultValue == nil {
			// default tag 没有当前提供商的配置，之后会当成普通字段进行处理
			return
		}
		// 如果字段值为零值，则设置为默认值
		if ctx.fieldValue.IsZero() {
			(*ctx.result)[ctx.jsonFieldName] = defaultValue
			isSkip = true
			return
		}
	}
	return
}

// flattenField 展开结构体字段
//
//	格式: `flatten:"provider1,provider2,provider3"`
//	如果是结构体类型的匿名字段，则默认展开
//
// 特别说明：
//
//	情况1: 要展开的结构体在前，普通字段在后 -> 普通字段可以覆盖要展开的结构体中的同名字段
//
//	struct定义:
//	  type BaseInfo struct {
//	      Name string `json:"name"`
//	      Age  int    `json:"age"`
//	  }
//	  type Person struct {
//	      BaseInfo                        // 展开字段在前
//	      Name     string `json:"name"`   // 普通字段在后，会覆盖BaseInfo中的Name
//	      Email    string `json:"email"`
//	  }
//	序列化结果: {"name": "Person.Name的值", "age": "BaseInfo.Age的值", "email": "Person.Email的值"}
//
//	情况2: 普通字段在前，要展开的结构体在后 -> 要展开的结构体中的同名字段无法覆盖普通字段
//
//	struct定义:
//	  type BaseInfo struct {
//	      Name string `json:"name"`
//	      Age  int    `json:"age"`
//	  }
//	  type Person struct {
//	      Name     string `json:"name"`   // 普通字段在前
//	      Email    string `json:"email"`
//	      BaseInfo                        // 展开字段在后，其中的Name无法覆盖前面的Name
//	  }
//	序列化结果: {"name": "Person.Name的值", "email": "Person.Email的值", "age": "BaseInfo.Age的值"}
func (s *Serializer) flattenField(fieldValue reflect.Value, visited *map[uintptr]bool, depth int) (flattenedFields map[string]any, shouldFlatten bool, err error) {
	// 处理指针
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return
		}
		fieldValue = fieldValue.Elem()
	}
	// 确保是结构体
	if fieldValue.Kind() != reflect.Struct {
		return
	}
	// 序列化结构体
	return s.serializeStruct(fieldValue, visited, depth)
}

// getCopyToTargetJsonFieldName 获取copyto标签字段对应的JSON字段名
//
//	格式1: `copyto:"targetField"` -> 将字段值复制到指定字段对应的JSON字段名
//	格式2: `copyto:"provider1:targetField1,provider2:targetField2,..."` -> 将字段值复制到提供商指定字段对应的JSON字段名
//	格式3: `copyto:"provider1|provider2|provider3|...:targetField1,provider4|provider5|...:targetField2"` -> 将字段值复制到提供商指定字段对应的JSON字段名
func (s *Serializer) getCopyToTargetJsonFieldName(copyToTag string, structValue reflect.Value) (targetJsonFieldName string, err error) {
	// 获取当前提供商指定tag的值
	var tagValue string
	if tagValue, err = s.getTagValue(copyToTag); err != nil {
		return
	}
	if tagValue == "" {
		// copyToTag 没有当前提供商的配置，返回空
		return
	}
	// 获取目标字段
	targetField, found := structValue.Type().FieldByName(tagValue)
	if !found {
		err = fmt.Errorf("target field %s not found", tagValue)
		return
	}
	// 获取目标字段的JSON字段名
	targetJsonFieldName = s.getJsonFieldName(targetField)
	return
}

// getDefaultValue 获取当前提供商的默认值
//
//	格式1: `default:"provider1:defaultValue1,provider2:defaultValue2,..."` -> 使用指定提供商的默认值
//	格式2: `default:"provider1|provider2|provider3|...:defaultValue1,provider4|provider5|...:defaultValue2"` -> 使用指定提供商的默认值
//	格式3: `default:"defaultValue"` -> 使用通用默认值
func (s *Serializer) getDefaultValue(defaultTag string) (result any, err error) {
	// 获取当前提供商指定tag的值
	var tagValue string
	if tagValue, err = s.getTagValue(defaultTag); err != nil {
		return
	}
	if tagValue == "" {
		// defaultTag 没有当前提供商的配置，返回空
		return
	}
	// 尝试解析JSON格式的默认值，如果失败则当作字符串
	if e := json.Unmarshal([]byte(tagValue), &result); e == nil {
		return
	}
	result = tagValue
	return
}

// shouldSkipField 判断是否应该跳过该字段
func (s *Serializer) shouldSkipField(field reflect.StructField, value reflect.Value) (shouldSkip bool) {
	// json:"-" 表示跳过此字段
	if field.Tag.Get("json") == "-" {
		return true
	}
	// 检查提供商是否被支持
	if !s.isSupported(field.Tag.Get("providers")) {
		return true
	}
	// 如果字段有default标签，则不会因为omitempty跳过（default优先级更高）
	if field.Tag.Get("default") != "" {
		return false
	}
	// 检查 JSON 标签的 omitempty
	return strings.Contains(field.Tag.Get("json"), "omitempty") && value.IsZero()
}

// isSupported 检查提供商是否被支持
//
//	格式: `tag:"provider1,provider2,provider3"`
//	如果tag为空，则默认支持所有提供商
func (s *Serializer) isSupported(tag string) (isSupported bool) {
	if tag == "" {
		return true
	}

	providerList := strings.Split(strings.ToLower(tag), ",")
	return slices.Contains(providerList, s.provider)
}

// getJsonFieldName 获取JSON字段名
//
//	格式1: `json:"jsonFieldName"` -> 使用json字段名
//	格式2: `json:"jsonFieldName" mapping:"provider1:mappedName1,provider2:mappedName2,..."` -> 使用指定提供商的映射名
//	格式3: `json:"jsonFieldName" mapping:"provider1|provider2|provider3|...:mappedName1,provider4|provider5|...:mappedName2"` -> 使用指定提供商的映射名
//	如果tag为空，则使用Go字段名
func (s *Serializer) getJsonFieldName(field reflect.StructField) (jsonFieldName string) {
	jsonTag := field.Tag.Get("json")
	switch jsonTag {
	case "":
		jsonFieldName = field.Name
	default:
		if parts := strings.Split(jsonTag, ","); parts[0] != "" {
			jsonFieldName = parts[0]
		} else {
			jsonFieldName = field.Name
		}
	}
	// 应用字段映射
	if mappingTag := field.Tag.Get("mapping"); mappingTag != "" {
		for mapping := range strings.SplitSeq(mappingTag, ",") {
			if parts := strings.Split(mapping, ":"); len(parts) == 2 {
				var (
					providerList = strings.Split(strings.TrimSpace(parts[0]), "|")
					mappedName   = strings.TrimSpace(parts[1])
				)
				if slices.Contains(providerList, s.provider) && mappedName != "" {
					jsonFieldName = mappedName
					return
				}
			}
		}
	}
	return
}

// getTagValue 获取当前提供商指定tag的值
//
//	格式1: `tag:"provider1:value1,provider2:value2,..."` -> 使用指定提供商的值
//	格式2: `tag:"provider1|provider2|provider3|...:value1,provider4|provider5|...:value2"` -> 使用指定提供商的值
//	格式3: `tag:"value"` -> 使用通用值
func (s *Serializer) getTagValue(tagValue string) (value string, err error) {
	// 检查是否包含提供商特定的值
	if strings.Contains(tagValue, ":") {
		for mapping := range strings.SplitSeq(tagValue, ",") {
			if parts := strings.Split(mapping, ":"); len(parts) == 2 {
				var providerList = strings.Split(strings.TrimSpace(parts[0]), "|")
				if slices.Contains(providerList, s.provider) {
					value = strings.TrimSpace(parts[1])
					if value == "" {
						// 配置了提供商，但值为空，报错
						err = fmt.Errorf("tag value[%s] is empty for provider[%s]", tagValue, s.provider)
						return
					}
					return
				}
			}
		}
		// 没有找到当前提供商的值
		return
	}
	// 通用值，所有提供商都使用
	value = tagValue
	return
}
