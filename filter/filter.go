package filter

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/parser/lexer"
	"github.com/pepabo/cwlq/parser"
)

const (
	timestampKey = "timestamp"
	messageKey   = "message"
	rawKey       = "raw"
)

type Filter struct {
	conds    []string
	fns      map[string]interface{}
	total    int64
	filtered int64
	err      error
}

func (f *Filter) Filter(ctx context.Context, in <-chan *parser.Parsed) <-chan *parser.Parsed {
	out := make(chan *parser.Parsed)
	go func() {
		defer close(out)
	L:
		for i := range in {
			tf, err := f.evalConds(i)
			if err != nil {
				f.err = err
				break L
			}
			f.total += 1
			if tf {
				f.filtered += 1
				out <- i
			}
			select {
			case <-ctx.Done():
				break L
			default:
			}
		}
	}()
	return out
}

func New(conds []string) *Filter {
	trimed := []string{}
	for _, c := range conds {
		trimed = append(trimed, trimComment(c))
	}
	return &Filter{
		conds: trimed,
		fns:   map[string]interface{}{},
	}
}

func (f *Filter) Total() int64 {
	return f.total
}

func (f *Filter) Filtered() int64 {
	return f.filtered
}

func (f *Filter) Err() error {
	return f.err
}

func (f *Filter) AddFunc(key string, fn interface{}) error {
	if key == timestampKey || key == messageKey || key == rawKey {
		return fmt.Errorf("'%s' is reserved", key)
	}
	if _, ok := f.fns[key]; ok {
		return fmt.Errorf("'%s' is already exists", key)
	}
	f.fns[key] = fn
	return nil
}

func (f *Filter) evalConds(p *parser.Parsed) (bool, error) {
	if len(f.conds) == 0 {
		return true, nil
	}
	for _, c := range f.conds {
		env := map[string]interface{}{
			timestampKey: p.Timestamp,
			messageKey:   p.Message,
			rawKey:       p.LogEvent.Raw,
		}
		for k, fn := range f.fns {
			env[k] = fn
		}
		tf, err := evalCond(c, env)
		if err != nil {
			return false, err
		}
		if tf {
			return true, nil
		}
	}
	return false, nil
}

func evalCond(cond string, env interface{}) (bool, error) {
	v, err := expr.Eval(cond, env)
	if err != nil {
		return false, err
	}
	switch vv := v.(type) {
	case bool:
		return vv, nil
	default:
		return false, nil
	}
}

func trimComment(cond string) string {
	const commentToken = "#"
	trimed := []string{}
	for _, l := range strings.Split(cond, "\n") {
		if strings.HasPrefix(strings.Trim(l, " "), commentToken) {
			continue
		}
		s := file.NewSource(l)
		tokens, err := lexer.Lex(s)
		if err != nil {
			trimed = append(trimed, l)
			continue
		}

		ccol := -1
		inClosure := false
	L:
		for _, t := range tokens {
			switch {
			case t.Kind == lexer.Bracket && t.Value == "{":
				inClosure = true
			case t.Kind == lexer.Bracket && t.Value == "}":
				inClosure = false
			case t.Kind == lexer.Operator && t.Value == commentToken && !inClosure:
				ccol = t.Column
				break L
			}
		}
		if ccol > 0 {
			trimed = append(trimed, strings.TrimSuffix(l[:ccol], " "))
			continue
		}

		trimed = append(trimed, l)
	}
	return strings.Join(trimed, "\n")
}
