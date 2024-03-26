package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Gobject interface {
	Name() string
	XY() (int, int)
	setXY(int, int)

	//Sprite stuff
	Sprite() *ebiten.Image
	SetImageSequence(string) error
	NextFrame()
	SetFrame(int)

	//Custom scripts
	OnCreate() func()
	OnUpdate() func()
	OnDraw() func(*ebiten.Image)
	DoesDrawScriptOverwriteSprite() bool

	isMarkedForDeletion() bool
}

type BaseGobject struct {
	name string
	x, y int

	sprites SpritePack

	markedForDeletion bool
}

func NewBaseGobject(name string, sprites SpritePack) BaseGobject {
	return BaseGobject{name, 0, 0, sprites, false}
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
	o.sprites = sp
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
	o.sprites.frameIndex =
		i % len(o.sprites.sequences[o.sprites.currentSequenceKey].frames)
}

func (o *BaseGobject) NextFrame() {
	o.SetFrame(o.sprites.frameIndex + 1)
}

func (o *BaseGobject) Sprite() *ebiten.Image {
	s := o.sprites.sequences[o.sprites.currentSequenceKey]
	return s.frames[o.sprites.frameIndex]
}

func (o *BaseGobject) isMarkedForDeletion() bool {
	return o.markedForDeletion
}

type BaseGobjectWithoutScripts struct {
	BaseGobject
}

func (o *BaseGobjectWithoutScripts) OnCreate() func() {
	return nil
}

func (o *BaseGobjectWithoutScripts) OnUpdate() func() {
	return nil
}

func (o *BaseGobjectWithoutScripts) OnDraw() func(*ebiten.Image) {
	return nil
}

func (o *BaseGobjectWithoutScripts) DoesDrawScriptOverwriteSprite() bool {
	return false
}

func (o BaseGobject) Build() Gobject {
	return &BaseGobjectWithoutScripts{
		BaseGobject: o,
	}
}
