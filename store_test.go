package symbols_test

import (
	"strings"
	"testing"
	"testing/quick"

	"github.com/vporoshok/symbols"
)

func TestSymbols_short(t *testing.T) {
	store := symbols.Store{}
	foo := store.AddString("foo")
	bar := store.AddString("bar")
	if store.GetString(foo) != "foo" {
		t.Logf(`expected "foo", got %q`, store.GetString(foo))
		t.Fail()
	}
	if store.GetString(bar) != "bar" {
		t.Logf(`expected "bar", got %q`, store.GetString(bar))
		t.Fail()
	}
}

func TestSymbols_long(t *testing.T) {
	store := symbols.Store{}
	foo := store.AddString(strings.Repeat("foo", 100))
	bar := store.AddString(strings.Repeat("bar", 100))
	if store.GetString(foo) != strings.Repeat("foo", 100) {
		t.Logf(`expected "foo"*100, got %q`, store.GetString(foo))
		t.Fail()
	}
	if store.GetString(bar) != strings.Repeat("bar", 100) {
		t.Logf(`expected "bar"*100, got %q`, store.GetString(bar))
		t.Fail()
	}
}

func TestSymbols_quick(t *testing.T) {
	store := symbols.Store{}
	fn := func(s string) bool {
		sym := store.AddString(s)
		return store.GetString(sym) == s
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSymbolsAdd(b *testing.B) {
	store := symbols.Store{}
	b.ReportAllocs()
	str := "foo asdf a daf asdf asdfasdfa asdfd"
	for i := 0; i < b.N; i++ {
		_ = store.AddString(str)
	}
}

func BenchmarkSymbolsGet(b *testing.B) {
	store := symbols.Store{}
	str := "foo asdf a daf asdf asdfasdfa asdfd"
	for i := 0; i < 1e6; i++ {
		_ = store.AddString(str)
	}
	sym := store.AddString(str)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := store.GetString(sym)
		_ = s
	}
}
