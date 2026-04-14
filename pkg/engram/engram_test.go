package engram

import (
	"context"
	"encoding/json"
	"testing"

	sdkengram "github.com/bubustack/bubu-sdk-go/engram"
	cfgpkg "github.com/bubustack/text-emitter-engram/pkg/config"
)

func TestProcessEmitsConfiguredMessage(t *testing.T) {
	engine := New()
	if err := engine.Init(context.Background(), cfgpkg.Config{Message: "hello"}, nil); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	res, err := engine.Process(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}
	data, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Data)
	}
	if data["type"] != "text.emit.v1" {
		t.Fatalf("expected type text.emit.v1, got %#v", data["type"])
	}
	if data["text"] != "hello" {
		t.Fatalf("expected text hello, got %#v", data["text"])
	}
}

func TestProcessDisabledReturnsEmptyMessage(t *testing.T) {
	enabled := false
	engine := New()
	if err := engine.Init(context.Background(), cfgpkg.Config{
		Message: "hello",
		Enabled: &enabled,
	}, nil); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	res, err := engine.Process(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}
	data, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Data)
	}
	if data["text"] != "" {
		t.Fatalf("expected empty text when disabled, got %#v", data["text"])
	}
}

func TestStreamEmitsOnce(t *testing.T) {
	engine := New()
	if err := engine.Init(context.Background(), cfgpkg.Config{Message: "hello"}, nil); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	in := make(chan sdkengram.InboundMessage, 1)
	out := make(chan sdkengram.StreamMessage, 2)
	in <- sdkengram.NewInboundMessage(sdkengram.StreamMessage{})
	close(in)

	if err := engine.Stream(context.Background(), in, out); err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected one emitted message, got %d", len(out))
	}
	msg := <-out
	var decoded map[string]any
	if err := json.Unmarshal(msg.Payload, &decoded); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if decoded["text"] != "hello" {
		t.Fatalf("expected streamed text hello, got %#v", decoded["text"])
	}
}
