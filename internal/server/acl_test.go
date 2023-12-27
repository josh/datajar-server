package server

import (
	"testing"
)

func TestCanAccessPath(t *testing.T) {
	var caps []Capabilities

	caps = []Capabilities{{Read: []string{"foo"}}}
	if CanAccessPath("/foo", caps, "read") != true {
		t.Error("foo could not read /foo")
	}

	caps = []Capabilities{{Read: []string{"foo"}}}
	if CanAccessPath("/bar", caps, "read") != false {
		t.Error("foo could read /bar")
	}

	caps = []Capabilities{{Read: []string{"foo", "bar"}}}
	if CanAccessPath("/foo", caps, "read") != true {
		t.Error("foo could not read /foo")
	}
	if CanAccessPath("/bar", caps, "read") != true {
		t.Error("bar could not read /bar")
	}

	caps = []Capabilities{{Read: []string{"foo/*"}}}
	if CanAccessPath("/foo", caps, "read") != false {
		t.Error("foo/* could read /foo")
	}
	if CanAccessPath("/foo/", caps, "read") != false {
		t.Error("foo/* could read /foo")
	}
	if CanAccessPath("/foo/bar", caps, "read") != true {
		t.Error("foo/* could not read /foo/bar")
	}
	if CanAccessPath("/foo/baz", caps, "read") != true {
		t.Error("foo/* could not read /foo/baz")
	}

	caps = []Capabilities{{Read: []string{"*"}}}
	if CanAccessPath("/", caps, "read") != true {
		t.Error("* could not read /")
	}
	if CanAccessPath("/foo", caps, "read") != true {
		t.Error("* could not read /foo")
	}
	if CanAccessPath("/foo/bar", caps, "read") != true {
		t.Error("* could not read /foo/bar")
	}
}
