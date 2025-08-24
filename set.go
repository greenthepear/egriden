package egriden

import "slices"

// Simple ordered set implementation
type gobjectSet struct {
	keys []Gobject
	m    map[Gobject]struct{}
}

func newGobjectSet() gobjectSet {
	return gobjectSet{
		make([]Gobject, 0),
		make(map[Gobject]struct{}),
	}
}

func (set *gobjectSet) Add(o Gobject) {
	set.keys = append(set.keys, o)
	set.m[o] = struct{}{}
}

func (set *gobjectSet) Delete(o Gobject) {
	set.keys = slices.DeleteFunc(set.keys,
		func(fo Gobject) bool {
			return fo == o
		})
	delete(set.m, o)
}
