package symbols_test

import (
	"testing"
	"testing/quick"

	"github.com/vporoshok/symbols"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestSymbols(t *testing.T) {
	src := make([]string, gofakeit.Number(10, 100))
	for i := range src {
		src[i] = gofakeit.Sentence(gofakeit.Number(1, 100))
	}
	store := symbols.Store{}
	sym := make([]symbols.Symbol, len(src))
	for i := range sym {
		sym[i] = store.AddString(src[i])
	}
	for i := range src {
		assert.Equal(t, src[i], store.GetString(sym[i]))
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
	str := ""
	for i := 0; i < b.N; i++ {
		_ = store.AddString(str)
	}
}

func BenchmarkSymbolsGet(b *testing.B) {
	store := symbols.Store{}
	str := "This library helps to store big number of strings in structure" +
		"with small number of pointers to make it friendly to Go garbage collector."
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
