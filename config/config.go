package config

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"working-project/common/kms"

	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct {
	kmsManager *kms.Manager
}

// NewLoader 创建配置加载器
// kmsManager: KMS管理器，用于解密配置中的敏感信息，如果为nil则不解密
func NewLoader(kmsManager *kms.Manager) *Loader {
	return &Loader{
		kmsManager: kmsManager,
	}
}

// LoadFromFile 从YAML文件加载配置，并自动解密KMS加密的字段
// filePath: 配置文件路径
// config: 配置结构体指针
func (l *Loader) LoadFromFile(ctx context.Context, filePath string, config interface{}) error {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	// 解析YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	// 如果配置了KMS，则解密配置中的敏感字段
	if l.kmsManager != nil {
		if err := l.decryptConfig(ctx, config); err != nil {
			return fmt.Errorf("failed to decrypt config: %w", err)
		}
	}

	return nil
}

// decryptConfig 递归解密配置结构体中所有以KMS前缀开头的字符串字段
func (l *Loader) decryptConfig(ctx context.Context, config interface{}) error {
	return l.decryptValue(ctx, reflect.ValueOf(config))
}

// decryptValue 递归处理反射值
func (l *Loader) decryptValue(ctx context.Context, v reflect.Value) error {
	// 如果是指针，获取指向的值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		// 遍历结构体的所有字段
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanSet() {
				continue
			}
			if err := l.decryptValue(ctx, field); err != nil {
				return err
			}
		}

	case reflect.Slice, reflect.Array:
		// 遍历数组/切片的所有元素
		for i := 0; i < v.Len(); i++ {
			if err := l.decryptValue(ctx, v.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		// 遍历map的所有值
		iter := v.MapRange()
		for iter.Next() {
			mapValue := iter.Value()
			mapKey := iter.Key()

			// map的value可能是字符串、结构体或指针
			if mapValue.Kind() == reflect.String {
				// 如果是字符串，检查是否需要解密
				strValue := mapValue.String()
				decrypted, err := l.kmsManager.DecryptIfNeeded(ctx, strValue)
				if err != nil {
					return fmt.Errorf("failed to decrypt map value: %w", err)
				}
				// 如果解密后的值不同，更新map
				if decrypted != strValue {
					v.SetMapIndex(mapKey, reflect.ValueOf(decrypted))
				}
			} else if mapValue.Kind() == reflect.Ptr {
				// 如果是指针，递归处理指向的值
				if err := l.decryptValue(ctx, mapValue); err != nil {
					return err
				}
			} else if mapValue.Kind() == reflect.Struct {
				// 如果是结构体，需要创建可设置的副本进行解密，然后更新map
				// 因为map的value不可直接设置，需要创建新的结构体
				structValue := reflect.New(mapValue.Type()).Elem()
				structValue.Set(mapValue)

				// 递归处理结构体的字段
				if err := l.decryptValue(ctx, structValue); err != nil {
					return err
				}

				// 更新map中的值
				v.SetMapIndex(mapKey, structValue)
			}
		}

	case reflect.String:
		// 字符串类型，检查是否需要解密
		if v.CanSet() {
			strValue := v.String()
			decrypted, err := l.kmsManager.DecryptIfNeeded(ctx, strValue)
			if err != nil {
				return fmt.Errorf("failed to decrypt string field: %w", err)
			}
			if decrypted != strValue {
				v.SetString(decrypted)
			}
		}
	}

	return nil
}

// MustLoadFromFile 从文件加载配置，如果失败则panic
func (l *Loader) MustLoadFromFile(ctx context.Context, filePath string, config interface{}) {
	if err := l.LoadFromFile(ctx, filePath, config); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}
