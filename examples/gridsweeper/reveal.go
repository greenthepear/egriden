package main

import (
	"github.com/greenthepear/egriden"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type oNumber struct {
	egriden.BaseGobject
	count int
}

func NewNumber(n int) egriden.Gobject {
	o := oNumber{
		egriden.NewBaseGobject("count", egriden.EmptySpritePack()),
		n}
	o.OnDrawFunc = func(
		self egriden.Gobject, i *ebiten.Image, l egriden.Layer) {

		selfcell := l.(*egriden.GridLayer).CellAt(self.GridPos().XY()).Anchor()
		op := text.DrawOptions{}
		op.GeoM.Translate(selfcell.XY())
		text.Draw(
			i, string(o.count), &text.GoTextFace{
				Source: silkscreen_reg,
				Size:   8,
			}, &op)
	}
	return o.Build()
}

func (g *Game) CountBombs() {
	l := g.Level().GridLayer(bombZ)
	w, h := l.Dimensions()
	for x := range w {
		for y := range h {
			cell := l.CellAt(x, y)
			if cell.HasGobject() {
				continue
			}
			bombCount := len(cell.GetNeighborsFunc(
				egriden.King, false, true,
				func(c egriden.Cell) bool {
					return c.HasGobject() && c.Gobject().Name() == "bomb"
				},
			))
			if bombCount == 0 {
				continue
			}
			l.AddGobject(NewNumber(bombCount), x, y)
		}
	}
}

// Shamelessly stolen from https://stackoverflow.com/questions/28541609
type cellStack []egriden.Cell

func (s *cellStack) Push(v egriden.Cell) {
	*s = append(*s, v)
}

func (s *cellStack) Pop() egriden.Cell {
	r := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return r
}

func (g *Game) RevealFlood(c egriden.Cell) {
	stack := make(cellStack, 0)
	stack.Push(c)

	for len(stack) != 0 {
		cell := stack.Pop()

		neighbors := cell.GetNeighborsSetFunc(egriden.King, false, true,
			func(c egriden.Cell) bool {
				if !c.HasGobject() {
					return false
				}
				o := g.Level().GridLayer(bombZ).CellAt(c.XY()).Gobject()
				if o == nil {
					return true
				}
				return o.Name() != "count"
			},
		)
		for neighbor := range neighbors {
			stack.Push(neighbor)
			g.Level().GridLayer(revealZ).DeleteAt(neighbor.XY())
		}
	}
}
