package local

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pepabo/cwlf/datasource"
)

var _ datasource.Datasource = (*Local)(nil)

type Local struct {
	root string
	err  error
}

type Log struct {
	LogEvents []*datasource.LogEvent `json:"logEvents"`
}

func (d *Local) Fetch(ctx context.Context) <-chan *datasource.LogEvent {
	out := make(chan *datasource.LogEvent)
	go func() {
		defer close(out)

		if err := filepath.WalkDir(d.root, func(path string, fi fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if fi.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			gr, err := gzip.NewReader(f)
			if err != nil {
				return err
			}
			defer gr.Close()

			dec := json.NewDecoder(gr)
			for {
				l := Log{}
				if err := dec.Decode(&l); err == io.EOF {
					break
				} else if err != nil {
					return err
				}
				for _, le := range l.LogEvents {
					out <- le
				}

				select {
				case <-ctx.Done():
					break
				default:
				}
			}
			return nil
		}); err != nil {
			d.err = err
		}
	}()
	return out
}

func (d *Local) Err() error {
	return d.err
}

func New(dsn string) (*Local, error) {
	if !strings.HasPrefix(dsn, "local://") {
		return nil, fmt.Errorf("invalid local dsn: %s", dsn)
	}
	p := strings.TrimPrefix(dsn, "local://")
	if p == "" {
		return nil, fmt.Errorf("invalid local dsn: %s", dsn)
	}
	root, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}
	return &Local{
		root: root,
	}, nil
}
