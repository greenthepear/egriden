package egriden

import (
	"slices"
	"testing"
)

func TestNeighbors(t *testing.T) {
	lv := NewBaseLevel("test")
	l1 := lv.CreateSimpleGridLayerOnTop("layer 1", 10, 5, 5, Sparse, 0, 0)

	c1 := l1.CellAt(0, 1)
	n := c1.GetNeighbors(Bishop, true, true)
	// |   n
	// | c
	// |   n
	if len(n) != 3 {
		t.Errorf("Wrong neighbors slice len in: %v", n)
	}
	if !slices.Contains(n, c1) {
		t.Errorf("Missing origin cell in: %v", n)
	}
	if !slices.Contains(n, l1.CellAt(1, 0)) || !slices.Contains(n, l1.CellAt(1, 2)) {
		t.Errorf("Missing cell in: %v", n)
	}

	l1.AddGobject(NewBaseGobject("hello!", EmptySpritePack()).Build(), 1, 1)
	// |
	// |   o
	// |     c
	n2 := l1.CellAt(2, 2).GetNeighborsSetFunc(King, false, true,
		func(c Cell) bool {
			if c.HasGobject() && c.Gobject().Name() == "hello!" {
				return true
			}
			return false
		},
	)
	if len(n2) != 1 {
		t.Errorf("Wrong neighbors set len of: %v", n2)
	}
	for k := range n2 {
		if k.HasGobject() && k.Gobject().Name() != "hello!" {
			t.Errorf("Wrong gobject in neighbor set: %v", n2)
		}
	}
}
