package tmpl

import (
	"testing"
)

func TestTmpl(t *testing.T) {
	tmpl := New(nil, true)
	if tmpl == nil {
		t.Fatal("create tample error")
	}

	text := "<html><body>Hello, {{ name | noescape }}!</body></html>"

	if LoadTmplString("hello.jet", text) != nil {
		t.Fatal("load valid tmpl error")
	}

	vars := tmpl.Vars()

	vars.Set("name", "World")

	res, err := Render("hello.jet", vars)
	if err != nil {
		t.Fatal("Render return error")
	}

	if res != "<html><body>Hello, World!</body></html>" {
		t.Fatal("Render invalid result")
	}

	defTmpl = nil

	_, err = Render("hello.jet", vars)
	if err == nil {
		t.Fatal("Render must return error")
	}

	if LoadTmplString("hello.jet", text) == nil {
		t.Fatal("LoadTmplString must return error")
	}

	Set(tmpl)
	if defTmpl == nil {
		t.Fatal("Set not work")
	}

	_, err = Render("not_found.jet", vars)
	if err == nil {
		t.Fatal("Template not found. But error not found")
	}

	vars = Vars()

	_, err = Render("hello.jet", vars)
	if err == nil {
		t.Fatal("Render must return error. Var not found")
	}
}
