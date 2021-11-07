package symbols

import "encoding/binary"

// AddStrings add every string in given store, convert slice of getted symbols to string and add it to store
func AddStrings(store interface{ AddString(string) Symbol }, ss []string) Symbol {
	res := make([]byte, 8*(len(ss)+1))
	n := binary.PutUvarint(res, uint64(len(ss)))
	for _, s := range ss {
		n += binary.PutUvarint(res[n:], uint64(store.AddString(s)))
	}
	return store.AddString(string(res[:n]))
}

// GetStrings restore slice of string added in store with AddStrings
func GetStrings(store interface{ GetString(Symbol) string }, sym Symbol) []string {
	encoded := []byte(store.GetString(sym))
	length, n := binary.Uvarint(encoded)
	res := make([]string, length)
	for i := range res {
		s, k := binary.Uvarint(encoded[n:])
		n += k
		res[i] = store.GetString(Symbol(s))
	}
	return res
}
