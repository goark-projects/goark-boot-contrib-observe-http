package gbcobservehttp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"goark.dev/arkarta/servlet"
	arkweb "goark.dev/arkarta/web"
	goweb "goark.dev/goark/web"
	"goark.dev/observe"
	observehttp "goark.dev/observe-http"
	observesdk "goark.dev/observe-sdk"
)

func TestServerFilterPropagatesContextAndRecordsRequest(t *testing.T) {
	t.Parallel()
	exporter := &captureExporter{}
	provider, _ := observesdk.NewProvider(observesdk.WithExporters(exporter))
	instrumenter, _ := observehttp.New(provider)
	registry := goweb.NewRegistry()
	if err := registry.GET("/users/{id}", arkweb.HandlerFunc(func(*arkweb.Context) (arkweb.Result, error) { return nil, nil })); err != nil {
		t.Fatalf("register route: %v", err)
	}
	filter := serverFilter(instrumenter, &routeResolver{registry: registry})
	request, _ := servlet.NewRequest(httptest.NewRequest(http.MethodGet, "http://example.test/users/42", nil))
	response := newTestResponse()
	var traced bool
	err := filter.Filter(t.Context(), request, response, servlet.ChainFunc(func(ctx context.Context, _ *servlet.Request, res servlet.Response) error {
		traced = observe.SpanContextFromContext(ctx).IsValid()
		res.SetStatus(http.StatusCreated)
		return nil
	}))
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if !traced {
		t.Fatal("handler context has no active span")
	}
	if err := provider.ForceFlush(t.Context()); err != nil {
		t.Fatalf("ForceFlush: %v", err)
	}
	exporter.mu.Lock()
	defer exporter.mu.Unlock()
	if len(exporter.spans) != 1 || len(exporter.metrics) < 3 {
		t.Fatalf("signals=spans:%d metrics:%d", len(exporter.spans), len(exporter.metrics))
	}
}

func TestClientInterceptorClonesRequestAndInjectsTraceParent(t *testing.T) {
	t.Parallel()
	exporter := &captureExporter{}
	provider, _ := observesdk.NewProvider(observesdk.WithExporters(exporter))
	instrumenter, _ := observehttp.New(provider)
	interceptor := clientInterceptor(instrumenter)
	request, _ := http.NewRequest(http.MethodGet, "https://example.test/data", nil)
	response, err := interceptor.Intercept(t.Context(), request, func(ctx context.Context, received *http.Request) (*http.Response, error) {
		if received == request {
			t.Fatal("request was not cloned")
		}
		if received.Header.Get("traceparent") == "" {
			t.Fatal("traceparent is missing")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(http.NoBody)}, nil
	})
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("Intercept response=%v err=%v", response, err)
	}
	if request.Header.Get("traceparent") != "" {
		t.Fatal("original request was mutated")
	}
}

type testResponse struct {
	header servlet.Header
	status int
}

func newTestResponse() *testResponse {
	return &testResponse{header: servlet.NewHeader(), status: http.StatusOK}
}
func (r *testResponse) Header() servlet.Header              { return r.header }
func (r *testResponse) SetStatus(status int)                { r.status = status }
func (r *testResponse) Status() int                         { return r.status }
func (*testResponse) Write(value []byte) (int, error)       { return len(value), nil }
func (*testResponse) WriteString(value string) (int, error) { return len(value), nil }
func (*testResponse) Flush() error                          { return nil }
func (*testResponse) Committed() bool                       { return false }
func (*testResponse) Reset() error                          { return nil }
func (*testResponse) BodyWriter() io.Writer                 { return io.Discard }

type captureExporter struct {
	mu      sync.Mutex
	spans   []observe.SpanSnapshot
	metrics []observe.MetricData
}

func (*captureExporter) Descriptor() observe.ExporterDescriptor {
	return observe.ExporterDescriptor{Name: "capture-gbc-http", Signals: observe.SignalTraces | observe.SignalMetrics, Stability: observe.StabilityExperimental, Capabilities: observe.ExporterCapabilities{Push: true, CumulativeTemporality: true, Histogram: true}}
}
func (*captureExporter) ForceFlush(context.Context) error { return nil }
func (*captureExporter) Shutdown(context.Context) error   { return nil }
func (e *captureExporter) ExportSpans(_ context.Context, values []observe.SpanSnapshot) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, value := range values {
		e.spans = append(e.spans, value.Clone())
	}
	return nil
}
func (e *captureExporter) ExportMetrics(_ context.Context, values []observe.MetricData) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, value := range values {
		e.metrics = append(e.metrics, value.Clone())
	}
	return nil
}
