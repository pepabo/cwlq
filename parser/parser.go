package parser

import (
	"context"

	"github.com/pepabo/cwlq/datasource"
)

const RDSAudit = "rdsaudit"

type Parsed struct {
	Timestamp int64                  `json:"timestamp"`
	Message   map[string]interface{} `json:"message"`
	LogEvent  *datasource.LogEvent   `json:"-"`
}

type Parser interface {
	Parse(context.Context, <-chan *datasource.LogEvent) <-chan *Parsed
	ParseLogEvent(*datasource.LogEvent) ([]*Parsed, error)
	NewFakeLogEvent() (*datasource.LogEvent, error)
	Err() error
}
