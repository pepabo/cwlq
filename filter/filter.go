package filter

import (
	"context"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/parser/lexer"
	"github.com/pepabo/cwlf/parser"
)

type Filter struct {
	conds []string
	err   error
}

func (f *Filter) Filter(ctx context.Context, in <-chan *parser.Parsed) <-chan *parser.Parsed {
	out := make(chan *parser.Parsed)
	go func() {
		defer close(out)
		if len(f.conds) == 0 {
			for i := range in {
				out <- i
			}
		} else {
		L:
			for i := range in {
				for _, c := range f.conds {
					env := map[string]interface{}{
						"timestamp": i.LogEvent.Timestamp.UnixMilli(),
						"message":   i.Data,
						"raw":       i.LogEvent.Raw,
					}
					tf, err := evalCond(c, env)
					if err != nil {
						f.err = err
						break L
					}
					if tf {
						out <- i
					}
				}
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
	}
}

func (f *Filter) Err() error {
	return f.err
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
