package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otel-bug/pkg"
)

var counter syncint64.Counter

func main() {
	ctx := context.Background()

	provider, err := pkg.InitMetrics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	meter := provider.Meter("name")

	counter, err = meter.SyncInt64().Counter("some_counter")
	if err != nil {
		log.Fatal(err)
	}

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(meter)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics(meter metric.Meter) {
	log.Printf("serving metrics at localhost:8080/metrics")
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		histogram, err := meter.SyncInt64().Histogram("http_client_total_duration")
		if err != nil {
			log.Fatal(err)
		}
		histogram.Record(request.Context(), 23, attribute.String("to_service", "catalog"))
		histogram.Record(request.Context(), 187, attribute.String("to_service", "ordering"))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
