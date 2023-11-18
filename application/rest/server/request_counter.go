package server

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("sonic-service-meter")

// Configures counter instrument for counter incoming requests.
func meterCounter() metric.Int64Counter {

	apiCounter, err := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	return apiCounter
}
