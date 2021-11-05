package symbols

import "github.com/cespare/xxhash/v2"

type Dictionary struct {
	Symbols
	shortIndex map[uint64]Symbol
	longIndex  map[string]Symbol
}

func (dict *Dictionary) AddString(s string) Symbol {
	if len(s) > LongGuard {
		return dict.Symbols.AddString(s)
	}
	h := xxhash.Sum64String(s)
	switch sym, ok := dict.shortIndex[h]; {
	case ok && dict.Symbols.GetString(sym) == s:
		return sym
	case !ok:
		sym = dict.Symbols.AddString(s)
		if dict.shortIndex == nil {
			dict.shortIndex = make(map[uint64]Symbol)
		}
		dict.shortIndex[h] = sym
		return sym
	}
	if sym, ok := dict.longIndex[s]; ok {
		return sym
	}
	sym := dict.Symbols.AddString(s)
	if dict.longIndex == nil {
		dict.longIndex = make(map[string]Symbol)
	}
	dict.longIndex[s] = sym
	return sym
}
