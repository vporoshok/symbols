package symbols

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"

	"github.com/pierrec/lz4/v4"
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

type dumpHeader struct {
	Index, Length, Pages, Longs uint32
}

// Dump data with comression
func (store Store) Dump() io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		lw := lz4.NewWriter(pw)
		defer func() {
			_ = lw.Close()
			_ = pw.Close()
		}()
		_ = binary.Write(lw, binary.LittleEndian, dumpHeader{
			Index:  uint32(store.index),
			Length: uint32(store.length),
			Pages:  uint32(len(store.pages)),
			Longs:  uint32(len(store.longStrings)),
		})
		for i := range store.pages {
			_, _ = lw.Write(store.pages[i][:])
		}
		for i := range store.longStrings {
			_ = binary.Write(lw, binary.LittleEndian, uint32(len(store.longStrings[i])))
			_, _ = lw.Write([]byte(store.longStrings[i]))
		}
	}()
	return pr
}

// Restore data from compressed dump
func Restore(r io.Reader) (Store, error) {
	var header dumpHeader
	lr := lz4.NewReader(r)
	if err := binary.Read(lr, binary.LittleEndian, &header); err != nil {
		return Store{}, fmt.Errorf("read header: %w", err)
	}
	store := Store{
		index:       int(header.Index),
		length:      int(header.Length),
		pages:       make([][PageSize]byte, header.Pages),
		longStrings: make([]string, header.Longs),
	}
	for i := range store.pages {
		if _, err := lr.Read(store.pages[i][:]); err != nil {
			return store, fmt.Errorf("read %d page: %w", i, err)
		}
	}
	for i := range store.longStrings {
		var l int32
		if err := binary.Read(lr, binary.LittleEndian, &l); err != nil {
			return store, fmt.Errorf("read %d long string length: %w", i, err)
		}
		buf := make([]byte, l)
		if _, err := lr.Read(buf); err != nil {
			return store, fmt.Errorf("read %d long string: %w", i, err)
		}
		//nolint:gosec // reuse exists memory to reduce allocations
		store.longStrings[i] = *(*string)(unsafe.Pointer(&buf))
	}
	return store, nil
}

// StoreState static information of the store
type StoreState struct {
	Symbols, Pages, LongStrings int
}

// State of the store
func (store Store) State() StoreState {
	return StoreState{
		Symbols:     store.length,
		Pages:       len(store.pages),
		LongStrings: len(store.longStrings),
	}
}
