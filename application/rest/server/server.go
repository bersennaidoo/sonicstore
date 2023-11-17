package server

import (
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type Server struct{}

func (s *Server) Welcome(w http.ResponseWriter, req *http.Request) {
	//tr := otel.Tracer("sonicstore/server")
	ctx := req.Context()
	span := trace.SpanFromContext(ctx)

	//_, span := tr.Start(context.Background(), "welcome", trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.KeyValue(semconv.HTTPMethod("GET")),
		attribute.KeyValue(semconv.NetProtocolVersion("1.1")),
		attribute.KeyValue(semconv.NetPeerName("example.com")),
		attribute.KeyValue(semconv.NetSockPeerAddr("101.10.9.5")),
	)
	defer span.End()

	_, _ = io.WriteString(w, "Welcome to sonic store!\n")
}

func (s *Server) Run() {

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(s.Welcome), "Welcome")

	http.Handle("/", otelHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
