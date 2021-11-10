package symbols

import "github.com/cespare/xxhash/v2"

// Dictionary deduplicate wrapper for Store
//
// Use it on filling Store and use dict.Store in runtime to reduce memory utilization.
type Dictionary struct {
	Store
	shortIndex map[uint64]Symbol
	longIndex  map[string]Symbol
}

// AddString check string to duplicate and return existed Symbol or create new
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

// DropIndex clean map to reduce memory utilization
func (dict *Dictionary) DropIndex() {
	dict.shortIndex = nil
	dict.longIndex = nil
}

// DictionaryState static information of the dictionary
type DictionaryState struct {
	StoreState
	ShortIndex, LongIndex int
}

// State of the dictionary
func (dict Dictionary) State() DictionaryState {
	return DictionaryState{
		StoreState: dict.Store.State(),
		ShortIndex: len(dict.shortIndex),
		LongIndex:  len(dict.longIndex),
	}
}
