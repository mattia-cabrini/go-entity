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

func getUpdateQueryByPKEM[T any](fields []int, pkFields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("update `")
	builder.WriteString(table)
	builder.WriteString("` set ")

	for i, fx := range fields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteString("` = ?")

		if i != len(fields)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(" where ")

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

func getUpdateQueryByLike[T any](fields []int, likeFields []int, obj T, table string, schema string) (query string) {
	var to = reflect.TypeOf(obj)
	var builder = strings.Builder{}

	builder.WriteString("update `")
	builder.WriteString(table)
	builder.WriteString("` set ")

	for i, fx := range fields {
		builder.WriteRune('`')
		builder.WriteString(to.Field(fx).Tag.Get(schema))
		builder.WriteString("` = ?")

		if i != len(fields)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(" where ")

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

func UpdateByPKEM[T any](objOriginal T, objValue T, query updateFunc, table string, schema string) (affectedRows int64, err error) {
	var fields, pkFields []int
	var queryStr string

	fields, err = getSchemaAttributes(objValue, schema)

	if err == nil {
		pkFields, err = getSchemaPK(objValue, schema)

		if err == nil {
			queryStr = getUpdateQueryByPKEM(fields, pkFields, objValue, table, schema)

			values := getSchemaValues(objValue, fields)
			values = append(values, getSchemaValues(objOriginal, pkFields)...)

			affectedRows, err = exec(query, queryStr, values...)
		}
	}

	return
}

func UpdateByLike[T any](objOriginal T, objValue T, query updateFunc, table string, schema string) (affectedRows int64, err error) {
	var fields, likeFields []int
	var queryStr string

	fields, err = getSchemaAttributes(objValue, schema)

	if err == nil {
		likeFields, err = getSchemaLike(objValue, schema)

		if err == nil {
			queryStr = getUpdateQueryByLike(fields, likeFields, objValue, table, schema)

			values := getSchemaValues(objValue, fields)
			values = append(values, getSchemaValues(objOriginal, likeFields)...)

			affectedRows, err = exec(query, queryStr, values...)
		}
	}

	return
}
