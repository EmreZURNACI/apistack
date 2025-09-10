package main

import (
	"context"
	"log"
	"time"

	"github.com/EmreZURNACI/apistack/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	config := zap.NewProductionConfig()                                         // or zap.NewDevelopmentConfig() or any other zap.Config
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339) // or time.RubyDate or "2006-01-02 15:04:05" or even freaking time.Kitchen
	logger, err := config.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	zap.L().Info("application starting...")
}

func main() {
	// Alınan traceler fonksiyonşardan geçirilecek
	tp := initTracer("stackapi")
	defer func() {
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
	}()

	server.Route()
}

func initTracer(service_name string) *sdktrace.TracerProvider {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(service_name),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

//Golang, functionları otomatik olarak trace etmeye izin vermez.Manuel olarka yapılandırılması lazım.
//fiber'dakicontextten gelen ctx'i userContext olarak almak gereklidir.

// ! Create Actor commit veya rollback trace bozuyor ve spanler ayrı
