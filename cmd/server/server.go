package main

import (
	"context"
	"log"

	"github.com/bersennaidoo/sonicstore/application/rest/server"
	"github.com/bersennaidoo/sonicstore/physical/opentel"
)

func main() {
	tp, err := opentel.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	srv := server.Server{}
	srv.Run()
}
