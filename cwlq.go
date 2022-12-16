package cwlq

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pepabo/cwlq/datasource"
	"github.com/pepabo/cwlq/datasource/fake"
	"github.com/pepabo/cwlq/datasource/local"
	"github.com/pepabo/cwlq/datasource/s3"
	"github.com/pepabo/cwlq/filter"
	"github.com/pepabo/cwlq/outer"
	"github.com/pepabo/cwlq/outer/stdout"
	"github.com/pepabo/cwlq/parser"
	"github.com/pepabo/cwlq/parser/rdsaudit"
)

const defaultRegion = "ap-northeast-1"

type Cwlf struct {
	d datasource.Datasource
	p parser.Parser
	f *filter.Filter
	o outer.Outer
}

func New(dsn, parserType string, filters []string) (*Cwlf, error) {
	// datasource
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	var d datasource.Datasource
	switch u.Scheme {
	case datasource.S3:
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, err
		}
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}
		d, err = s3.New(cfg, dsn)
		if err != nil {
			return nil, err
		}
	case datasource.Local:
		d, err = local.New(dsn)
		if err != nil {
			return nil, err
		}
	case datasource.Fake:
		d, err = fake.New(dsn)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsuppoted scheme: %s", dsn)
	}

	// parser
	var p parser.Parser
	switch parserType {
	case parser.RDSAudit:
		p = rdsaudit.New()
	default:
		return nil, fmt.Errorf("unsuppoted parser: %s", parserType)
	}

	// filter
	f := filter.New(filters)

	return &Cwlf{
		d: d,
		p: p,
		f: f,
		o: stdout.New(),
	}, nil
}

func (c *Cwlf) Outer(o outer.Outer) {
	c.o = o
}

func (c *Cwlf) Total() int64 {
	return c.f.Total()
}

func (c *Cwlf) Filtered() int64 {
	return c.f.Filtered()
}

func (c *Cwlf) Run(ctx context.Context) (err error) {
	c.o.Write(ctx, c.f.Filter(ctx, c.p.Parse(ctx, c.d.Fetch(ctx))))

	defer func() {
		if cerr := c.o.Close(); cerr != nil {
			err = cerr
		}
	}()

	if err := c.d.Err(); err != nil {
		return err
	}

	if err := c.p.Err(); err != nil {
		return err
	}

	if err := c.f.Err(); err != nil {
		return err
	}

	if err := c.o.Err(); err != nil {
		return err
	}

	return nil
}
