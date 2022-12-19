package stdout

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pepabo/cwlq/outer"
	"github.com/pepabo/cwlq/parser"
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
	for p := range in {
		b, err := json.Marshal(p)
		if err != nil {
			o.err = err
			break
		}
		if _, err := fmt.Fprintf(o.out, "%s\n", string(b)); err != nil {
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
