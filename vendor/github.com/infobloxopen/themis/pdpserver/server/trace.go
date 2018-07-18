package server

import (
	"strings"

	ot "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

func initTracing(tracingType, tracingEP string) (ot.Tracer, error) {
	if tracingEP == "" {
		return nil, nil
	}

	switch strings.ToLower(tracingType) {
	default:
		return nil, newTracingTypeError(tracingType)

	case "zipkin":
		return setupZipkin(tracingEP)
	}
}

func setupZipkin(tracingEP string) (ot.Tracer, error) {
	if strings.Index(tracingEP, "http") == -1 {
		tracingEP = "http://" + tracingEP + "/api/v1/spans"
	}

	collector, err := zipkin.NewHTTPCollector(tracingEP)
	if err != nil {
		return nil, err
	}

	recorder := zipkin.NewRecorder(collector, false, "", "PDP")
	return zipkin.NewTracer(recorder)
}
