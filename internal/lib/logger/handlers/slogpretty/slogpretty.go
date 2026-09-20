package slogpretty

import (
	// "context" 
	"io"
	stdLog "log"
	"log/slog"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	h := PrettyHandler{
		opts: opts,
		l: stdLog.New(out, "", 0),
		attrs: make([]slog.Attr, 0, 4),
	}
	return &h;
}

// func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
// 	level := r.Level.String() + ": ";
// 	switch r.Level {
// 	case slog.LevelDebug:
// 		level = "DEBUG";
// 	case slog.LevelInfo:
// 		level = "INFO";
// 	case slog.LevelWarn:
// 		level = "WARN";
// 	case slog.LevelError:
// 		level = "ERROR";
// 	}
// 	fields := make(map[string]interface{},r,r.NumAttrs());

// 	return nil;
// }

