package main

import (
	"testing"
)

func TestBasic(t *testing.T) {
	v, err := Parse("foo=bar")
	if err != nil {
		t.Fatal(err)
	}

	check("foo", "bar", v.Values, t)
}

func Test2(t *testing.T) {
	s := `
foo=bar
foo2=bar2

# this is a comment

[section1]
foo  = bar
foo3 = bar3

[section2]
foo=3 adsasdf asdf asdf
foo3=bar3`

	v, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	if len(v.Sections) != 2 {
		t.Fail()
	}

	if v.Sections[0] != "section1" {
		t.Fail()
	}

	if v.Sections[1] != "section2" {
		t.Fail()
	}

	check("foo", "bar", v.Values, t)
	check("foo2", "bar2", v.Values, t)
	check("section1.foo", "bar", v.Values, t)
	check("section1.foo3", "bar3", v.Values, t)
	check("section2.foo", "3 adsasdf asdf asdf", v.Values, t)
	check("section2.foo3", "bar3", v.Values, t)
}

func check(key, expected string, values map[string]string, t *testing.T) {
	v := values[key]
	if v != expected {
		t.Fatalf("Invalid value for %s: %s expected %s", key, v, expected)
	}
}
