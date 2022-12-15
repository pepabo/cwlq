package filter

import (
	"fmt"
	"testing"

	"github.com/pepabo/cwlq/datasource"
	"github.com/pepabo/cwlq/parser"
)

func TestEvalConds(t *testing.T) {
	tests := []struct {
		conds []string
		fns   map[string]interface{}
		in    *parser.Parsed
		want  bool
	}{
		{
			[]string{},
			map[string]interface{}{},
			newParsed(t, map[string]interface{}{}),
			true,
		},
		{
			[]string{"message.scope == 'world'"},
			map[string]interface{}{},
			newParsed(t, map[string]interface{}{
				"scope": "world",
			}),
			true,
		},
		{
			[]string{"message.scope != 'world'", "message.scope == 'world'"},
			map[string]interface{}{},
			newParsed(t, map[string]interface{}{
				"scope": "world",
			}),
			true,
		},
		{
			[]string{"truefunc()"},
			map[string]interface{}{
				"truefunc": func() bool { return true },
			},
			newParsed(t, map[string]interface{}{}),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.conds), func(t *testing.T) {
			f := New(tt.conds)
			for k, fn := range tt.fns {
				f.AddFunc(k, fn)
			}
			got, err := f.evalConds(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func newParsed(t *testing.T, msg map[string]interface{}) *parser.Parsed {
	t.Helper()
	return &parser.Parsed{
		Message: msg,
		LogEvent: &datasource.LogEvent{
			ID:      datasource.NewFakeID(),
			Message: "fake message",
			Raw:     "fake raw",
		},
	}
}
