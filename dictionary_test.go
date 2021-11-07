package symbols_test

import (
	"testing"
	"testing/quick"

	"github.com/vporoshok/symbols"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestDictionary(t *testing.T) {
	src := make([]string, 1e4)
	for i := range src {
		src[i] = gofakeit.Sentence(gofakeit.Number(1, 100))
	}
	dict := symbols.Dictionary{}
	sym := make([]symbols.Symbol, len(src))
	for i := range sym {
		sym[i] = dict.AddString(src[i])
	}
	for i := range src {
		assert.Equal(t, src[i], dict.GetString(sym[i]))
	}
}

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
	str := "This library helps to store big number of strings in structure" +
		"with small number of pointers to make it friendly to Go garbage collector."
	for i := 0; i < b.N; i++ {
		_ = ss.AddString(str)
	}
}
