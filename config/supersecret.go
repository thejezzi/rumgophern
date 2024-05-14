package config

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type ContextKey string

const MyKey ContextKey = "trololololololololo"

type MyInnerConfig map[string]map[string]map[string]int

type SuperSecretConfig struct {
	inner *atomic.Pointer[MyInnerConfig]
}

func NewSuperSecretConfig() *SuperSecretConfig {
	initialConfig := make(MyInnerConfig)
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
	for k1, v1 := range m {
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

	if level1, ok := (*config)[key1]; ok {
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
	if config == nil {
		newConfig := make(MyInnerConfig)
		s.inner.Store(&newConfig)
		config = s.inner.Load()
	}

	if _, ok := (*config)[key1]; !ok {
		(*config)[key1] = make(map[string]map[string]int)
	}
	if _, ok := (*config)[key1][key2]; !ok {
		(*config)[key1][key2] = make(map[string]int)
	}

	(*config)[key1][key2][key3] = value
}

func (s *SuperSecretConfig) UpdateValue(key1, key2, key3 string, value int) bool {
	config := s.inner.Load()
	if config == nil {
		return false
	}

	if level1, ok := (*config)[key1]; ok {
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

	if level1, ok := (*config)[key1]; ok {
		if level2, ok := level1[key2]; ok {
			if _, ok := level2[key3]; ok {
				delete(level2, key3)
				if len(level2) == 0 {
					delete(level1, key2)
				}
				if len(level1) == 0 {
					delete(*config, key1)
				}
				return true
			}
		}
	}

	return false
}
