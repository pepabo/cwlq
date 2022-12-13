package s3

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pepabo/cwlf/datasource"
)

var _ datasource.Datasource = (*S3)(nil)

type S3 struct {
	client *s3.Client
	bucket string
	prefix string
	err    error
}

type Log struct {
	LogEvents []*datasource.LogEvent `json:"logEvents"`
}

func (s *S3) Fetch(ctx context.Context) <-chan *datasource.LogEvent {
	out := make(chan *datasource.LogEvent)
	go func() {
		defer close(out)
		var t *string
	L:
		for {
			o, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:            aws.String(s.bucket),
				Prefix:            aws.String(s.prefix),
				ContinuationToken: t,
			})
			if err != nil {
				s.err = err
				break L
			}
			for _, c := range o.Contents {
				if err := func() error {
					obj, err := s.client.GetObject(ctx, &s3.GetObjectInput{
						Bucket: aws.String(s.bucket),
						Key:    c.Key,
					})
					if err != nil {
						return err
					}
					defer obj.Body.Close()
					gr, err := gzip.NewReader(obj.Body)
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
					}
					return nil
				}(); err != nil {
					s.err = err
					break L
				}
			}
			if o.NextContinuationToken == nil {
				break L
			}
			t = o.NextContinuationToken
		}
	}()
	return out
}

func (s *S3) Err() error {
	return s.err
}

func New(cfg aws.Config, dsn string) (*S3, error) {
	if !strings.HasPrefix(dsn, "s3://") {
		return nil, fmt.Errorf("invalid s3 bucket url: %s", dsn)
	}
	splitted := strings.SplitN(strings.TrimPrefix(dsn, "s3://"), "/", 2)
	if len(splitted) == 0 || splitted[0] == "" {
		return nil, fmt.Errorf("invalid s3 bucket url: %s", dsn)
	}
	bucket := splitted[0]
	prefix := splitted[1]
	return &S3{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		prefix: prefix,
	}, nil
}
