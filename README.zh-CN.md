# Goark Boot Observe HTTP

[English](README.md)

`goark.dev/gbc-observe-http` 为 Goark Boot 自动配置 HTTP 服务端和客户端观测。HTTP 协议语义由 `goark.dev/observe-http` 维护，本模块只负责 Goark Web 与 Boot 配置适配。

## 使用

```go
boot.WithAutoConfiguration(
    gbcobserve.AutoConfigure(gbcobserve.WithExporters(exporter)),
    gbcobservehttp.AutoConfigure(),
    gbcweb.AutoConfigure(),
)
```

starter 注册 Servlet Filter，把提取后的 context 传入处理链；同时注册客户端拦截器，在注入 header 前克隆请求。已注册 MVC 路由模板只解析一次并作为有界标签，未命中路径统一使用 `unknown`。

`goark.observe.http.enabled`、`.server.enabled`、`.client.enabled` 默认均为 `true`。基础 observe starter 缺失时自动补充。禁用基础观测后，HTTP 行为通过 no-op Provider 保持不变。

本项目采用 Apache License 2.0。
