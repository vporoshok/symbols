package symbols_test

import (
	"testing"

	"github.com/vporoshok/symbols"
)

func TestStrings(t *testing.T) {
	src := []string{
		"foo",
		"bar",
		"buz",
		"qux",
	}
	ss := symbols.Store{}
	sym := symbols.AddStrings(&ss, src)
	res := symbols.GetStrings(ss, sym)
	if len(res) != len(src) {
		t.Fatalf("result length %d missmatch source %d", len(res), len(src))
	}
	for i := range res {
		if res[i] != src[i] {
			t.Logf(`expected %q but got %q`, src[i], res[i])
			t.Fail()
		}
	}
}

func TestDictionaryStrings(t *testing.T) {
	src := []string{
		"foo",
		"bar",
		"buz",
		"qux",
	}
	ss := symbols.Dictionary{}
	sym := symbols.AddStrings(&ss, src)
	res := symbols.GetStrings(ss, sym)
	if len(res) != len(src) {
		t.Fatalf("result length %d missmatch source %d", len(res), len(src))
	}
	for i := range res {
		if res[i] != src[i] {
			t.Logf(`expected %q but got %q`, src[i], res[i])
			t.Fail()
		}
	}
}
