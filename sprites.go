package egriden

import (
	"embed"
	"fmt"
	"image"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ImageSequence is a sequence of images aka frames.
type ImageSequence struct {
	Name   string
	Frames []*ebiten.Image
}

// SpritePack is a collection of ImageSequences and controls things like the
// frame index.
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
//
// Uses [ebitenutil.NewImageFromFile], which means you might need to import
// decoders like `import _ "image/png"`
func CreateImageSequenceFromPaths(name string, paths ...string) (ImageSequence, error) {
	if len(paths) < 1 {
		return ImageSequence{}, fmt.Errorf("no paths provided")
	}
	frameSlice := make([]*ebiten.Image, len(paths))
	for i, p := range paths {
		img, _, err := ebitenutil.NewImageFromFile(p)
		if err != nil {
			return ImageSequence{}, fmt.Errorf("while importing file from path `%v`: %v", p, err)
		}
		frameSlice[i] = img
	}
	return ImageSequence{name, frameSlice}, nil
}

// Searches for PNG files in the folder and creates an ImageSequence,
// with frame order based on the alphabetical order of the file names.
//
// Deprecated: Use [CreateImageSequenceFromGlob]
func CreateImageSequenceFromFolder(
	name, folderPath string) (ImageSequence, error) {

	found, err := filepath.Glob(folderPath + "/*.png")
	if err != nil {
		return ImageSequence{}, fmt.Errorf(
			"while searching for files in folder `%v`: %v", folderPath, err)
	}
	if len(found) == 0 {
		return ImageSequence{}, fmt.Errorf(
			"no PNG files found in `%v`", folderPath)
	}
	return CreateImageSequenceFromPaths(name, found...)
}

// Searches for files using a glob pattern (such as `Graphics/*.png`) with [filepath.Glob].
// Uses [ebitenutil.NewImageFromFile], which means you might need to import decoders like `import _ "image/png"`
func CreateImageSequenceFromGlob(
	name, globPattern string) (ImageSequence, error) {

	found, err := filepath.Glob(globPattern)
	if err != nil {
		return ImageSequence{}, fmt.Errorf(
			"while searching for files using pattern `%v`: %v", globPattern, err)
	}
	if len(found) == 0 {
		return ImageSequence{}, fmt.Errorf(
			"no PNG files found using `%v`", globPattern)
	}
	return CreateImageSequenceFromPaths(name, found...)
}

// Creates image sequence from std's image.Image using ebiten.NewImageFromImage
func CreateImageSequenceFromImages(
	name string, images ...image.Image) (ImageSequence, error) {

	if len(images) < 1 {
		return ImageSequence{}, fmt.Errorf("no images provided")
	}
	frameSlice := make([]*ebiten.Image, len(images))

	for i, img := range images {
		frameSlice[i] = ebiten.NewImageFromImage(img)
	}

	return ImageSequence{name, frameSlice}, nil
}

func openAndDecode(in embed.FS, path string) (image.Image, error) {
	f, err := in.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

func openAndDecodeMany(in embed.FS, paths ...string) ([]image.Image, error) {
	r := make([]image.Image, len(paths))

	for i, path := range paths {
		img, err := openAndDecode(in, path)
		if err != nil {
			return nil, err
		}
		r[i] = img
	}

	return r, nil
}

// Create a map mapping SpritePack names to built SpritePacks, populated with
// ImageSequences containing frames from files in an [embed.FS],
// based on YAML data.
// The YAML requires a specific structure, example of which:
//
//	spritepacks:
//	- name: card bases
//	  sequences:
//	    - name: back
//	      paths:
//	        - "Graphics/card_base_back.png"
//	    - name: front
//	      paths:
//	        - "Graphics/CardsBase/card_front_base_pink.png"
//	        - "Graphics/CardsBase/card_front_base_green.png"
//	        - "Graphics/CardsBase/card_front_base_blue.png"
//	        - "Graphics/CardsBase/card_front_base_orange.png"
//	- name: card patterns
//	...
func CreateSpritePacksFromYaml(
	yamlBytes []byte, fs embed.FS) (map[string]SpritePack, error) {

	var data struct {
		Spritepacks []struct {
			Name      string
			Sequences []struct {
				Name  string
				Paths []string
			}
		}
	}
	err := yaml.Unmarshal(yamlBytes, &data)
	if err != nil {
		return nil, err
	}

	spritemap := make(map[string]SpritePack, len(data.Spritepacks))
	for _, spritepack := range data.Spritepacks {
		finalSequences := make([]ImageSequence, len(spritepack.Sequences))
		for i, sequence := range spritepack.Sequences {
			images, err := openAndDecodeMany(fs, sequence.Paths...)
			if err != nil {
				return nil, err
			}
			createdSequence, err := CreateImageSequenceFromImages(
				sequence.Name, images...)
			if err != nil {
				return nil, err
			}
			finalSequences[i] = createdSequence
		}
		spritemap[spritepack.Name] = NewSpritePackWithSequences(
			finalSequences...)
	}
	return spritemap, nil
}

func NewSpritePack() SpritePack {
	return SpritePack{
		sequences:  make(map[string]*ImageSequence),
		frameIndex: 0,

		currentSequenceKey: "",

		visible: true,
		empty:   false, // TODO: maybe remove empty, its confusing

		DrawOptions: &ebiten.DrawImageOptions{},
		Offset:      imggg.Point[float64]{},
	}
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

// Return specific frame of a sequence
func (sp SpritePack) FrameAt(sequenceKey string, frame int) *ebiten.Image {
	return sp.sequences[sequenceKey].Frames[frame]
}

// Return the current sprite
func (sp SpritePack) Sprite() *ebiten.Image {
	if !sp.visible {
		return nil
	}

	return sp.FrameAt(sp.currentSequenceKey, sp.frameIndex)
}
