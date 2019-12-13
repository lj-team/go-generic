package params

import (
	"testing"
)

func TestParams(t *testing.T) {

	params := Params(map[string]string{
		"uid":   "123",
		"float": "12.5",
		"hello": "hello",
		"other": "val",
		"bool1": "true",
		"bool2": "1",
		"bool3": "0",
	})

	if params.GetString("hello") != "hello" {
		t.Fatal("GetString faliled")
	}

	if params.GetString("hello1") != "" {
		t.Fatal("GetString faliled")
	}

	if params.GetFloat("float") != 12.5 {
		t.Fatal("GetFloat faliled")
	}

	if params.GetFloat("uid") != 123 {
		t.Fatal("GetFloat faliled")
	}

	if params.GetFloat("float2") != 0 {
		t.Fatal("GetFloat faliled")
	}

	if params.GetFloat("hello") != 0 {
		t.Fatal("GetFloat faliled")
	}

	if params.GetInt("uid") != 123 {
		t.Fatal("GetInt faliled")
	}

	if params.GetInt("hello") != 0 {
		t.Fatal("GetInt faliled")
	}

	if params.GetInt("hello1") != 0 {
		t.Fatal("GetInt faliled")
	}

	if !params.GetBool("bool1") {
		t.Fatal("GetInt faliled")
	}

	if !params.GetBool("bool2") {
		t.Fatal("GetInt faliled")
	}

	if params.GetBool("bool3") {
		t.Fatal("GetInt faliled")
	}

	if params.GetBool("hello1") {
		t.Fatal("GetBool faliled")
	}

	if params.GetBool("hello") {
		t.Fatal("GetBool faliled")
	}
}
