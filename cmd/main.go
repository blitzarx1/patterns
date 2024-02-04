package main

import (
	"context"
	"log"
	"os"

	"github.com/boson-research/patterns/internal/models"
	"github.com/boson-research/patterns/internal/processor"
	"github.com/boson-research/patterns/internal/telemetry"
	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

const (
	alphabetPath = "input/alphabet"
	textPath     = "input/text"
)

func main() {
	ctx, closer, err := telemetry.Init(context.Background(), telemetry.Config{
		Name:               "patterns",
		Version:            "v0.1.1", // todo take from git tag + commit hash
		JaegerOTLPEndpoint: "jaeger:4318",
	})
	if err != nil {
		log.Fatalf("failed to initialize tracing: %s", err)
	}
	defer closer(ctx)

	ctx, span := otel.Tracer("main").Start(ctx, "main")
	defer span.End()
	l := logger.Logger(ctx)

	alphabetRaw, err := os.ReadFile(alphabetPath)
	if err != nil {
		log.Fatalf("failed to read alphabet: %s", err)
	}

	l.Info("alphabet loaded")

	text, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatalf("failed to read text: %s", err)
	}

	l.Info("text loaded")

	p := processor.NewProcessor(ctx, models.Alphabet(alphabetRaw))

	patternsLocation := p.PatternsLocations(ctx, text)

	l.Infof("found patterns locations: %v", patternsLocation)
}
