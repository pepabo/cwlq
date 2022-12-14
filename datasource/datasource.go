package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-faker/faker/v4"
)

const (
	S3    = "s3"
	Local = "local"
	Fake  = "fake"
)

type LogEvent struct {
	ID        string
	Timestamp time.Time
	Message   string
	Raw       string
}

type Datasource interface {
	Fetch(context.Context) <-chan *LogEvent
	Err() error
}

func NewFakeID() string {
	return fmt.Sprintf("%d%d", faker.RandomUnixTime(), faker.RandomUnixTime())
}

func (e LogEvent) MarshalJSON() ([]byte, error) {
	s := struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Message   string `json:"message"`
	}{
		ID:        e.ID,
		Timestamp: e.Timestamp.UnixMilli(),
		Message:   e.Message,
	}
	return json.Marshal(s)
}

func (e *LogEvent) UnmarshalJSON(b []byte) error {
	s := struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Message   string `json:"message"`
	}{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	e.ID = s.ID
	e.Timestamp = time.UnixMilli(s.Timestamp)
	e.Message = s.Message
	e.Raw = string(b)
	return nil
}
