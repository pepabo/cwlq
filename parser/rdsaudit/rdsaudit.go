package rdsaudit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/pepabo/cwlf/datasource"
	"github.com/pepabo/cwlf/parser"
)

var _ parser.Parser = (*RDSAudit)(nil)

// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html
type AuditLog struct {
	Timestamp      string `faker:"rds_timestamp"`
	Serverhost     string `faker:"ipv4"`
	Username       string `faker:"username"`
	Host           string `faker:"ipv4"`
	Connectionid   string `faker:"uuid_digit"`
	Queryid        string `faker:"uuid_digit"`
	Operation      string `faker:"oneof: CONNECT, QUERY, READ, WRITE, CREATE, ALTER, RENAME, DROP"`
	Database       string `faker:"domain_name"`
	Object         string `faker:"username"`
	Retcode        int64  `faker:"oneof: 0, 1"`
	ConnectionType int64  `faker:"oneof: 0, 1, 2, 3, 4, 5"`
}

type RDSAudit struct {
	err error
}

func (r *RDSAudit) Parse(ctx context.Context, in <-chan *datasource.LogEvent) <-chan *parser.Parsed {
	out := make(chan *parser.Parsed)
	go func() {
		defer close(out)
		for e := range in {
			a, err := parseMessage(e.Message)
			if err != nil {
				r.err = err
				break
			}
			out <- &parser.Parsed{
				Data:     a.ToMap(),
				LogEvent: e,
			}
		}
	}()
	return out
}

func New() *RDSAudit {
	return &RDSAudit{}
}

func (r *RDSAudit) NewFakeLogEvent() (*datasource.LogEvent, error) {
	al := AuditLog{}
	if err := faker.FakeData(&al); err != nil {
		return nil, err
	}
	if al.Operation == "QUERY" {
		al.Object = fmt.Sprintf("SELECT * FROM %s;", al.Object)
	}
	msg, err := al.ToCSV()
	if err != nil {
		return nil, err
	}
	le := &datasource.LogEvent{
		ID:        datasource.NewFakeID(),
		Timestamp: time.Unix(faker.RandomUnixTime(), 0),
		Message:   msg,
	}
	raw, err := json.Marshal(le)
	if err != nil {
		return nil, err
	}
	le.Raw = raw

	return le, nil
}

func (a *AuditLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":       a.Timestamp,
		"serverhost":      a.Serverhost,
		"username":        a.Username,
		"host":            a.Host,
		"connectionid":    a.Connectionid,
		"queryid":         a.Queryid,
		"operation":       a.Operation,
		"database":        a.Database,
		"object":          a.Object,
		"retcode":         a.Retcode,
		"connection_type": a.ConnectionType,
	}
}

func (a *AuditLog) ToCSV() (string, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	d := []string{
		a.Timestamp,
		a.Serverhost,
		a.Username,
		a.Host,
		a.Connectionid,
		a.Queryid,
		a.Operation,
		a.Database,
		a.Object,
		fmt.Sprintf("%d", a.Retcode),
		fmt.Sprintf("%d", a.ConnectionType),
	}
	if err := w.Write(d); err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func parseMessage(msg string) (*AuditLog, error) {
	r := csv.NewReader(strings.NewReader(msg))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 11 {
			return nil, fmt.Errorf("invalid message: %#v", msg)
		}
		retcode, err := strconv.ParseInt(record[9], 10, 64)
		if err != nil {
			return nil, err
		}
		connectionType, err := strconv.ParseInt(record[10], 10, 64)
		if err != nil {
			return nil, err
		}
		return &AuditLog{
			Timestamp:      record[0],
			Serverhost:     record[1],
			Username:       record[2],
			Host:           record[3],
			Connectionid:   record[4],
			Queryid:        record[5],
			Operation:      record[6],
			Database:       record[7],
			Object:         record[8],
			Retcode:        retcode,
			ConnectionType: connectionType,
		}, nil
	}
	return nil, fmt.Errorf("invalid message: %#v", msg)
}

func init() {
	_ = faker.AddProvider("rds_timestamp", func(v reflect.Value) (interface{}, error) {
		return time.Unix(faker.RandomUnixTime(), 0).Format("20060102 03:04:05"), nil
	})
}
