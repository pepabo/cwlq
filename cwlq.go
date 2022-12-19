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
	D datasource.Datasource
	P parser.Parser
	F *filter.Filter
	O outer.Outer
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
		D: d,
		P: p,
		F: f,
		O: stdout.New(),
	}, nil
}

func (c *Cwlf) Outer(o outer.Outer) {
	c.O = o
}

func (c *Cwlf) Total() int64 {
	return c.F.Total()
}

func (c *Cwlf) Filtered() int64 {
	return c.F.Filtered()
}

func (c *Cwlf) Run(ctx context.Context) (err error) {
	c.O.Write(ctx, c.F.Filter(ctx, c.P.Parse(ctx, c.D.Fetch(ctx))))

	defer func() {
		if cerr := c.O.Close(); cerr != nil {
			err = cerr
		}
	}()

	if err := c.D.Err(); err != nil {
		return err
	}

	if err := c.P.Err(); err != nil {
		return err
	}

	if err := c.F.Err(); err != nil {
		return err
	}

	if err := c.O.Err(); err != nil {
		return err
	}

	return nil
}
