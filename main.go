package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otel-bug/pkg"
)

var histogram syncint64.Histogram

func main() {
	ctx := context.Background()

	provider, err := pkg.InitMetrics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	meter := provider.Meter("name")

	histogram, err = meter.SyncInt64().Histogram("http.client.total.duration")
	if err != nil {
		log.Fatal(err)
	}

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:8080/metrics")
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		histogram.Record(request.Context(), 23, attribute.String("service", "catalog"))
		histogram.Record(request.Context(), 187, attribute.String("service", "ordering"))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
