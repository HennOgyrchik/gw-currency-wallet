package logs

import (
	"io"
	"log/slog"
)

type Log struct {
	logger *slog.Logger
}
type Attr struct {
	Key   string
	Value string
}

func New(w io.Writer) *Log {
	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)
	return &Log{logger: logger}
}

func (l *Log) Err(msg string, err error) {
	slog.Error(msg, "error", err.Error())
}

func (l *Log) Info(msg string, attr Attr) {

	slog.Info(msg, attr.Key, attr.Value)
}
