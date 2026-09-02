package gbcobservehttp

import (
	"context"
	"net/http"

	webclient "goark.dev/goark/web/client"
	observehttp "goark.dev/observe-http"
)

func clientInterceptor(instrumenter *observehttp.Instrumenter) webclient.Interceptor {
	return webclient.InterceptorFunc(func(ctx context.Context, req *http.Request, next webclient.ExchangeFunc) (*http.Response, error) {
		if req == nil {
			return next(ctx, req)
		}
		cloned := req.Clone(ctx)
		cloned.Header = req.Header.Clone()
		scheme, host, protocol := "", "", ""
		if cloned.URL != nil {
			scheme = cloned.URL.Scheme
			host = cloned.URL.Host
		}
		protocol = cloned.Proto
		ctx, operation := instrumenter.StartClient(ctx, observehttp.HeaderCarrier(cloned.Header), observehttp.ClientRequest{Method: cloned.Method, Scheme: scheme, Host: host, Protocol: protocol})
		cloned = cloned.WithContext(ctx)
		response, err := next(ctx, cloned)
		status := 0
		if response != nil {
			status = response.StatusCode
		}
		operation.End(status, err)
		return response, err
	})
}
