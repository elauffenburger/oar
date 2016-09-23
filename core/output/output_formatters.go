package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"io"

	res "github.com/elauffenburger/oar/core/results"
)

type OutputFormatter interface {
	Format(results *res.Results) string
	FormatToStream(results *res.Results, stream io.Writer)
}

type JsonOutputFormatter struct {
}

func (formatter *JsonOutputFormatter) Format(results *res.Results) string {
	jsonArray := formatter.ToJsonArray(results)
	marshalledbytes, err := json.Marshal(&jsonArray)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling results: '%s", err))
	}

	return string(marshalledbytes)
}

func (formatter *JsonOutputFormatter) FormatToStream(results *res.Results, stream io.Writer) {
	bytes := []byte(formatter.Format(results))

	stream.Write(bytes)
}

type JsonObject map[string]string
type JsonArray []*JsonObject

func (obj JsonObject) Keys() []string {
	keys := make([]string, len(obj))

	i := 0
	for key, _ := range obj {
		keys[i] = key

		i++
	}

	return keys
}

func (formatter *JsonOutputFormatter) ToJsonArray(results *res.Results) JsonArray {
	result := make(JsonArray, results.NumRows())

	for i, set := range results.Rows {
		object := make(JsonObject)

		for _, entry := range set.Values {
			object[entry.Name] = entry.Value
		}

		result[i] = &object
	}

	return result
}

const maxLinesPerSqlStmt = 1000

type SqlOutputFormatter struct {
	TableName string
}

func (formatter *SqlOutputFormatter) Format(results *res.Results) string {
	out := ""

	formatter.formatInternal(results, func(str *string) {
		out += *str
	})

	return out
}

func (formatter *SqlOutputFormatter) formatInternal(results *res.Results, writeFn func(str *string)) {
	if len(results.Rows) == 0 {
		return
	}

	// generate initial "insert into dbo.foobar(...) values" stmt
	insertHeaderStr := fmt.Sprintf("insert into %s (", formatter.TableName)
	{
		row := results.Rows[0]
		for i, val := range row.Values {
			insertHeaderStr += fmt.Sprintf("[%s]", val.Name)

			if i != len(row.Values)-1 {
				insertHeaderStr += ","
			}
		}
	}

	insertHeaderStr += ") values \n"
	writeFn(&insertHeaderStr)

	// generate subsequent "(...)" stmts and start new "insert..." stmts as necessary
	n := len(results.Rows)
	for i, row := range results.Rows {

		// Generate (...) stmt for this row
		rowstr := "("
		for i, val := range row.Values {
			// escape 's in values
			value := strings.Replace(val.Value, "'", "''", -1)

			rowstr += fmt.Sprintf("'%s'", value)

			if i != len(row.Values)-1 {
				rowstr += ","
			}
		}
		rowstr += ")"

		// if we're writing the max allowed insert statement (1000), end the current stmt
		if (i+1)%maxLinesPerSqlStmt == 0 {
			rowstr += ";\n"

			// if we're not writing the last line, start a new insert stmt
			if i != n-1 {
				rowstr += insertHeaderStr
			}
		} else {
			// otherwise...

			// if we're not writing the last line, add a comma separator
			if i != n-1 {
				rowstr += ",\n"
			}
		}

		writeFn(&rowstr)
	}
}

func (formatter *SqlOutputFormatter) FormatToStream(results *res.Results, stream io.Writer) {
	formatter.formatInternal(results, func(str *string) {
		bytes := []byte(*str)

		stream.Write(bytes)
	})
}

func NewSqlOutputFormatter(tablename string) *SqlOutputFormatter {
	return &SqlOutputFormatter{TableName: tablename}
}
