package symbols

import (
	"unsafe"
)

const (
	// PageSize 1MB
	PageSize = 1 << 20
	// LongGuard length of string to store outside of pages and not to check on duplicates
	LongGuard = 255
)

// Store of strings
type Store struct {
	index, length int
	pages         [][PageSize]byte
	longStrings   []string
}

// Symbol reference to string stored in Symbols
type Symbol uint64

// AddString to store
func (store *Store) AddString(s string) Symbol {
	store.length++
	if len(s) > LongGuard {
		i := len(store.longStrings)
		store.longStrings = append(store.longStrings, s)
		return 0x8000000000000000 | Symbol(i)
	}
	if len(store.pages) == 0 || store.index+len(s)+1 > PageSize {
		store.pages = append(store.pages, [PageSize]byte{})
		store.index = 0
	}
	pageNumber := len(store.pages) - 1
	index := store.index
	store.pages[pageNumber][index] = byte(len(s))
	store.index += 1 + copy(store.pages[pageNumber][index+1:], s)
	return Symbol(index) | Symbol(pageNumber)<<32
}

// GetString from store by reference
func (store Store) GetString(sym Symbol) string {
	if sym>>63 == 1 {
		return store.longStrings[sym^1<<63]
	}
	page := store.pages[sym>>32][:]
	index := int(sym << 32 >> 32)
	length := int(page[index])
	b := page[index+1 : index+1+length]
	//nolint:gosec // reuse exists memory to reduce allocations
	return *(*string)(unsafe.Pointer(&b))
}

// Len count of strings in store
func (store Store) Len() int {
	return store.length
}
