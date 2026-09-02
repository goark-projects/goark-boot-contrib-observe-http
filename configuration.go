package gbcobservehttp

import (
	"context"

	"goark.dev/boot"
	gbcobserve "goark.dev/gbc-observe"
	gbcweb "goark.dev/gbc-web"
	goarkcontainer "goark.dev/goark/container"
	appcontext "goark.dev/goark/context"
	goweb "goark.dev/goark/web"
	webclient "goark.dev/goark/web/client"
	"goark.dev/observe"
	observehttp "goark.dev/observe-http"
)

// AutoConfigure 创建 HTTP 观测自动配置，并在缺失时补充基础可观测配置。
func AutoConfigure(options ...Option) boot.AutoConfiguration {
	copied := append([]Option(nil), options...)
	return boot.NewAutoConfiguration(StarterID, func(ctx context.Context, app *appcontext.ApplicationContext) error {
		if !hasConfiguration(app, gbcobserve.StarterID+".configuration") {
			if err := gbcobserve.AutoConfigure().Configure(ctx, app); err != nil {
				return err
			}
		}
		return app.RegisterConfiguration(configuration{options: copied})
	}, boot.WithAutoConfigurationOrder(100))
}

type configuration struct{ options []Option }

func (configuration) Name() string { return StarterID + ".configuration" }
func (configuration) Order() int   { return -9000 }
func (c configuration) Register(ctx context.Context, registry *goarkcontainer.Registry) error {
	return c.RegisterWithContext(ctx, appcontext.NewConfigurationContext(nil, registry))
}
func (c configuration) RegisterWithContext(_ context.Context, config appcontext.ConfigurationContext) error {
	resolved, err := newSettings(config.Environment(), c.options)
	if err != nil {
		return err
	}
	if !*resolved.serverEnabled && !*resolved.clientEnabled {
		return nil
	}
	if err := goarkcontainer.Register[*observehttp.Instrumenter](config.Registry(), BeanNameInstrumenter, func(ctx context.Context, resolver goarkcontainer.Resolver) (*observehttp.Instrumenter, error) {
		provider, err := goarkcontainer.Get[observe.Provider](ctx, resolver, gbcobserve.BeanNameProvider)
		if err != nil {
			return nil, err
		}
		return observehttp.New(provider)
	}, goarkcontainer.WithDependsOn(gbcobserve.BeanNameProvider)); err != nil {
		return err
	}
	if *resolved.serverEnabled {
		if err := registerServer(config.Registry()); err != nil {
			return err
		}
	}
	if *resolved.clientEnabled {
		return registerClient(config.Registry())
	}
	return nil
}
func registerServer(registry *goarkcontainer.Registry) error {
	return goarkcontainer.Register[goweb.Configurer](registry, BeanNameServerConfigurer, func(ctx context.Context, resolver goarkcontainer.Resolver) (goweb.Configurer, error) {
		instrumenter, err := goarkcontainer.Get[*observehttp.Instrumenter](ctx, resolver, BeanNameInstrumenter)
		if err != nil {
			return nil, err
		}
		return goweb.ConfigurerFunc(func(_ context.Context, webRegistry *goweb.Registry) error {
			webRegistry.AddFilter(serverFilter(instrumenter, &routeResolver{registry: webRegistry}))
			return nil
		}), nil
	}, goarkcontainer.WithDependsOn(BeanNameInstrumenter))
}
func registerClient(registry *goarkcontainer.Registry) error {
	return goarkcontainer.Register[gbcweb.HTTPClientBuilderCustomizer](registry, BeanNameClientCustomizer, func(ctx context.Context, resolver goarkcontainer.Resolver) (gbcweb.HTTPClientBuilderCustomizer, error) {
		instrumenter, err := goarkcontainer.Get[*observehttp.Instrumenter](ctx, resolver, BeanNameInstrumenter)
		if err != nil {
			return nil, err
		}
		return gbcweb.HTTPClientBuilderCustomizerFunc(func(_ context.Context, builder *webclient.Builder) (*webclient.Builder, error) {
			return builder.Apply(webclient.WithInterceptor(clientInterceptor(instrumenter))), nil
		}), nil
	}, goarkcontainer.WithDependsOn(BeanNameInstrumenter))
}
func hasConfiguration(app *appcontext.ApplicationContext, name string) bool {
	for _, descriptor := range app.Configurations() {
		if descriptor.Name == name {
			return true
		}
	}
	return false
}
