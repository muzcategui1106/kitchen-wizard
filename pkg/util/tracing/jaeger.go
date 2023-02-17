package tracing

import (
	"context"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaegerTracer(ctx context.Context) error {
	// Recommended configuration for production.
	logger.Log.Info("setting up opentracing global tracer")
	cfg := jaegercfg.Configuration{}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"collector-collector.observability.svc:14268",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		closer.Close()
	}()

	return nil
}
