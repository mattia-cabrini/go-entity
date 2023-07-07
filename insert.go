/**
 * Copyright 2023 Mattia Cabrini
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package entity

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

func getInsertQuery[T any](fields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("insert into `")
	builder.WriteString(table)
	builder.WriteString("` (")

	for i, fx := range fields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteRune('`')

		if i != len(fields)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(") values (")

	for i := range fields {
		builder.WriteRune('?')

		if i != len(fields)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(");")

	return builder.String()
}

func Insert[T any](obj T, query updateFunc, table string, schema string) (lastInsertedId int64, err error) {
	var affectedRows int64
	var res sql.Result
	fields, err := getSchemaAttributes(obj, schema)

	if err == nil {
		queryStr := getInsertQuery(fields, obj, table, schema)

		res, err = query(queryStr, getSchemaValues(obj, fields)...)

		if err == nil {
			affectedRows, err = res.RowsAffected()

			if err == nil {
				if affectedRows == 0 {
					err = errors.New("no row affected")
				} else {
					lastInsertedId, err = res.LastInsertId()
				}
			}
		}
	}

	return
}
