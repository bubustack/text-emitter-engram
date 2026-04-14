package engram

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	sdk "github.com/bubustack/bubu-sdk-go"
	sdkengram "github.com/bubustack/bubu-sdk-go/engram"
	cfgpkg "github.com/bubustack/text-emitter-engram/pkg/config"
)

const componentName = "text-emitter-engram"

// TextEmitterEngram emits a single text message per stream instance.
type TextEmitterEngram struct {
	cfg  cfgpkg.Config
	sent bool
}

func New() *TextEmitterEngram { return &TextEmitterEngram{} }

func (e *TextEmitterEngram) Init(_ context.Context, cfg cfgpkg.Config, _ *sdkengram.Secrets) error {
	e.cfg = cfgpkg.Normalize(cfg)
	return nil
}

func (e *TextEmitterEngram) Process(
	ctx context.Context,
	_ *sdkengram.ExecutionContext,
	_ map[string]any,
) (*sdkengram.Result, error) {
	logger := sdk.LoggerFromContext(ctx).With(
		"component", componentName,
		"mode", "batch",
	)

	if !e.canEmit() {
		return sdkengram.NewResultFrom(map[string]any{
			"type": "text.emit.v1",
			"text": "",
		}), nil
	}

	if err := waitDelay(ctx, e.cfg.DelayMs); err != nil {
		return nil, err
	}

	logger.Info("Text emitted", "message", truncateText(e.cfg.Message, 160))
	return sdkengram.NewResultFrom(map[string]any{
		"type": "text.emit.v1",
		"text": e.cfg.Message,
	}), nil
}

func (e *TextEmitterEngram) Stream(
	ctx context.Context,
	in <-chan sdkengram.InboundMessage,
	out chan<- sdkengram.StreamMessage,
) error {
	logger := sdk.LoggerFromContext(ctx).With(
		"component", componentName,
		"mode", "stream",
	)

	if e.shouldEmit() {
		if err := e.emitText(ctx, out, logger); err != nil {
			return err
		}
	}

	// Drain input so the connector stream does not back up.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-in:
			if !ok {
				return nil
			}
			msg.Done()
		}
	}
}

func (e *TextEmitterEngram) shouldEmit() bool {
	if e.sent {
		return false
	}
	return e.canEmit()
}

func (e *TextEmitterEngram) canEmit() bool {
	if e.cfg.Enabled != nil && !*e.cfg.Enabled {
		return false
	}
	if strings.TrimSpace(e.cfg.Message) == "" {
		return false
	}
	return true
}

func (e *TextEmitterEngram) emitText(
	ctx context.Context,
	out chan<- sdkengram.StreamMessage,
	logger *slog.Logger,
) error {
	if err := waitDelay(ctx, e.cfg.DelayMs); err != nil {
		return err
	}

	payload := map[string]any{
		"type": "text.emit.v1",
		"text": e.cfg.Message,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	metadata := map[string]string{
		"type":     "text.emit.v1",
		"provider": "text-emitter",
	}

	select {
	case out <- sdkengram.StreamMessage{
		Metadata: metadata,
		Inputs:   bytes,
		Payload:  bytes,
	}:
		e.sent = true
		logger.Info("Text emitted", "message", truncateText(e.cfg.Message, 160))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func waitDelay(ctx context.Context, delayMs int) error {
	if delayMs <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(delayMs) * time.Millisecond):
		return nil
	}
}

func truncateText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "…"
}
