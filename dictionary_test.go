package symbols_test

import (
	"testing"
	"testing/quick"

	"github.com/vporoshok/symbols"
)

func TestDictionary_quick(t *testing.T) {
	ss := symbols.Dictionary{}
	fn := func(s string) bool {
		sym := ss.AddString(s)
		return ss.GetString(sym) == s
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkDictionaryAdd(b *testing.B) {
	ss := symbols.Dictionary{}
	b.ReportAllocs()
	str := "foo asdf a daf asdf asdfasdfa asdfd"
	for i := 0; i < b.N; i++ {
		_ = ss.AddString(str)
	}
}
