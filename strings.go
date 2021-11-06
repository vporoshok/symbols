package symbols

import "encoding/binary"

func AddStrings(store interface{ AddString(string) Symbol }, ss []string) Symbol {
	res := make([]byte, 8*(len(ss)+1))
	n := binary.PutUvarint(res, uint64(len(ss)))
	for _, s := range ss {
		n += binary.PutUvarint(res[n:], uint64(store.AddString(s)))
	}
	return store.AddString(string(res[:n]))
}

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
