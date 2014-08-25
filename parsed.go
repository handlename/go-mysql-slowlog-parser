package slowlog

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Parsed is struct store values of parsed slowlog.
type Parsed struct {
	Datetime     int64   `json:"datetime"`
	User         string  `json:"user"`
	Host         string  `json:"host"`
	QueryTime    float32 `json:"query_time"`
	LockTime     float32 `json:"lock_time"`
	RowsSent     int32   `json:"rows_sent"`
	RowsExamined int32   `json:"rows_examined"`
	Sql          string  `json:"sql"`
}

// AsLTSV returns parsed slowlog as LTSV format.
func (p *Parsed) AsLTSV() string {
	return strings.Join([]string{
		fmt.Sprintf("datetime:%d", p.Datetime),
		fmt.Sprintf("user:%s", p.User),
		fmt.Sprintf("host:%s", p.Host),
		fmt.Sprintf("query_time:%f", p.QueryTime),
		fmt.Sprintf("lock_time:%f", p.LockTime),
		fmt.Sprintf("rows_sent:%d", p.RowsSent),
		fmt.Sprintf("rows_examined:%d", p.RowsExamined),
		fmt.Sprintf("sql:%s", p.Sql),
	}, "\t")
}

// AsJSON returns parsed slowlog as JSON format.
func (p *Parsed) AsJSON() string {
	j, _ := json.Marshal(p)
	return string(j)
}
