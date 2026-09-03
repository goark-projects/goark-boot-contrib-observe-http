package gbcobservehttp

import (
	"sort"

	"goark.dev/arkarta/servlet"
)

// servletHeaderCarrier 将 Arkarta 标准请求头适配为可观测性传播载体。
type servletHeaderCarrier struct {
	header servlet.Header
}

func (c servletHeaderCarrier) Get(key string) string {
	if c.header == nil {
		return ""
	}
	return c.header.Get(key)
}

func (c servletHeaderCarrier) Set(key, value string) {
	if c.header != nil {
		c.header.Set(key, value)
	}
}

func (c servletHeaderCarrier) Keys() []string {
	if c.header == nil {
		return nil
	}
	seen := make(map[string]struct{})
	keys := make([]string, 0)
	c.header.Visit(func(name, _ string) bool {
		if _, exists := seen[name]; exists {
			return true
		}
		seen[name] = struct{}{}
		keys = append(keys, name)
		return true
	})
	sort.Strings(keys)
	return keys
}
