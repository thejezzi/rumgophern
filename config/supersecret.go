package config

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type (
	ContextKey string
	SpecialMap map[string]map[string]map[string]int
)

const MyKey ContextKey = "trololololololololo"

type MyInnerConfig struct {
	InnerMap SpecialMap
}

type SuperSecretConfig struct {
	inner *atomic.Pointer[MyInnerConfig]
}

func NewSuperSecretConfig() *SuperSecretConfig {
	initialMap := make(SpecialMap)
	initialConfig := MyInnerConfig{
		InnerMap: initialMap,
	}
	configPointer := &atomic.Pointer[MyInnerConfig]{}
	configPointer.Store(&initialConfig)
	return &SuperSecretConfig{
		inner: configPointer,
	}
}

func (s *SuperSecretConfig) Unwrap() *MyInnerConfig {
	return s.inner.Load()
}

func (m MyInnerConfig) String() string {
	sb := strings.Builder{}
	for k1, v1 := range m.InnerMap {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				sb.WriteString(fmt.Sprintf("%s/%s/%s: %d\n", k1, k2, k3, v3))
			}
		}
	}
	return sb.String()
}

func (s *SuperSecretConfig) GetValue(key1, key2, key3 string) (int, bool) {
	config := s.inner.Load()
	if config == nil {
		return 0, false
	}

	if level1, ok := (*config).InnerMap[key1]; ok {
		if level2, ok := level1[key2]; ok {
			if value, ok := level2[key3]; ok {
				return value, true
			}
		}
	}

	return 0, false
}

func (s *SuperSecretConfig) SetValue(key1, key2, key3 string, value int) {
	config := s.inner.Load()

	newConfig := config.DeepCopy()
	if newConfig == nil {
		panic("config is nil")
	}

	if _, ok := (*newConfig).InnerMap[key1]; !ok {
		(*newConfig).InnerMap[key1] = make(map[string]map[string]int)
	}

	if _, ok := (*newConfig).InnerMap[key1][key2]; !ok {
		(*newConfig).InnerMap[key1][key2] = make(map[string]int)
	}

	(*newConfig).InnerMap[key1][key2][key3] = value
	s.inner.Store(newConfig)
}

func (m *MyInnerConfig) DeepCopy() *MyInnerConfig {
	newConfig := MyInnerConfig{
		InnerMap: make(SpecialMap),
	}
	for k1, v1 := range m.InnerMap {
		newConfig.InnerMap[k1] = make(map[string]map[string]int)
		for k2, v2 := range v1 {
			newConfig.InnerMap[k1][k2] = make(map[string]int)
			for k3, v3 := range v2 {
				newConfig.InnerMap[k1][k2][k3] = v3
			}
		}
	}

	return &newConfig
}

func (s *SuperSecretConfig) UpdateValue(key1, key2, key3 string, value int) bool {
	config := s.inner.Load()
	if config == nil {
		return false
	}

	if level1, ok := (*config).InnerMap[key1]; ok {
		if level2, ok := level1[key2]; ok {
			if _, ok := level2[key3]; ok {
				level2[key3] = value
				return true
			}
		}
	}

	return false
}

func (s *SuperSecretConfig) DeleteValue(key1, key2, key3 string) bool {
	config := s.inner.Load()
	if config == nil {
		return false
	}

	if level1, ok := (*config).InnerMap[key1]; ok {
		if level2, ok := level1[key2]; ok {
			if _, ok := level2[key3]; ok {
				delete(level2, key3)
				if len(level2) == 0 {
					delete(level1, key2)
				}
				if len(level1) == 0 {
					delete((*config).InnerMap, key1)
				}
				return true
			}
		}
	}

	return false
}
