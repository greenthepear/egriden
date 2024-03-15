package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Gobject interface {
	Name() string
	XY() (int, int)
	setXY(int, int)

	Sprite() *ebiten.Image
	SetImageSequence(string) error
	NextFrame()
	SetFrame(int)
}

type BaseGobject struct {
	name string
	x, y int

	sprites *SpritePack
}

func NewBaseGobject(name string, sprites SpritePack) BaseGobject {
	return BaseGobject{name, 0, 0, &sprites}
}

func (o *BaseGobject) Name() string {
	return o.name
}

func (o *BaseGobject) XY() (int, int) {
	return o.x, o.y
}

func (o *BaseGobject) setXY(x, y int) {
	o.x, o.y = x, y
}

func (o *BaseGobject) SetSpritePack(sp SpritePack) {
	o.sprites = &sp
}

func (o *BaseGobject) SetImageSequence(name string) error {
	_, ok := o.sprites.sequences[name]
	if !ok {
		return fmt.Errorf("ImageSequence '%v' doesn't exist", name)
	}
	o.sprites.currentSequenceKey = name
	return nil
}

func (o *BaseGobject) SetFrame(i int) {
	s := o.sprites.sequences[o.sprites.currentSequenceKey]
	s.frameIndex = i % len(s.frames)
}

func (o *BaseGobject) NextFrame() {
	o.SetFrame(
		o.sprites.sequences[o.sprites.currentSequenceKey].frameIndex + 1)
}

func (o *BaseGobject) Sprite() *ebiten.Image {
	s := o.sprites.sequences[o.sprites.currentSequenceKey]
	return s.frames[s.frameIndex]
}
