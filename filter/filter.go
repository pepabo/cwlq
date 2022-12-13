package filter

import (
	"context"

	"github.com/pepabo/cwlf/parser"
)

type Filter struct {
	conds []string
}

func (f *Filter) Filter(ctx context.Context, in <-chan *parser.Parsed) <-chan *parser.Parsed {
	out := make(chan *parser.Parsed)
	go func() {
		defer close(out)
		for i := range in {
			// TODO: use f.conds
			out <- i
		}
	}()
	return out
}

func New(conds []string) *Filter {
	return &Filter{
		conds: conds,
	}
}
