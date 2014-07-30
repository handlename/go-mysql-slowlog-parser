package slowlog

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

// Parser ...
type Parser struct {
}

// NewParser ... 新しい Parser をつくる
func NewParser() Parser {
	return Parser{}
}

// Parse ... ファイルをパースして結果を流す
func (p Parser) Parse(r io.Reader) <-chan Parsed {
	ch := make(chan Parsed)

	reg1, _ := regexp.Compile(`^#? User\@Host:\s+(\S+)\s+\@\s+(\S+).*`)
	reg2, _ := regexp.Compile(`^# Query_time: ([0-9.]+)\s+Lock_time: ([0-9.]+)\s+Rows_sent: ([0-9.]+)\s+Rows_examined: ([0-9.]+).*`)

	go func() {
		scanner := bufio.NewScanner(r)
		parsed := Parsed{}

		for scanner.Scan() {
			line := scanner.Text()

			if shouldIgnore(line) {
				continue
			}

			if strings.HasPrefix(line, "# Time:") {
				parsed = Parsed{}
			}

			// DateTime
			if strings.HasPrefix(line, "# Time:") {
				t, err := time.Parse(
					"060102 15:04:05.999999",
					strings.Replace(line, "# Time: ", "", 1))

				if err != nil {
					log.Println(err)
				}

				parsed.Datetime = t.Unix()

				continue
			}

			// User, Host
			if r := reg1.FindStringSubmatch(line); r != nil {
				parsed.User = r[1]
				parsed.Host = r[2]
				continue
			}

			// QueryTime, LockTime, RowsSent, RowsExamined
			if r := reg2.FindStringSubmatch(line); r != nil {
				parsed.QueryTime = stringToFloat32(r[1])
				parsed.LockTime = stringToFloat32(r[2])
				parsed.RowsSent = stringToInt32(r[3])
				parsed.RowsExamined = stringToInt32(r[4])
				continue
			}

			// Sql
			if !strings.HasPrefix(line, "#") {
				parsed.Sql += strings.Trim(line, " \r\n") + " "

				if strings.HasSuffix(line, ";") && parsed.Sql != "" {
					parsed.Sql = strings.Trim(parsed.Sql, " ")
					ch <- parsed
					parsed.Sql = ""
				}

				continue
			}
		}

		ch <- parsed

		close(ch)
	}()

	return ch
}

func shouldIgnore(line string) bool {
	uppered := strings.ToUpper(line)

	return strings.HasPrefix(uppered, "USE") ||
		strings.HasPrefix(uppered, "SET TIMESTAMP=")
}