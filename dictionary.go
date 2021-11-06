package symbols

import "github.com/cespare/xxhash/v2"

type Dictionary struct {
	Store
	shortIndex map[uint64]Symbol
	longIndex  map[string]Symbol
}

func (dict *Dictionary) AddString(s string) Symbol {
	if len(s) > LongGuard {
		return dict.Store.AddString(s)
	}
	h := xxhash.Sum64String(s)
	switch sym, ok := dict.shortIndex[h]; {
	case ok && dict.Store.GetString(sym) == s:
		return sym
	case !ok:
		sym = dict.Store.AddString(s)
		if dict.shortIndex == nil {
			dict.shortIndex = make(map[uint64]Symbol)
		}
		dict.shortIndex[h] = sym
		return sym
	}
	if sym, ok := dict.longIndex[s]; ok {
		return sym
	}
	sym := dict.Store.AddString(s)
	if dict.longIndex == nil {
		dict.longIndex = make(map[string]Symbol)
	}
	dict.longIndex[s] = sym
	return sym
}
