package egriden

import "github.com/hajimehoshi/ebiten/v2"

func (l GridLayer) draw(screen *ebiten.Image) {
	switch l.mode {
	case Sparce:
		for vec, o := range l.mapMat {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				float64(vec.x)+l.xOffset,
				float64(vec.y)+l.yOffset)

			screen.DrawImage(o.Sprite(), op)
		}
	case Dense:
		for y := range l.height {
			for x := range l.width {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(
					float64(x)+l.xOffset,
					float64(y)+l.yOffset)

				screen.DrawImage(
					l.sliceMat[y][x].Sprite(), op)
			}
		}
	case Static:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			l.xOffset,
			l.yOffset)
		screen.DrawImage(l.staticImage, op)
	}
}
