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
	"testing"
)

func TestSelect(t *testing.T) {
	obj := testType{}

	fields, _ := getSchemaAttributes(obj, "schema")

	exp := "select `a`, `b`, `c` from tab;"
	query := getSelectAllQuery(fields, obj, "tab", "schema")

	if query != exp {
		t.Errorf("\ngot ```%s```\nexp ```%s```", query, exp)
	}

	db := openTestConnection()

	a, err := SelectAll(obj, db.Query, "tab", "schema")

	if err != nil {
		t.Error(err)
		return
	}

	if len(a) != 3 {
		t.Errorf("expected: 3, got %d", len(a))
	}
}
