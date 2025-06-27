package tracing

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// var (
// 	defaultSampleRatio float64 = 0.01
// )

// Init returns a newly configured tracer
func Init(serviceName, host string) (*sdktrace.TracerProvider, error) {

	// ratio := defaultSampleRatio
	// if val, ok := os.LookupEnv("JAEGER_SAMPLE_RATIO"); ok {
	// 	ratio, _ = strconv.ParseFloat(val, 64)
	// 	if ratio > 1 {
	// 		ratio = 1.0
	// 	}
	// }

	// log.Info().Msgf("Jaeger client: adjusted sample ratio %f", ratio)
	// tempCfg := &config.Configuration{
	// 	ServiceName: serviceName,
	// 	Sampler: &config.SamplerConfig{
	// 		Type:  "probabilistic",
	// 		Param: ratio,
	// 	},
	// 	Reporter: &config.ReporterConfig{
	// 		LogSpans:            false,
	// 		BufferFlushInterval: 1 * time.Second,
	// 		LocalAgentHostPort:  host,
	// 	},
	// }

	// log.Info().Msg("Overriding Jaeger config with env variables")
	// cfg, err := tempCfg.FromEnv()

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create OpenTelemetry exporter: %v", err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return provider, nil
}
