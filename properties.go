package gbcobservehttp

const (
	StarterID                = "goark.boot.observe.http"
	BeanNameInstrumenter     = "goark.observe.http.instrumenter"
	BeanNameServerConfigurer = "goark.observe.http.serverConfigurer"
	BeanNameClientCustomizer = "goark.observe.http.clientCustomizer"
)

const (
	PropertyEnabled       = "goark.observe.http.enabled"
	PropertyServerEnabled = "goark.observe.http.server.enabled"
	PropertyClientEnabled = "goark.observe.http.client.enabled"
)

const DefaultEnabled = true
