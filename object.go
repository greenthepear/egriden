package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Gobject is an object that exists in a layer
type Gobject interface {
	Name() string
	XY() (int, int)
	setXY(int, int)

	//Sprite stuff

	IsVisible() bool
	Sprite() *ebiten.Image
	SetImageSequence(string) error
	NextFrame()
	SetFrame(int)

	//Custom scripts

	OnUpdate() func()                        //Runs every game.Update() call
	OnDraw() func(*ebiten.Image, *GridLayer) //Runs ever game.Draw() call
	DrawSprite(*ebiten.Image, *GridLayer)    //Default sprite drawing function

	//objects are referenced outside of the grid sometimes, if they get deleted from it, these must be called and checked

	isMarkedForDeletion() bool
	markForDeletion()
}

// The BaseGobject. Use it for simple Gobjects or implement your own Gobject by embedding this struct in your own.
type BaseGobject struct {
	name string
	x, y int

	sprites SpritePack

	markedForDeletion bool
}

// Create a new BaseGobject. Use BaseGobject.Build() to create a scriptless Gobject that can be added to a layer
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

// Assigns Sprite Pack. Should not be used during game updates.
func (o *BaseGobject) SetSpritePack(sp SpritePack) {
	o.sprites = sp
}

// Sets Image Sequence under name, returns error if the name key is not present
func (o *BaseGobject) SetImageSequence(name string) error {
	_, ok := o.sprites.sequences[name]
	if !ok {
		return fmt.Errorf("ImageSequence '%v' doesn't exist", name)
	}
	o.sprites.currentSequenceKey = name
	return nil
}

// Sets the frame to `i % len(frames)`
func (o *BaseGobject) SetFrame(i int) {
	o.sprites.frameIndex =
		i % len(o.sprites.sequences[o.sprites.currentSequenceKey].frames)
}

// Increments frame by one. Wraps back to 0 if out of range.
func (o *BaseGobject) NextFrame() {
	o.SetFrame(o.sprites.frameIndex + 1)
}

func (o *BaseGobject) IsVisible() bool {
	return o.sprites.visible
}

// Returns the current sprite. Used for built-in drawing.
func (o *BaseGobject) Sprite() *ebiten.Image {
	return o.sprites.Sprite()
}

func (o *BaseGobject) isMarkedForDeletion() bool {
	return o.markedForDeletion
}

func (o *BaseGobject) markForDeletion() {
	o.markedForDeletion = true
}

type baseGobjectWithoutScripts struct {
	BaseGobject
}

func (o *baseGobjectWithoutScripts) OnUpdate() func() {
	return nil
}

func (o *baseGobjectWithoutScripts) OnDraw() func(*ebiten.Image, *GridLayer) {
	return nil
}

// Default function for drawing the sprite in the grid, shouldn't be overwritten.
func (o *BaseGobject) DrawSprite(on *ebiten.Image, l *GridLayer) {
	x, y := o.XY()
	on.DrawImage(o.Sprite(),
		createDrawImageOptionsForXY(
			float64(x)*float64(l.SquareLength)+l.XOffset,
			float64(y)*float64(l.SquareLength)+l.YOffset))
}

// Makes a BaseGobject respect Gobject interface by assigning it nil draw update etc. functions
//
// For Gobjects with custom scripts you must create your own build function.
func (o BaseGobject) Build() Gobject {
	return &baseGobjectWithoutScripts{
		BaseGobject: o,
	}
}

/// Layer interactions

// Returns Gobject at x y, nil if empty. Panics if out of bounds.
func (l GridLayer) GobjectAt(x, y int) Gobject {
	if !l.IsXYwithinBounds(x, y) {
		panic("GobjectAt() panic! Out of bounds.")
	}
	if l.mode == Sparce {
		return l.mapMat[vec{x, y}]
	}
	return l.sliceMat[y][x]
}

func (l GridLayer) IsOccupiedAt(x, y int) bool {
	return l.GobjectAt(x, y) != nil
}

// Adds Gobject to the layer at x y. Will overwrite the any existing Gobject there.
func (l *GridLayer) AddGobject(o Gobject, x, y int) {
	o.setXY(x, y)
	if l.mode == Sparce {
		if l.mapMat[vec{x, y}] != nil {
			l.mapMat[vec{x, y}].markForDeletion()
		}
		l.mapMat[vec{x, y}] = o
		return
	}
	if l.sliceMat[y][x] != nil {
		l.sliceMat[y][x].markForDeletion()
	}
	l.sliceMat[y][x] = o
}

func (l *GridLayer) internalDeleteAt(x, y int, markForDeletion bool) {
	if !l.IsXYwithinBounds(x, y) {
		panic("not within layer bounds")
	}

	if l.mode == Sparce {
		if markForDeletion {
			l.mapMat[vec{x, y}].markForDeletion()
		}
		delete(l.mapMat, vec{x, y})
		return
	}

	if markForDeletion {
		l.sliceMat[y][x].markForDeletion()
	}
	l.sliceMat[y][x] = nil
}

func (l *GridLayer) DeleteAt(x, y int) {
	l.internalDeleteAt(x, y, true)
}

func (l *GridLayer) MoveGobjectTo(o Gobject, x, y int) {
	if !l.IsXYwithinBounds(x, y) {
		panic("not within layer bounds")
	}
	fromX, fromY := o.XY()
	fromGobject := l.GobjectAt(fromX, fromY)
	if fromGobject != o {
		panic(fmt.Sprintf(
			`Gobject '%s' is not the same as in the layer (%p != %p).
			Are you referencing one from another layer or one that wasn't added yet?`,
			o.Name(), o, fromGobject,
		))
	}

	l.internalDeleteAt(fromX, fromY, false)
	l.AddGobject(o, x, y)
}
