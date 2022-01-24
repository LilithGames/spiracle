package telemetry

import (
	"log"
	"net/http"
	"time"
	"context"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/attribute"
)

func InitMeter(ctx context.Context) {
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.Key("service_name").String("spiracle"),
	))
	if err != nil {
		log.Panicf("InitMeter resource.New err: %v", err)
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(time.Second * 5),
		controller.WithResource(res),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())

	http.HandleFunc("/metrics", exporter.ServeHTTP)
}

func RunMetricServer(ctx context.Context, addr string) {
	log.Printf("Prometheus server running on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("metric server err: ", err)
	}
}
