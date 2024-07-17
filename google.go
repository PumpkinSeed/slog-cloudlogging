package slogcloudlogging

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultGoogleAutoFlushInterval = 500
)

var (
	ErrUninitializedLogger = errors.New("uninitialized logger error")
)

type Google struct {
	ProjectID         string
	LogName           string
	AutoFlushInterval int
	client            *logging.Client
	logger            *logging.Logger

	Handler                slog.Handler
	ForwardHandler         bool
	UseOpenTelemetryTracer bool
	TracePrefix            string
}

type Line struct {
	Level     slog.Level             `json:"level,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Time      string                 `json:"time,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func (g *Google) Print(ctx context.Context, main Line) {
	// Make sure we have a set google client
	g.init()

	var severity logging.Severity
	switch main.Level {
	case slog.LevelError:
		severity = logging.Error
	case slog.LevelWarn:
		severity = logging.Warning
	case slog.LevelInfo:
		severity = logging.Info
	case slog.LevelDebug:
		severity = logging.Debug
	}

	// Create payload
	payload := make(map[string]interface{})
	for k, v := range main.Data {
		payload[k] = v
	}
	payload["timestamp"] = main.Timestamp
	payload["time"] = main.Time
	for fieldKey, fieldValue := range main.Data {
		if v, ok := fieldValue.(string); ok && v != "" {
			payload[fieldKey] = v
		}
	}

	entry := logging.Entry{
		Severity: severity,
		Payload:  payload,
	}

	if g.UseOpenTelemetryTracer {
		if s := trace.SpanContextFromContext(ctx); s.IsValid() {
			entry.Trace = g.TracePrefix + s.TraceID().String()
			entry.SpanID = s.SpanID().String()
			entry.TraceSampled = s.TraceFlags().IsSampled()
		}
	}

	// Adds an entry to the log buffer.
	g.logger.Log(entry)
}

func (g *Google) AutoFlush() chan bool {
	g.init()
	if g.AutoFlushInterval == 0 {
		g.AutoFlushInterval = DefaultGoogleAutoFlushInterval
	}
	ticker := time.NewTicker(time.Duration(g.AutoFlushInterval) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				if g.logger != nil {
					if err := g.logger.Flush(); err != nil {
						panic(err)
					}
				}
			}
		}
	}()

	return done
}

func (g *Google) Flush() error {
	if g != nil && g.logger != nil {
		return g.logger.Flush()
	}
	return ErrUninitializedLogger
}

func (g *Google) init() {
	if g.client == nil || g.logger == nil {
		ctx := context.Background()

		// Creates a client.
		var err error
		g.client, err = logging.NewClient(ctx, g.ProjectID)
		if err != nil {
			panic(err)
		}

		g.logger = g.client.Logger(g.LogName)
	}
}
