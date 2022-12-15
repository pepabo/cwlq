package datasource

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJSON(t *testing.T) {
	want := &LogEvent{
		ID:        NewFakeID(),
		Timestamp: time.Unix(faker.UnixTime(), 0),
		Message:   faker.Email(),
	}
	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	got := &LogEvent{}
	if err := json.Unmarshal(b, got); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(LogEvent{}, "Raw")); diff != "" {
		t.Errorf("%s", diff)
	}
}
