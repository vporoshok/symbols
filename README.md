# Symbols

[![Test](https://github.com/vporoshok/symbols/actions/workflows/test.yml/badge.svg)](https://github.com/vporoshok/symbols/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vporoshok/symbols)](https://goreportcard.com/report/github.com/vporoshok/symbols)
[![GoDoc](http://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/vporoshok/symbols)
[![codecov](https://codecov.io/gh/vporoshok/symbols/branch/main/graph/badge.svg)](https://codecov.io/gh/vporoshok/symbols)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)

> GC-friendly string store

This library helps to store big number of strings in structure with small number of pointers to make it friendly to Go garbage collector.

## Motivation

As described in article [Avoiding high GC overhead with large heaps](https://blog.gopheracademy.com/advent-2018/avoid-gc-overhead-large-heaps/), a large number of pointers make Go Garbage Collector life harder. So what if you need to deal with strings, a lot of strings? Every string is a pointer. You need to hide this pointers from Go runtime.

It's helpful for persistent caches and text resources.

### Philosophy

As small as possible, as fast as possible.

## Usage

```go
package main

import (
	"fmt"

	"github.com/vporoshok/symbols"
)

// User is a domain model
type User struct {
	ID     int
	Name   string
	Labels []string
}

// UserCache is a persistent cache of users
type UserCache struct {
	data map[int]cachedUser
	dict symbols.Dictionary // using Dictionary to deduplicate labels
}

// cachedUser is a special struct to glue symbols with domain model
type cachedUser struct {
	name, labels symbols.Symbol
}

// Add user to cache
func (cache *UserCache) Add(user User) {
	if _, ok := cache.data[user.ID]; !ok {
		if cache.data == nil {
			cache.data = make(map[int]cachedUser)
		}
		cache.data[user.ID] = cachedUser{
			name:   cache.dict.AddString(user.Name),
			labels: symbols.AddStrings(&cache.dict, user.Labels),
		}
	}
}

// Get user from cache by id
func (cache UserCache) Get(id int) (User, bool) {
	if user, ok := cache.data[id]; ok {
		return User{
			ID:     id,
			Name:   cache.dict.GetString(user.name),
			Labels: symbols.GetStrings(cache.dict, user.labels),
		}, true
	}
	return User{}, false
}

func main() {
	cache := new(UserCache)
	cache.Add(User{
		ID:     1,
		Name:   "John",
		Labels: []string{"developer", "golang", "bicycle"},
	})
	cache.Add(User{
		ID:     2,
		Name:   "Mary",
		Labels: []string{"developer", "golang", "running"},
	})
	cache.Add(User{
		ID:     3,
		Name:   "Albert",
		Labels: []string{"manager", "bicycle", "running"},
	})
	fmt.Println(cache.Get(1))
	fmt.Println(cache.Get(2))
	fmt.Println(cache.Get(3))
	// Dictionary now contain strings
	// John, Mary, Albert
	// developer, golang, bicycle, running, manager
	// and composition of John's, Mary's and Albert's labels as 3-byte strings
	fmt.Println(cache.dict.Len()) // 11
}

```

## How it works

### Store

Store has two parts:
1. slice of pages for small strings;
2. slice of long strings as is;

Short strings (less or equal 1<<12 - 1 bytes) stored in pages (1 MB byte arrays), for example:
```
[J, o, h, n, d, e, v, e, l, o, p, e, r, g, o, l, a, n, g, b, i, c, y, c, l, e, M, a, r, y, ...]
```
And Symbol for it represents next bits:
- 1 bit with zero (see long strings store);
- 31 bit with page number;
- 12 bit with length of string;
- 20 bit with index of string in the page;

If string to be add is longer then rest of current page, a new page is being created, so pages may be incomplete (less than 2KB per 1MB page, less than 0.2%).

Long strings stored as is in separate slice and Symbol for it represents index in this slice with 1 at first bit.

### Dictionary

Dictionary is a Store Wrapper to deduplicate strings. It use map of xxhashes of strings as short index and long index `map[string]Symbol` on collision short index.

## See also

Original Phil's Pearl projects [stringbank](https://github.com/philpearl/stringbank) and [intern](https://github.com/philpearl/intern)
- store all strings independ of it length in pages that may cause a lot of unused space;
- use slices as pages that add more pointers and unneeded length int per page;
- use maphash to deduplicate, that is slower than xxhash;

## Roadmap

- Optimize long string store;
