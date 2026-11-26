package observability

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", serviceName),
	))
	if err != nil {
		return nil, err
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		provider := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)
		otel.SetTracerProvider(provider)
		return provider.Shutdown, nil
	}

	options := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" {
		options = append(options, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}

func TracingMiddleware(serviceName string) echo.MiddlewareFunc {
	tracer := otel.Tracer(serviceName)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := tracer.Start(c.Request().Context(), c.Request().Method+" "+c.Request().URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", c.Request().Method),
					attribute.String("http.route", c.Request().URL.Path),
				),
			)
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			err := next(c)
			span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
			if err != nil {
				span.RecordError(err)
			}
			return err
		}
	}
}
