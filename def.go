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
	"fmt"
	"reflect"
	"strings"
)

type queryFunc func(string, ...any) (*sql.Rows, error)
type updateFunc func(string, ...any) (sql.Result, error)

func getSchemaAttributes[T any](obj T, schema string) (fields []int, _ error) {
	to := reflect.TypeOf(obj)
	voPtr := reflect.ValueOf(&obj)
	elem := voPtr.Elem()

	if elem.Kind() != reflect.Struct {
		return nil, errors.New("obj is not a struct")
	}

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)

		if !field.IsExported() {
			return nil, fmt.Errorf("field %d (%s) is not exported", i, field.Name)
		}

		// not a schema attribute
		if _, b := field.Tag.Lookup(schema); !b {
			continue
		}

		voField := elem.Field(i)
		if !voField.CanSet() {
			return nil, fmt.Errorf("field %d (%s) can't set", i, field.Name)
		}

		fields = append(fields, i)
	}

	return
}

func getSchemaFieldsWithSpecificAttribute[T any](obj T, schema string, specific string) (fields []int, _ error) {
	to := reflect.TypeOf(obj)
	voPtr := reflect.ValueOf(&obj)
	elem := voPtr.Elem()

	if elem.Kind() != reflect.Struct {
		return nil, errors.New("obj is not a struct")
	}

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)

		// not a schema attribute
		if _, b := field.Tag.Lookup(schema); !b {
			continue
		}

		if !field.IsExported() {
			return nil, fmt.Errorf("field %d (%s) is not exported", i, field.Name)
		}

		if str, b := field.Tag.Lookup(specific); b && strings.ToLower(str) == "true" {
			fields = append(fields, i)
		}
	}

	return
}

func getSchemaPK[T any](obj T, schema string) ([]int, error) {
	return getSchemaFieldsWithSpecificAttribute(obj, schema, "pk")
}

func getSchemaLike[T any](obj T, schema string) ([]int, error) {
	return getSchemaFieldsWithSpecificAttribute(obj, schema, "like")
}

func getSchemaValues[T any](obj T, fields []int) (values []interface{}) {
	var vo = reflect.ValueOf(obj)

	for _, fx := range fields {
		values = append(values, vo.Field(fx).Interface())
	}

	return
}

func exec(query updateFunc, queryStr string, values ...interface{}) (affectedRows int64, err error) {
	res, err := query(queryStr, values...)

	if err == nil {
		affectedRows, err = res.RowsAffected()
	}

	return
}

func inquery[T any](obj T, query queryFunc, fields []int, queryStr string, values ...interface{}) (A []T, err error) {
	rows, err := query(queryStr, values...)

	if err == nil {
		defer rows.Close()

		for err == nil && rows.Next() {
			err = decodeRow(rows, &obj, fields)

			if err == nil {
				A = append(A, obj)
			}
		}
	}

	return
}
