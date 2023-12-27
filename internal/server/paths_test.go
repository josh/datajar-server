package server

import (
	"reflect"
	"testing"
)

func TestGetValueByPath(t *testing.T) {
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
	if value, _ := GetValueByPath(store, "/foo/bar"); reflect.DeepEqual(value, "baz") != true {
		t.Error("expected /foo/bar to be baz, got", value)
	}
	if value, _ := GetValueByPath(store, "/bar/answer"); reflect.DeepEqual(value, 42) != true {
		t.Error("expected /bar/answer to be 42, got", value)
	}
	if value, err := GetValueByPath(store, "/foo/baz"); value != nil || err == nil {
		t.Error("expected /foo/baz to be nil, got", value, err)
	}
	if value, err := GetValueByPath(store, "/foo/baz/biz"); value != nil || err == nil {
		t.Error("expected /foo/baz/biz to be nil, got", value, err)
	}
	if value, _ := GetValueByPath(store, "/foo"); reflect.DeepEqual(value, foo) != true {
		t.Error("expected /foo to be foo, got", value)
	}
	if value, _ := GetValueByPath(store, "/"); reflect.DeepEqual(value, store) != true {
		t.Error("expected / to be store, got", value)
	}
	if value, _ := GetValueByPath(store, ""); reflect.DeepEqual(value, store) != true {
		t.Error("expected / to be store, got", value)
	}
}
