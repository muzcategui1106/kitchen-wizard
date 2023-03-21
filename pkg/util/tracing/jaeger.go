package tracing

import (
	"context"

	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaegerTracer(ctx context.Context, collectorAddress string) error {
	// Recommended configuration for production.
	logger.Log.Sugar().Infof("setting up opentracing global tracer to send traaces to %v", collectorAddress)
	cfg := jaegercfg.Configuration{
		ServiceName: "kitchen-wizard-api",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: collectorAddress,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Extractor(opentracing.HTTPHeaders, propagator),
		jaegercfg.Injector(opentracing.HTTPHeaders, propagator),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)

	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	go func() {
		<-ctx.Done()
		logger.Log.Info("closing open tracer")
		closer.Close()
	}()

	return nil
}
