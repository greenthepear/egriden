package egriden

import (
	"container/list"
	"fmt"

	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Gobject is an object that exists in a layer
type Gobject interface {
	Name() string
	// Coordinates within a GridLayer, if Gobject is in a FreeLayer
	// this should be always (0, 0)
	GridPos() imggg.Point[int]
	setGridPos(int, int)
	// Position on the screen, in a FreeLayer it's just the position, for
	// GridLayer it's the draw anchor.
	ScreenPos(Layer) imggg.Point[float64]
	setScreenPos(float64, float64)

	//Sprite stuff

	IsVisible() bool
	Sprite() *ebiten.Image
	SpritePack() SpritePack
	SetDrawOptions(*ebiten.DrawImageOptions)
	SetDrawOffsets(float64, float64)
	SetImageSequence(string) error
	NextFrame()
	SetFrame(int)

	//Custom scripts

	//Runs during (Layer).RunThinkers() call
	OnUpdate() func(self Gobject, l Layer)
	//Runs every (Layer).Draw() call
	OnDraw() func(self Gobject, i *ebiten.Image, l Layer)
	//Default sprite drawing function
	DrawSprite(*ebiten.Image, Layer)

	thinkerElement() *list.Element
	setThinkerElement(*list.Element)

	gobjectElement() *list.Element
	setGobjectElement(*list.Element)
}

// The BaseGobject. Use it for simple Gobjects or implement your own Gobject by
// embedding this struct in your own.
type BaseGobject struct {
	name      string
	gridPos   imggg.Point[int]
	screenPos imggg.Point[float64]

	sprites SpritePack

	OnDrawFunc   func(self Gobject, i *ebiten.Image, l Layer)
	OnUpdateFunc func(self Gobject, l Layer)

	gobjectElem *list.Element // Referenced by FreeLayer
	thinkerElem *list.Element // Referenced by thinker list in layers
}

// Create a new BaseGobject. Use BaseGobject.Build() to create a scriptless
// Gobject that can be added to a layer
func NewBaseGobject(name string, sprites SpritePack) BaseGobject {
	return BaseGobject{name,
		imggg.Pt[int](0, 0),
		imggg.Pt[float64](0, 0), sprites, nil, nil, nil, nil}
}

func (o *BaseGobject) Name() string {
	return o.name
}

// Grid coordinate in a GridLayer, returns (0, 0) for free layer objects.
func (o *BaseGobject) GridPos() imggg.Point[int] {
	return o.gridPos
}

func (o *BaseGobject) setGridPos(x, y int) {
	o.gridPos.X, o.gridPos.Y = x, y
}

// Position on the screen, in a FreeLayer it's just the position, for
// GridLayer it's the draw anchor.
func (o *BaseGobject) ScreenPos(layer Layer) imggg.Point[float64] {
	xoffset, yoffset := layer.CurrentAnchor().XY()
	spriteXoffset, spriteYoffset := o.SpritePack().Offset.XY()
	switch l := layer.(type) {
	case *GridLayer:
		x, y := o.gridPos.X, o.gridPos.Y
		return imggg.Point[float64]{
			X: float64(x)*(float64(l.cellDimensions.Width)+l.Padding.X) +
				xoffset + spriteXoffset,
			Y: float64(y)*(float64(l.cellDimensions.Height)+l.Padding.Y) +
				yoffset + spriteYoffset,
		}
	case *FreeLayer:
		return o.screenPos
	default:
		panic("layer is not GridLayer or FreeLayer somehow")
	}
}

func (o *BaseGobject) setScreenPos(x, y float64) {
	o.screenPos.X, o.screenPos.Y = x, y
}

// Assigns Sprite Pack. Should not be used during game updates.
func (o *BaseGobject) SetSpritePack(sp SpritePack) {
	o.sprites = sp
}

// Sets the frame to `i % len(frames)`
func (o *BaseGobject) SetFrame(i int) {
	o.sprites.frameIndex =
		i % len(o.sprites.sequences[o.sprites.currentSequenceKey].Frames)
}

// Sets Image Sequence under name, returns error if the name key is not present.
//
// Resets frame to 0.
func (o *BaseGobject) SetImageSequence(name string) error {
	_, ok := o.sprites.sequences[name]
	if !ok {
		return fmt.Errorf("ImageSequence '%v' doesn't exist", name)
	}
	o.sprites.currentSequenceKey = name
	o.SetFrame(0)
	return nil
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

func (o *BaseGobject) SpritePack() SpritePack {
	return o.sprites
}

// Set custom ebiten draw options. Remember that tx and ty get translated
// depending on the grid position and layers offset.
func (o *BaseGobject) SetDrawOptions(op *ebiten.DrawImageOptions) {
	o.sprites.DrawOptions = op
}

// Quick way to make the sprite draw with x and y added to the screen position.
func (o *BaseGobject) SetDrawOffsets(x, y float64) {
	o.sprites.Offset = imggg.Pt(x, y)
}

func (o *BaseGobject) OnDraw() func(Gobject, *ebiten.Image, Layer) {
	return o.OnDrawFunc
}

func (o *BaseGobject) OnUpdate() func(Gobject, Layer) {
	return o.OnUpdateFunc
}

// Default function for drawing the sprite in the grid.
func (o *BaseGobject) DrawSprite(on *ebiten.Image, l Layer) {
	l.DrawSprite(o, on)
}

func (o BaseGobject) thinkerElement() *list.Element {
	return o.thinkerElem
}

func (o *BaseGobject) setThinkerElement(e *list.Element) {
	o.thinkerElem = e
}

func (o BaseGobject) gobjectElement() *list.Element {
	return o.gobjectElem
}

func (o *BaseGobject) setGobjectElement(e *list.Element) {
	o.gobjectElem = e
}

// Makes a copy of the Gobject.
func (o BaseGobject) Build() Gobject {
	copy := o
	if !o.SpritePack().empty {
		drawOp := *o.SpritePack().DrawOptions
		copy.sprites.DrawOptions = &drawOp
	} else {
		copy.sprites = EmptySpritePack()
	}
	return &copy
}

/// Layer interactions

// Returns Gobject at x y, nil if empty. Panics if out of bounds.
func (l GridLayer) GobjectAt(x, y int) Gobject {
	if !l.IsXYwithinBounds(x, y) {
		panic(fmt.Sprintf("GobjectAt() panic! (%d , %d) Out of bounds.", x, y))
	}
	if l.mode == Sparse {
		return l.mapMat[imggg.Point[int]{X: x, Y: y}]
	}
	return l.sliceMat[y][x]
}

func (l GridLayer) IsOccupiedAt(x, y int) bool {
	return l.GobjectAt(x, y) != nil
}

func (l *GridLayer) internalAddGobject(
	o Gobject, x, y int) {

	o.setGridPos(x, y)

	if o.OnUpdate() != nil {
		o.setThinkerElement(l.thinkers.PushBack(o))
	}

	if l.mode == Sparse {
		l.mapMat[imggg.Pt(x, y)] = o
		return
	}
	l.sliceMat[y][x] = o
}

// Adds Gobject to the layer at x y.
// Will overwrite the any existing Gobject there.
func (l *GridLayer) AddGobject(o Gobject, x, y int) {
	l.internalAddGobject(o, x, y)
}

func (l *GridLayer) internalDeleteAt(x, y int) {
	if !l.IsXYwithinBounds(x, y) {
		panic("not within layer bounds")
	}

	o := l.GobjectAt(x, y)
	if o != nil && o.thinkerElement() != nil {
		l.thinkers.Remove(o.thinkerElement())
	}

	if l.mode == Sparse {
		delete(l.mapMat, imggg.Pt(x, y))
		return
	}

	l.sliceMat[y][x] = nil
}

func (l *GridLayer) DeleteAt(x, y int) {
	l.internalDeleteAt(x, y)
}

// Moves Gobject by first finding itself in the layer with its XY coordinates,
// but will panic if the Gobject in that cell is not the same, so you cannot use
// this with Gobjects that are not in the layer, obviously.
func (l *GridLayer) MoveGobjectTo(o Gobject, x, y int) {
	if !l.IsXYwithinBounds(x, y) {
		panic("not within layer bounds")
	}
	fromX, fromY := o.GridPos().XY()
	fromGobject := l.GobjectAt(fromX, fromY)
	if fromGobject != o {
		panic(fmt.Sprintf(
			`Gobject '%s' is not the same as in the layer (%p != %p).
			Are you referencing one from another layer or one that wasn't added yet?`,
			o.Name(), o, fromGobject,
		))
	}

	l.internalDeleteAt(fromX, fromY)
	l.AddGobject(o, x, y)
}

// Swaps objects between two grid positions, if either is empty it will be
// basically the same as moving the object. Panics if out of bounds.
func (l *GridLayer) SwapGobjectsAt(x1, y1, x2, y2 int) {
	o1 := l.GobjectAt(x1, y1)
	o2 := l.GobjectAt(x2, y2)

	if o1 != nil {
		l.internalAddGobject(o1, x2, y2)
		if o2 == nil {
			l.internalDeleteAt(x1, y1)
		}
	}

	if o2 != nil {
		l.internalAddGobject(o2, x1, y1)
		if o1 == nil {
			l.internalDeleteAt(x2, y2)
		}
	}
}

// Swaps objects between two cells, if either is empty it will be basically the
// same as moving the object. Panics if out of bounds.
func (l *GridLayer) SwapObjectsAtCells(cell1, cell2 Cell) {
	x1, y1 := cell1.XY()
	x2, y2 := cell2.XY()
	l.SwapGobjectsAt(x1, y1, x2, y2)
}
