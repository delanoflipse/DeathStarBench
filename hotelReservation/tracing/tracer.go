package tracing

import (
	"context"
	"errors"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// var (
// 	defaultSampleRatio float64 = 0.01
// )

// Init returns a newly configured tracer
func Init(serviceName, host string) (*trace.TracerProvider, error) {

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
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost(),         // Discover and provide host information.
	)

	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		log.Println(err) // Log non-fatal issues.
	} else if err != nil {
		log.Fatalln(err) // The error may be fatal.
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSyncer(traceExporter),
		trace.WithSampler(trace.AlwaysSample()),
	)

	return tracerProvider, nil
}
