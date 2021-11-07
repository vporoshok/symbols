package symbols_test

import (
	"testing"

	"github.com/vporoshok/symbols"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	src := make([]string, gofakeit.Number(10, 100))
	for i := range src {
		src[i] = gofakeit.Sentence(gofakeit.Number(1, 100))
	}
	ss := symbols.Store{}
	sym := symbols.AddStrings(&ss, src)
	assert.Equal(t, src, symbols.GetStrings(ss, sym))
}

func TestDictionaryStrings(t *testing.T) {
	src := make([]string, gofakeit.Number(10, 100))
	for i := range src {
		src[i] = gofakeit.Sentence(gofakeit.Number(1, 10))
	}
	ss := symbols.Dictionary{}
	sym := symbols.AddStrings(&ss, src)
	assert.Equal(t, src, symbols.GetStrings(ss, sym))
	assert.Equal(t, sym, symbols.AddStrings(&ss, src))
}
