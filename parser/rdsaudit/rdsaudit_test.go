package rdsaudit

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pepabo/cwlq/datasource"
)

func TestParse(t *testing.T) {
	r := New()
	in := make(chan *datasource.LogEvent, 1)
	le, err := r.NewFakeLogEvent()
	if err != nil {
		t.Fatal(err)
	}
	in <- le
	out := r.Parse(context.Background(), in)
	p := <-out
	if p == nil {
		t.Fatal(r.Err())
	}
	got := p.LogEvent
	if diff := cmp.Diff(got, le, nil); diff != "" {
		t.Errorf("%s", diff)
	}
}
