package symbols_test

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

func Example() {
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
	// Output:
	// {1 John [developer golang bicycle]} true
	// {2 Mary [developer golang running]} true
	// {3 Albert [manager bicycle running]} true
	// 11
}
