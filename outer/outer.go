package outer

import (
	"context"

	"github.com/pepabo/cwlq/parser"
)

type Outer interface {
	Write(ctx context.Context, in <-chan *parser.Parsed)
	Close() error
	Err() error
}
