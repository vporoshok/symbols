package symbols

import (
	"unsafe"
)

const (
	// should be less than 32
	pageSizeLog = 20
	// PageSize 1MB
	PageSize = 1 << pageSizeLog
	// LongGuard length of string to store outside of pages and not to check on duplicates
	LongGuard = 1<<(32-pageSizeLog) - 1
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
		return Symbol(i) | 1<<63
	}
	if len(store.pages) == 0 || store.index+len(s) > PageSize {
		store.pages = append(store.pages, [PageSize]byte{})
		store.index = 0
	}
	pageNumber := len(store.pages) - 1
	index := store.index
	store.index += copy(store.pages[pageNumber][index:], s)
	return Symbol(pageNumber)<<32 | Symbol(index|len(s)<<pageSizeLog)
}

// GetString from store by reference
func (store Store) GetString(sym Symbol) string {
	if sym>>63 == 1 {
		return store.longStrings[sym^1<<63]
	}
	page := store.pages[sym>>32][:]
	index := sym & (1<<pageSizeLog - 1)
	length := (sym & (1<<32 - 1)) >> pageSizeLog
	b := page[index : index+length]
	//nolint:gosec // reuse exists memory to reduce allocations
	return *(*string)(unsafe.Pointer(&b))
}

// Len count of strings in store
func (store Store) Len() int {
	return store.length
}
