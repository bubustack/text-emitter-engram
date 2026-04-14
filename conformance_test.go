package main

import (
	"testing"

	"github.com/bubustack/bubu-sdk-go/conformance"
	cfgpkg "github.com/bubustack/text-emitter-engram/pkg/config"
	"github.com/bubustack/text-emitter-engram/pkg/engram"
)

func TestConformance(t *testing.T) {
	suite := conformance.BatchSuite[cfgpkg.Config, map[string]any]{
		Engram: engram.New(),
		Config: cfgpkg.Config{},
		Inputs: map[string]any{},
	}
	suite.Run(t)
}
