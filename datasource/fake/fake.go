package fake

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/k1LoW/duration"
	"github.com/pepabo/cwlq/datasource"
	"github.com/pepabo/cwlq/parser"
	"github.com/pepabo/cwlq/parser/rdsaudit"
)

const defaultDuration = "1min"

var _ datasource.Datasource = (*Fake)(nil)

type Fake struct {
	parser   parser.Parser
	duration time.Duration
	err      error
}

func (f *Fake) Fetch(ctx context.Context) <-chan *datasource.LogEvent {
	out := make(chan *datasource.LogEvent)
	timer := time.NewTimer(f.duration)

	go func() {
		defer close(out)
	L:
		for {
			le, err := f.parser.NewFakeLogEvent()
			if err != nil {
				f.err = err
				break
			}
			out <- le
			select {
			case <-ctx.Done():
				break L
			case <-timer.C:
				break L
			default:
			}
		}
	}()

	return out
}

func (f *Fake) Err() error {
	return f.err
}

func New(dsn string) (*Fake, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	if u.Scheme != datasource.Fake {
		return nil, fmt.Errorf("invalid fake url: %s", dsn)
	}
	ds := defaultDuration
	if u.Query().Has("duration") {
		ds = u.Query().Get("duration")
	}
	d, err := duration.Parse(ds)
	if err != nil {
		return nil, err
	}
	var p parser.Parser
	pt := u.Host
	switch pt {
	case parser.RDSAudit:
		p = rdsaudit.New()
	default:
		return nil, fmt.Errorf("unsupported parser: %s", pt)
	}

	return &Fake{
		parser:   p,
		duration: d,
	}, nil
}
