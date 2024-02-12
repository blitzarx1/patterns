package main

import (
	pkgContext "context"
	"log"
	"os"

	"github.com/boson-research/patterns/internal"
	"github.com/boson-research/patterns/internal/alphabet"
	"github.com/boson-research/patterns/internal/context"
	"github.com/boson-research/patterns/internal/processor"
	"github.com/boson-research/patterns/internal/telemetry"
)

const (
	alphabetPath = "input/alphabet"
	textPath     = "input/text"
)

func main() {
	ctx, closer, err := telemetry.Init(context.New(pkgContext.Background()), telemetry.Config{
		Name:               "patterns",
		Version:            internal.MustGetGitVersion(),
		JaegerOTLPEndpoint: "jaeger:4318",
	})
	if err != nil {
		log.Fatalf("failed to initialize tracing: %s", err)
	}
	defer closer(ctx)

	ctx, span := ctx.StartSpan("main")
	defer span.End()

	alphabetRaw, err := os.ReadFile(alphabetPath)
	if err != nil {
		log.Fatalf("failed to read alphabet: %s", err)
	}

	ctx.Logger().Info("alphabet loaded")

	text, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatalf("failed to read text: %s", err)
	}

	ctx.Logger().Info("text loaded")

	p := processor.New(ctx)
	p.AnalyzeAlphabet(ctx, alphabet.Alphabet(alphabetRaw))
	p.AnalyzeText(ctx, text)
}
