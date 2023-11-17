package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
}

func (c *Client) Browse(url string) {

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	var body []byte

	tr := otel.Tracer("sonicstore/client")

	ctx, span := tr.Start(context.Background(), "web request", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.KeyValue(semconv.HTTPMethod("GET")),
		attribute.KeyValue(semconv.NetProtocolVersion("1.1")),
		attribute.KeyValue(semconv.NetPeerName("example.com")),
		attribute.KeyValue(semconv.NetSockPeerAddr("101.10.9.5")),
	)
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	fmt.Printf("Sending request...\n")
	span.AddEvent("about to send a request")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	span.AddEvent("request sent", trace.WithAttributes(attribute.String("url", url)))
	span.SetAttributes(
		attribute.Int("status.code", res.StatusCode),
	)
	body, err = io.ReadAll(res.Body)
	_ = res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response Received: %s\n\n\n", body)
	fmt.Printf("Waiting for few seconds to export spans ...\n\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("Inspect traces on stdout\n")
}
