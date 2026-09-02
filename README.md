# Goark Boot Observe HTTP

[简体中文](README.zh-CN.md)

`goark.dev/gbc-observe-http` auto-configures HTTP server and client instrumentation for Goark Boot. Protocol semantics remain in `goark.dev/observe-http`; this module only adapts them to Goark Web and Boot configuration.

## Usage

```go
boot.WithAutoConfiguration(
    gbcobserve.AutoConfigure(gbcobserve.WithExporters(exporter)),
    gbcobservehttp.AutoConfigure(),
    gbcweb.AutoConfigure(),
)
```

The starter registers a Servlet Filter that passes extracted context into the handler chain and a client interceptor that clones requests before injecting headers. Registered MVC route templates are resolved once and used as bounded labels; unmatched paths become `unknown`.

`goark.observe.http.enabled`, `.server.enabled`, and `.client.enabled` default to `true`. The base observe starter is installed when absent. Disabling base observability preserves HTTP behavior through the no-op provider.

Licensed under Apache License 2.0.
