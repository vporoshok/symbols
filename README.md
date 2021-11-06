# Symbols

<!-- [![Travis Build](https://travis-ci.com/vporoshok/project.svg?branch=master)](https://travis-ci.com/vporoshok/project) -->
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

Short strings (less or equal 255 bytes) stored in pages (1 MB byte arrays) with prefix byte with it length, for example:
```
[4, J, o, h, n, 9, d, e, v, e, l, o, p, e, r, 6, g, o, l, a, n, g, 7, b, c, y, c, l, e, 4, M, a, r, y, ...]
```
And Symbol for it represents two uin32 numbers (page number and index of byte with length) stored as an uint64.

If string to be add is longer then rest of current page, a new page is being created, so pages may be incomplete (less than 256 bytes per page, less than 0.01%).

Long strings stored as is in separate slice and Symbol for it represents index in this slice with 1 at first bit.

### Dictionary

## Key features

List of features that may help to decide use this project.

## See also

List of inspiration projects, alternatives, and analogs. Maybe some comparisons if it's important.

## Roadmap

- Configure Github Actions for project;
- Store string length in Symbol instead of page;
- Optimize long string store;
- May be add methods to save/restore Store to/from file or writer/reader (how to use it?);
