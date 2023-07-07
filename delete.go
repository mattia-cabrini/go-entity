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
	"reflect"
	"strings"
)

func getDeleteQueryByEM[T any](pkFields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("delete from `")
	builder.WriteString(table)
	builder.WriteString("` where ")

	for i, fx := range pkFields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteString("` = ?")

		if i != len(pkFields)-1 {
			builder.WriteString(" and ")
		}
	}

	builder.WriteString(";")

	return builder.String()
}

func getDeleteQueryByLike[T any](likeFields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("delete from `")
	builder.WriteString(table)
	builder.WriteString("` where ")

	for i, fx := range likeFields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteString("` like concat('%', ?, '%')")

		if i != len(likeFields)-1 {
			builder.WriteString(" and ")
		}
	}

	builder.WriteString(";")

	return builder.String()
}

func DeleteByPKEM[T any](obj T, query updateFunc, table string, schema string) (affectedRows int64, err error) {
	var pkFields []int
	var queryStr string

	pkFields, err = getSchemaPK(obj, "schema")

	if err == nil {
		queryStr = getDeleteQueryByEM(pkFields, obj, table, schema)
		affectedRows, err = exec(query, queryStr, getSchemaValues(obj, pkFields)...)
	}

	return
}

func DeleteByLike[T any](obj T, query updateFunc, table string, schema string) (affectedRows int64, err error) {
	var likeFields []int
	var queryStr string

	likeFields, err = getSchemaLike(obj, "schema")

	if err == nil {
		queryStr = getDeleteQueryByLike(likeFields, obj, table, schema)
		affectedRows, err = exec(query, queryStr, getSchemaValues(obj, likeFields)...)
	}

	return
}
