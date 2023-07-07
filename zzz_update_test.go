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

func TestUpdateEM(t *testing.T) {
	objValue := testType{A: "0", B: "1", C: ""}
	objOriginal := testType{A: "0", B: "1", C: "3"}

	fields, _ := getSchemaAttributes(objValue, "schema")
	pkFields, _ := getSchemaPK(objValue, "schema")

	exp := "update `tabUpdate` set `a` = ?, `b` = ?, `c` = ? where `a` = ? and `b` = ?;"
	query := getUpdateQueryByPKEM(fields, pkFields, objValue, "tabUpdate", "schema")

	if query != exp {
		t.Errorf("\ngot ```%s```\nexp ```%s```", query, exp)
	}

	db := openTestConnection()

	_, err := UpdateByPKEM(objOriginal, objValue, db.Exec, "tabUpdate", "schema")

	if err != nil {
		t.Error(err)
		return
	}

	objValue.B = "1"
	objValue.C = "0"
	affectedRows, err := UpdateByPKEM(objOriginal, objValue, db.Exec, "tabUpdate", "schema")

	if err != nil {
		t.Error(err)
		return
	}
	if affectedRows == 0 {
		t.Fail()
	}
}

func TestUpdateLike(t *testing.T) {
	objValue := testType{A: "0", B: "1", C: ""}

	fields, _ := getSchemaAttributes(objValue, "schema")
	likeFields, _ := getSchemaLike(objValue, "schema")

	exp := "update `tabUpdate` set `a` = ?, `b` = ?, `c` = ? where `c` like concat('%', ?, '%');"
	query := getUpdateQueryByLike(fields, likeFields, objValue, "tabUpdate", "schema")

	if query != exp {
		t.Errorf("\ngot ```%s```\nexp ```%s```", query, exp)
	}
}
