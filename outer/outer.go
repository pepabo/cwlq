package outer

import (
	"context"

	"github.com/pepabo/cwlf/parser"
)

type Outer interface {
	Write(ctx context.Context, in <-chan *parser.Parsed)
	Close() error
	Err() error
}
