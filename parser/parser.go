package parser

import (
	"context"

	"github.com/pepabo/cwlq/datasource"
)

const RDSAudit = "rdsaudit"

type Parsed struct {
	Message  map[string]interface{}
	LogEvent *datasource.LogEvent
}

type Parser interface {
	Parse(context.Context, <-chan *datasource.LogEvent) <-chan *Parsed
	NewFakeLogEvent() (*datasource.LogEvent, error)
	Err() error
}
