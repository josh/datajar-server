package shortcuts

import (
	"testing"
)

func TestConvertToJSONPath(t *testing.T) {
	if ConvertToJSONPath("") != "" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("/") != "" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("foo") != "foo" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("/foo") != "foo" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("foo/bar") != "foo.bar" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("/foo/bar") != "foo.bar" {
		t.Error("ConvertToJSONPath failed")
	}

	if ConvertToJSONPath("foo/bar/baz") != "foo.bar.baz" {
		t.Error("ConvertToJSONPath failed")
	}
}
