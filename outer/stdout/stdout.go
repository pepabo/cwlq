package stdout

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pepabo/cwlf/outer"
	"github.com/pepabo/cwlf/parser"
)

var _ outer.Outer = (*Stdout)(nil)

type Stdout struct {
	out io.Writer
	err error
}

func New() *Stdout {
	return &Stdout{out: os.Stdout}
}

func (o *Stdout) Write(ctx context.Context, in <-chan *parser.Parsed) {
	for e := range in {
		if _, err := fmt.Fprintf(o.out, "%s\n", e.LogEvent.Raw); err != nil {
			o.err = err
			break
		}
	}
}

func (o *Stdout) Close() error {
	return nil
}

func (o *Stdout) Err() error {
	return o.err
}
