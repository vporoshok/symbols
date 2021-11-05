package symbols

import "encoding/binary"

func AddStrings(symbols interface{ AddString(string) Symbol }, ss []string) Symbol {
	res := make([]byte, 8*(len(ss)+1))
	n := binary.PutUvarint(res, uint64(len(ss)))
	for _, s := range ss {
		n += binary.PutUvarint(res[n:], uint64(symbols.AddString(s)))
	}
	return symbols.AddString(string(res[:n]))
}

func GetStrings(symbols interface{ GetString(Symbol) string }, sym Symbol) []string {
	encoded := []byte(symbols.GetString(sym))
	length, n := binary.Uvarint(encoded)
	res := make([]string, length)
	for i := range res {
		s, k := binary.Uvarint(encoded[n:])
		n += k
		res[i] = symbols.GetString(Symbol(s))
	}
	return res
}
