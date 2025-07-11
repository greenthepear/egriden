package egriden

import (
	"fmt"
	"image"
	_ "image/png"
	"path/filepath"

	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ImageSequence is a sequence of images aka frames.
type ImageSequence struct {
	Name   string
	Frames []*ebiten.Image
}

// SpritePack is a collection of ImageSequences and controls things like the frame index.
type SpritePack struct {
	sequences  map[string]*ImageSequence
	frameIndex int

	currentSequenceKey string
	visible            bool
	empty              bool

	DrawOptions *ebiten.DrawImageOptions
	Offset      imggg.Point[float64]
}

// Create an ImageSequence using multiple (or just one) file paths.
func CreateImageSequenceFromPaths(name string, paths ...string) (ImageSequence, error) {
	if len(paths) < 1 {
		return ImageSequence{}, fmt.Errorf("no paths provided")
	}
	frameSlice := make([]*ebiten.Image, 0, len(paths))
	for _, p := range paths {
		img, _, err := ebitenutil.NewImageFromFile(p)
		if err != nil {
			return ImageSequence{}, fmt.Errorf("while importing file from path `%v`: %v", p, err)
		}
		frameSlice = append(frameSlice, img)
	}
	return ImageSequence{name, frameSlice}, nil
}

// Searches for PNG files in the folder and creates an ImageSequence, with frame order
// based on the alphabetical order of the file names.
func CreateImageSequenceFromFolder(name, folderPath string) (ImageSequence, error) {
	found, err := filepath.Glob(folderPath + "/*.png")
	if err != nil {
		return ImageSequence{}, fmt.Errorf("while searching for files in folder `%v`: %v", folderPath, err)
	}
	if len(found) == 0 {
		return ImageSequence{}, fmt.Errorf("no PNG files found in `%v`", folderPath)
	}
	return CreateImageSequenceFromPaths(name, found...)
}

// Creates image sequence from std's image.Image using ebiten.NewImageFromImage
func CreateImageSequenceFromImages(
	name string, images ...image.Image) (ImageSequence, error) {

	if len(images) < 1 {
		return ImageSequence{}, fmt.Errorf("no images provided")
	}
	frameSlice := make([]*ebiten.Image, 0, len(images))

	for _, img := range images {
		frameSlice = append(frameSlice, ebiten.NewImageFromImage(img))
	}

	return ImageSequence{name, frameSlice}, nil
}

func NewSpritePack() SpritePack {
	return SpritePack{make(map[string]*ImageSequence), 0, "", true, false, &ebiten.DrawImageOptions{}, imggg.Point[float64]{}}
}

// Assigns an ImageSequence to SpritePack
func (ip *SpritePack) AddImageSequence(is ImageSequence) {
	ip.sequences[is.Name] = &is
	if ip.currentSequenceKey == "" {
		ip.currentSequenceKey = is.Name
	}
}

// A sprite pack that will not render anything
func EmptySpritePack() SpritePack {
	return SpritePack{visible: false, empty: true}
}

// Create SpritePack and assign sequence
func NewSpritePackWithSequence(is ImageSequence) SpritePack {
	ip := NewSpritePack()
	ip.AddImageSequence(is)
	return ip
}

// Create SpritePack and assign multiple sequences
func NewSpritePackWithSequences(is ...ImageSequence) SpritePack {
	ip := NewSpritePack()
	for _, seq := range is {
		ip.AddImageSequence(seq)
	}
	return ip
}

// Return specific sprite in a seqence and frame
func (sp SpritePack) SpriteAt(sequenceKey string, frame int) *ebiten.Image {
	return sp.sequences[sequenceKey].Frames[frame]
}

// Return the current sprite
func (sp SpritePack) Sprite() *ebiten.Image {
	if !sp.visible {
		return nil
	}

	return sp.SpriteAt(sp.currentSequenceKey, sp.frameIndex)
}
