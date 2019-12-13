package main

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	a := []byte("\b021B69")

	b := []byte("\x1B\x69")

	d, err := decode(a)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(d, b) {
		t.Fatal("different")
	}
}

func TestDecodeLong(t *testing.T) {
	a := []byte("o€\b00000000021B69€")

	b := []byte("o€\x1B\x69€")

	d, err := decode(a)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(d, b) {
		t.Fatal("different")
	}
}
