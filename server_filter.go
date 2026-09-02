package gbcobservehttp

import (
	"context"

	"goark.dev/arkarta/servlet"
	"goark.dev/observe-http"
)

func serverFilter(instrumenter *observehttp.Instrumenter, routes *routeResolver) servlet.Filter {
	return servlet.FilterFunc(func(ctx context.Context, req *servlet.Request, res servlet.Response, chain servlet.Chain) error {
		if chain == nil {
			return servlet.ErrNilHandler
		}
		if req == nil || res == nil {
			return chain.Next(ctx, req, res)
		}
		route := routes.Resolve(req.Method(), req.Path())
		ctx, operation := instrumenter.StartServer(ctx, observehttp.HeaderCarrier(req.Header()), observehttp.ServerRequest{Method: req.Method(), Route: route, Scheme: req.Scheme(), Host: req.Host(), Protocol: req.Protocol()})
		err := chain.Next(ctx, req, res)
		operation.End(res.Status(), err)
		return err
	})
}
