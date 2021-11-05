package symbols

import (
	"unsafe"
)

const (
	PageSize  = 1 << 20
	LongGuard = 255
)

type Symbols struct {
	index, length int
	pages         [][PageSize]byte
	longStrings   []string
}

// Symbol reference to string stored in Symbols
//
// If first bit is 1 (s >> 63 == 1) this is long string placed in longStrings store.
// Otherwise Symbol should be interpreted as 31-bit number of page and 32-bit index in this page.
type Symbol uint64

func (symbols *Symbols) AddString(s string) Symbol {
	symbols.length++
	if len(s) > LongGuard {
		i := len(symbols.longStrings)
		symbols.longStrings = append(symbols.longStrings, s)
		return Symbol(1<<63) | Symbol(i)
	}
	if len(symbols.pages) == 0 || symbols.index+len(s)+1 > PageSize {
		symbols.pages = append(symbols.pages, [PageSize]byte{})
		symbols.index = 0
	}
	pageNumber := len(symbols.pages) - 1
	index := symbols.index
	symbols.pages[pageNumber][index] = byte(len(s))
	symbols.index += 1 + copy(symbols.pages[pageNumber][index+1:], s)
	return Symbol(index) | Symbol(pageNumber)<<32
}

func (symbols Symbols) GetString(sym Symbol) string {
	if sym>>63 == 1 {
		return symbols.longStrings[sym^1<<63]
	}
	page := symbols.pages[sym>>32]
	index := int(sym << 32 >> 32)
	length := int(page[index])
	b := page[index+1 : index+1+length]
	return *(*string)(unsafe.Pointer(&b))
}

func (symbols Symbols) Len() int {
	return symbols.length
}
