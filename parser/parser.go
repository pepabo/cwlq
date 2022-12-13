package parser

import (
	"context"

	"github.com/pepabo/cwlf/datasource"
)

type Parsed struct {
	Data     map[string]interface{}
	LogEvent *datasource.LogEvent
}

type Parser interface {
	Parse(context.Context, <-chan *datasource.LogEvent) <-chan *Parsed
	NewFakeLogEvent() (*datasource.LogEvent, error)
}
