package filter

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/pepabo/cwlq/datasource"
	"github.com/pepabo/cwlq/parser"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		in           []*parser.Parsed
		conds        []string
		wantTotal    int64
		wantFiltered int64
	}{
		{
			[]*parser.Parsed{
				newParsed(t, map[string]interface{}{
					"scope": "world",
				}),
				newParsed(t, map[string]interface{}{
					"scope": "world",
				}),
			},
			[]string{},
			2,
			2,
		},
		{
			[]*parser.Parsed{
				newParsed(t, map[string]interface{}{
					"scope": "world",
				}),
				newParsed(t, map[string]interface{}{
					"scope": "space",
				}),
				newParsed(t, map[string]interface{}{
					"scope": "local",
				}),
			},
			[]string{"message.scope == 'world'"},
			3,
			1,
		},
	}
	ctx := context.Background()
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f := New(tt.conds)
			in := make(chan *parser.Parsed)
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for range f.Filter(ctx, in) {
				}
				wg.Done()
			}()
			for _, p := range tt.in {
				in <- p
			}
			close(in)
			wg.Wait()
			{
				got := f.Total()
				if got != tt.wantTotal {
					t.Errorf("got %v\nwant %v", got, tt.wantTotal)
				}
			}
			{
				got := f.Filtered()
				if got != tt.wantFiltered {
					t.Errorf("got %v\nwant %v", got, tt.wantFiltered)
				}
			}
		})
	}
}

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
				if err := f.AddFunc(k, fn); err != nil {
					t.Fatal(err)
				}
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
