package acl

import (
	"reflect"
	"testing"
)

func TestCanReadPath(t *testing.T) {
	var caps []Capabilities

	caps = []Capabilities{{Read: []string{"foo"}}}
	if CanReadPath("/foo", caps) != true {
		t.Error("foo could not read /foo")
	}

	caps = []Capabilities{{Read: []string{"foo"}}}
	if CanReadPath("/bar", caps) != false {
		t.Error("foo could read /bar")
	}

	caps = []Capabilities{{Read: []string{"foo", "bar"}}}
	if CanReadPath("/foo", caps) != true {
		t.Error("foo could not read /foo")
	}
	if CanReadPath("/bar", caps) != true {
		t.Error("bar could not read /bar")
	}

	caps = []Capabilities{{Read: []string{"foo/*"}}}
	if CanReadPath("/foo", caps) != false {
		t.Error("foo/* could read /foo")
	}
	if CanReadPath("/foo/bar", caps) != true {
		t.Error("foo/* could not read /foo/bar")
	}
	if CanReadPath("/foo/baz", caps) != true {
		t.Error("foo/* could not read /foo/baz")
	}

	caps = []Capabilities{{Read: []string{"*"}}}
	if CanReadPath("/", caps) != true {
		t.Error("* could not read /")
	}
	if CanReadPath("/foo", caps) != true {
		t.Error("* could not read /foo")
	}
	if CanReadPath("/foo/bar", caps) != true {
		t.Error("* could not read /foo/bar")
	}
}

func TestGetPath(t *testing.T) {
	foo := map[string]interface{}{
		"bar": "baz",
	}
	bar := map[string]interface{}{
		"answer": 42,
	}
	store := map[string]interface{}{
		"foo": foo,
		"bar": bar,
	}
	if value := GetPath(store, "/foo/bar"); reflect.DeepEqual(value, "baz") != true {
		t.Error("expected /foo/bar to be baz, got", value)
	}
	if value := GetPath(store, "/bar/answer"); reflect.DeepEqual(value, 42) != true {
		t.Error("expected /bar/answer to be 42, got", value)
	}
	if value := GetPath(store, "/foo/baz"); value != nil {
		t.Error("expected /foo/baz to be nil, got", value)
	}
	if value := GetPath(store, "/foo"); reflect.DeepEqual(value, foo) != true {
		t.Error("expected /foo to be foo, got", value)
	}
	if value := GetPath(store, "/"); reflect.DeepEqual(value, store) != true {
		t.Error("expected / to be store, got", value)
	}
	if value := GetPath(store, ""); reflect.DeepEqual(value, store) != true {
		t.Error("expected / to be store, got", value)
	}
}
