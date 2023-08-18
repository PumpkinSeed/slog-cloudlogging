package slogcloudlogging

import (
	"context"
	"log/slog"
	"time"
)

var _ slog.Handler = &Google{}

type Opts struct {
	Handler           slog.Handler
	ForwardHandler    bool
	AutoFlushInterval int
}

func NewHandler(project string, logName string, opts *Opts) *Google {
	g := &Google{
		ProjectID: project,
		LogName:   logName,
	}
	if opts != nil {
		g.Handler = opts.Handler
		g.AutoFlushInterval = opts.AutoFlushInterval
		g.ForwardHandler = opts.ForwardHandler
	}
	g.init()
	return g
}

func (g *Google) Enabled(ctx context.Context, level slog.Level) bool {
	return g.Handler.Enabled(ctx, level)
}

func (g *Google) Handle(ctx context.Context, record slog.Record) error {
	var data = make(map[string]any)
	record.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.String()
		return true
	})
	data["message"] = record.Message
	l := Line{
		Level:     record.Level,
		Timestamp: record.Time.Unix(),
		Time:      record.Time.Format(time.RFC3339),
		Data:      data,
	}
	g.Print(l)

	if g.ForwardHandler {
		return g.Handler.Handle(ctx, record)
	}
	return nil
}

func (g *Google) WithAttrs(attrs []slog.Attr) slog.Handler {
	return g.Handler.WithAttrs(attrs)
}

func (g *Google) WithGroup(name string) slog.Handler {
	return g.Handler.WithGroup(name)
}
