# cwlq

cwlq is a tool/package for querying logs (of Amazon CloudWatch Logs) stored in various datasources.

## Usage

``` console
$ cwlq s3://myrds-audit-logs/2022/12/11/ --parser rdsaudit --filter "message.host == '10.0.1.123'" --filter "message.object contains 'INSERT'"
```

## Support datasource

### Amazon S3

`s3://bucket/path/to`

### Local file or directory

`local://path/to` `local:///root/path/to`

### Fake datasource

`fake://rdsaudit?duration=3sec`

### Amazon CloudWatch Logs directly

WIP

> **Note**
> Perhaps it would be better to use [CloudWatch Logs Insights](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AnalyzingLogData.html).

## Support Parser

### `rdsaudit`

Parser for gziped logs via MariaDB Audit Plugin for Amazon RDS.

| Field | Example | Description |
| --- | --- | --- |
| `timestamp` | `1670717181000` | The Unix time stamp for the logged event with microsecond precision. |
| `message.timestamp`  | `20221211 00:06:21` | [The Unix time stamp for the logged event with microsecond precision????](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.serverhost` | `ip-10-0-0-123` | [The name of the instance that the event is logged for.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.username` | `redash` | [The connected user name of the user.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.host` | `10.0.1.123` | [The host that the user connected from.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.connectionid` | `502547196` | [The connection ID number for the logged operation.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.queryid` | `84996781288` | [The query ID number, which can be used for finding the relational table events and related queries. For `TABLE` events, multiple lines are added.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.operation` | `QUERY` | [The recorded action type. Possible values are: `CONNECT`, `QUERY`, `READ`, `WRITE`, `CREATE`, `ALTER`, `RENAME`, and `DROP`.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.database` | `dbname` | [The active database, as set by the `USE` command.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.object` | `SELECT * FROM accounts;` | [For `QUERY` events, this value indicates the query that the database performed. For `TABLE` events, it indicates the table name.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.retcode` | `0` | [The return code of the logged operation.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `message.connection_type` | `1` | [The security state of the connection to the server.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.MySQL.Options.AuditPlugin.html) |
| `raw` | `` | Raw data of log event. |

## Install

**homebrew tap:**

```console
$ brew install pepabo/tap/cwlq
```

**manually:**

Download binany from [releases page](https://github.com/pepabo/cwlq/releases)

**go install:**

```console
$ go install github.com/pepabo/cwlq@latest
```
