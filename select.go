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
	"fmt"
	"reflect"
	"strings"

	"github.com/mattia-cabrini/go-utility"
)

func getSelectAllQuery[T any](fields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("select ")

	for i, fx := range fields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteRune('`')

		if i != len(fields)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(" from ")
	builder.WriteString(table)
	builder.WriteRune(';')

	return builder.String()
}

func decodeRow[T any](rows *sql.Rows, obj *T, fields []int) (err error) {
	var vo = reflect.ValueOf(obj).Elem()
	var to = reflect.TypeOf(*obj)
	var dest = make([]any, len(fields))

	for i, field := range fields {
		if v, b := utility.TypeDefault(vo.Field(field).Kind()); b {
			dest[i] = &v
		} else {
			return fmt.Errorf("field %d (%s) has unusable kind: %v", field, to.Field(field).Name, vo.Field(field).Kind())
		}
	}

	err = rows.Scan(dest...)

	if err == nil {
		for i, field := range fields {
			fieldVO := vo.Field(field)
			value := utility.IPtrToI(fieldVO, dest[i])
			fieldVO.Set(reflect.ValueOf(value))
		}
	}

	return
}

func SelectAll[T any](obj T, query queryFunc, table string, schema string) (A []T, err error) {
	fields, err := getSchemaAttributes(obj, schema)

	if err == nil {
		queryStr := getSelectAllQuery(fields, obj, table, schema)
		A, err = inquery(obj, query, fields, queryStr)
	}

	return
}
