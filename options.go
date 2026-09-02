package gbcobservehttp

// Option 定制 HTTP 观测自动配置。
type Option func(*settings)

// WithEnabled 同时控制服务端和客户端 HTTP 观测。
func WithEnabled(enabled bool) Option {
	return func(settings *settings) { settings.enabled = &enabled }
}

// WithServerEnabled 控制服务端 HTTP 观测。
func WithServerEnabled(enabled bool) Option {
	return func(settings *settings) { settings.serverEnabled = &enabled }
}

// WithClientEnabled 控制客户端 HTTP 观测。
func WithClientEnabled(enabled bool) Option {
	return func(settings *settings) { settings.clientEnabled = &enabled }
}
