package cwlf

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pepabo/cwlf/datasource"
	"github.com/pepabo/cwlf/datasource/fake"
	"github.com/pepabo/cwlf/datasource/local"
	"github.com/pepabo/cwlf/datasource/s3"
	"github.com/pepabo/cwlf/filter"
	"github.com/pepabo/cwlf/outer"
	"github.com/pepabo/cwlf/outer/stdout"
	"github.com/pepabo/cwlf/parser"
	"github.com/pepabo/cwlf/parser/rdsaudit"
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

func (c *Cwlf) Run(ctx context.Context) (err error) {
	c.o.Write(ctx, c.f.Filter(ctx, c.p.Parse(ctx, c.d.Fetch(ctx))))

	defer func() {
		err = c.o.Close()
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

	return nil
}
