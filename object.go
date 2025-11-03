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

	//Runs during (Layer).RunThinkers() call.
	OnUpdate() func(self Gobject, l Layer)

	//Runs every (Layer).Draw() call.
	OnDraw() func(self Gobject, i *ebiten.Image, l Layer)

	//Runs when Gobject gets added to a layer.
	OnAdd() func(self Gobject, l Layer)

	//Runs on (*FreeLayer).DeleteGobject and (*GridLayer).DeleteAt calls,
	//but not when gobjects get overwritten by new or moving Gobjects in a
	//GridLayer.
	OnDelete() func(self Gobject, l Layer)

	//Default sprite drawing function.
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

	//Runs every (Layer).Draw() call.
	OnDrawFunc func(self Gobject, i *ebiten.Image, l Layer)
	//Runs during (Layer).RunThinkers() call.
	OnUpdateFunc func(self Gobject, l Layer)
	//Runs when Gobject gets added to a layer.
	OnAddFunc func(self Gobject, l Layer)
	//Runs on (*FreeLayer).DeleteGobject and (*GridLayer).DeleteAt calls,
	//but not when gobjects get overwritten by new or moving Gobjects in a
	//GridLayer. Also doesn't run during (Layer).Clear().
	OnDeleteFunc func(self Gobject, l Layer)

	gobjectElem *list.Element // Referenced by FreeLayer
	thinkerElem *list.Element // Referenced by thinker list in layers
}

// Create a new BaseGobject. Use BaseGobject.Build() to create copies for that
// can be added to layers.
func NewBaseGobject(name string, sprites SpritePack) BaseGobject {
	return BaseGobject{name,
		imggg.Pt[int](0, 0),
		imggg.Pt[float64](0, 0),
		sprites, nil, nil, nil, nil, nil, nil}
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

func (o BaseGobject) OnDraw() func(Gobject, *ebiten.Image, Layer) {
	return o.OnDrawFunc
}

func (o BaseGobject) OnUpdate() func(Gobject, Layer) {
	return o.OnUpdateFunc
}

func (o BaseGobject) OnAdd() func(Gobject, Layer) {
	return o.OnAddFunc
}

func (o BaseGobject) OnDelete() func(Gobject, Layer) {
	return o.OnDeleteFunc
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
