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
	"testing"

	"github.com/mattia-cabrini/go-utility"
	_ "github.com/mattn/go-sqlite3"
)

type testType struct {
	A string `schema:"a" pk:"true"`
	B string `schema:"b" pk:"true"`
	C string `schema:"c" like:"true"`
	D string
	E string
}

func openTestConnection() *sql.DB {
	var db, err = sql.Open("sqlite3", "entity_test.db")

	if err != nil {
		utility.Mypanic(err)
	}

	return db
}

func TestGetSchemaAttribute(t *testing.T) {
	obj := testType{}

	fields, err := getSchemaAttributes(obj, "schema")

	if err != nil {
		t.Error(err)
		return
	}

	if fields[0] != 0 {
		t.Error(0)
	}
	if fields[1] != 1 {
		t.Error(1)
	}
	if fields[2] != 2 {
		t.Error(2)
	}

	_, err = getSchemaAttributes(0, "schema")

	if err == nil {
		t.Fail()
	}

	objNS := struct {
		A string `schema:"a" pk:"true"`
		b string `schema:"b" pk:"true"` // b unexported: must cause getSchemaAttributes to err out
		C string `schema:"c"`
		D string
		E string
	}{}

	_, err = getSchemaAttributes(objNS, "schema")

	if err == nil {
		t.Fail()
	}
}

func TestGetSchemaPK(t *testing.T) {
	obj := testType{}

	fields, err := getSchemaPK(obj, "schema")

	if err != nil {
		t.Error(err)
		return
	}

	if len(fields) != 2 {
		t.Error(-1)
		return
	}

	if fields[0] != 0 {
		t.Error(0)
	}
	if fields[1] != 1 {
		t.Error(1)
	}

	_, err = getSchemaPK(0, "schema")

	if err == nil {
		t.Fail()
	}

	objNS := struct {
		A string `schema:"a" pk:"true"`
		b string `schema:"b" pk:"true"` // b unexported: must cause getSchemaAttributes to err out
		C string `schema:"c"`
		D string
		E string
	}{}

	_, err = getSchemaPK(objNS, "schema")

	if err == nil {
		t.Fail()
	}
}

func TestGetSchemaLike(t *testing.T) {
	obj := testType{}

	fields, err := getSchemaLike(obj, "schema")

	if err != nil {
		t.Error(err)
		return
	}

	if len(fields) != 1 {
		t.Error(-1)
		return
	}

	if fields[0] != 2 {
		t.Error(0)
	}

	_, err = getSchemaLike(0, "schema")

	if err == nil {
		t.Fail()
	}

	objNS := struct {
		A string `schema:"a" pk:"true"`
		b string `schema:"b" pk:"true" like:"true"` // b unexported: must cause getSchemaAttributes to err out
		C string `schema:"c"`
		D string
		E string
	}{}

	_, err = getSchemaLike(objNS, "schema")

	if err == nil {
		t.Fail()
	}
}
