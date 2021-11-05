package symbols_test

import (
	"strings"
	"testing"
	"testing/quick"

	"github.com/vporoshok/symbols"
)

func TestSymbols_short(t *testing.T) {
	ss := symbols.Symbols{}
	foo := ss.AddString("foo")
	bar := ss.AddString("bar")
	if ss.GetString(foo) != "foo" {
		t.Logf(`expected "foo", got %q`, ss.GetString(foo))
		t.Fail()
	}
	if ss.GetString(bar) != "bar" {
		t.Logf(`expected "bar", got %q`, ss.GetString(bar))
		t.Fail()
	}
}

func TestSymbols_long(t *testing.T) {
	ss := symbols.Symbols{}
	foo := ss.AddString(strings.Repeat("foo", 100))
	bar := ss.AddString(strings.Repeat("bar", 100))
	if ss.GetString(foo) != strings.Repeat("foo", 100) {
		t.Logf(`expected "foo"*100, got %q`, ss.GetString(foo))
		t.Fail()
	}
	if ss.GetString(bar) != strings.Repeat("bar", 100) {
		t.Logf(`expected "bar"*100, got %q`, ss.GetString(bar))
		t.Fail()
	}
}

func TestSymbols_quick(t *testing.T) {
	ss := symbols.Symbols{}
	fn := func(s string) bool {
		sym := ss.AddString(s)
		return ss.GetString(sym) == s
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSymbolsAdd(b *testing.B) {
	ss := symbols.Symbols{}
	b.ReportAllocs()
	str := "foo asdf a daf asdf asdfasdfa asdfd"
	for i := 0; i < b.N; i++ {
		_ = ss.AddString(str)
	}
}
